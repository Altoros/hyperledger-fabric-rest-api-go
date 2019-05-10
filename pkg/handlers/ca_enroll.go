package handlers

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

var oidEmailAddress = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 9, 1}

type EnrollRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type CaEnrollRequest struct {
	CertificateRequest string `json:"certificate_request"`
	Profile            string `json:"profile"`
	CrlOverride        string `json:"crl_override"`
	Label              string `json:"Label"`
	//"NotBefore": "0001-01-01T00:00:00Z",
	//"NotAfter": "0001-01-01T00:00:00Z",
	CAName string `json:"CAName"`
}

type CaEnrollResponseResult struct {
	Cert       string                           `json:"Cert"`
	ServerInfo CaEnrollResponseResultServerInfo `json:"ServerInfo"`
}

type CaEnrollResponseResultServerInfo struct {
	CAName                    string `json:"CAName"`
	CAChain                   string `json:"CAChain"`
	IssuerPublicKey           string `json:"IssuerPublicKey"`
	IssuerRevocationPublicKey string `json:"IssuerRevocationPublicKey"`
	Version                   string `json:"Version"`
}

type CaEnrollResponse struct {
	Result   CaEnrollResponseResult `json:"result"`
	Errors   []string               `json:"errors"`
	Messages []string               `json:"messages"`
	Success  bool                   `json:"success"`
}

func PostCaEnrollHandler(ec echo.Context) error {
	c := ec.(*ApiContext)

	enrollRequest := new(EnrollRequest)
	if err := c.Bind(enrollRequest); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Private key generation
	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	// Storing private key in PEM format
	privatePem, _ := PrivateKeyToPEM(privateKey)

	home, err := os.UserHomeDir()
	if err != nil {
		return errors.Wrap(err, "Cannot obtain user home dir")
	}
	restPath := home + "/.fabric-rest-api-go/"

	err = os.MkdirAll(restPath, os.ModePerm)
	if err != nil {
		return errors.Wrap(err, "Key saving error")
	}
	err = ioutil.WriteFile(restPath+"/admin_private.pem", privatePem, 0644)
	if err != nil {
		return errors.Wrap(err, "Key saving error")
	}

	// create CSR
	emailAddress := "admin@example.com"
	subj := pkix.Name{
		CommonName:         enrollRequest.Login,
		Country:            []string{"US"},
		Province:           []string{"North Carolina"},
		Organization:       []string{"Hyperledger"},
		OrganizationalUnit: []string{"Fabric"},
	}
	rawSubj := subj.ToRDNSequence()
	rawSubj = append(rawSubj, []pkix.AttributeTypeAndValue{
		{Type: oidEmailAddress, Value: emailAddress},
	})

	asn1Subj, _ := asn1.Marshal(rawSubj)
	template := x509.CertificateRequest{
		RawSubject: asn1Subj,
		EmailAddresses:     []string{emailAddress},
		SignatureAlgorithm: x509.ECDSAWithSHA256,
		DNSNames:           []string{},
	}

	csrBytes, err := x509.CreateCertificateRequest(rand.Reader, &template, privateKey)
	if err != nil {
		return errors.Wrap(err, "CSR creating error")
	}

	buf := new(bytes.Buffer)

	err = pem.Encode(buf, &pem.Block{Type: "CERTIFICATE REQUEST", Bytes: csrBytes})
	if err != nil {
		return err
	}

	// send CSR to CA /enroll
	enrollCaRequest := CaEnrollRequest{
		CertificateRequest: buf.String(),
	}

	jsonEnrollRequest, err := json.Marshal(enrollCaRequest)
	if err != nil {
		return err
	}

	caEnrollUrl := fmt.Sprintf("http://%s/api/v1/enroll", c.Fsc().ApiConfig.Ca.Host)

	req, err := http.NewRequest("POST", caEnrollUrl, bytes.NewBuffer(jsonEnrollRequest))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(enrollRequest.Login, enrollRequest.Password)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return errors.Errorf("CA returned status %s, but 201 expected", resp.Status)
	}

	body, _ := ioutil.ReadAll(resp.Body)

	caEnrollResponse := CaEnrollResponse{}
	err = json.Unmarshal(body, &caEnrollResponse)
	if err != nil {
		return errors.Wrap(err, "CA response unmarshal error")
	}

	if !caEnrollResponse.Success {
		return errors.Errorf("CA response with errors: %s", strings.Join(caEnrollResponse.Errors, ", "))
	}

	signCertPem, _ := B64Decode(caEnrollResponse.Result.Cert)

	err = ioutil.WriteFile(restPath+"/admin_signcert.pem", signCertPem, 0644)
	if err != nil {
		return errors.Wrap(err, "Key saving error")
	}

	return c.JSONBlob(http.StatusOK, []byte( fmt.Sprintf(`{"result": "%s"}`, caEnrollResponse.Result.Cert)))
}

// PrivateKeyToPEM converts the private key to PEM format.
// EC private keys are converted to PKCS#8 format.
// RSA private keys are converted to PKCS#1 format.
func PrivateKeyToPEM(privateKey *ecdsa.PrivateKey) ([]byte, error) {

	// Validate inputs
	if privateKey == nil {
		return nil, errors.New("Invalid key. It must be different from nil.")
	}

	// get the oid for the curve
	oidNamedCurve, ok := oidFromNamedCurve(privateKey.Curve)
	if !ok {
		return nil, errors.New("unknown elliptic curve")
	}

	// based on https://golang.org/src/crypto/x509/sec1.go
	privateKeyBytes := privateKey.D.Bytes()

	paddedPrivateKey := make([]byte, (privateKey.Curve.Params().N.BitLen()+7)/8)
	copy(paddedPrivateKey[len(paddedPrivateKey)-len(privateKeyBytes):], privateKeyBytes)
	// omit NamedCurveOID for compatibility as it's optional
	asn1Bytes, err := asn1.Marshal(ecPrivateKey{
		Version:    1,
		PrivateKey: paddedPrivateKey,
		PublicKey:  asn1.BitString{Bytes: elliptic.Marshal(privateKey.Curve, privateKey.X, privateKey.Y)},
	})

	if err != nil {
		return nil, fmt.Errorf("error marshaling EC key to asn1 [%s]", err)
	}

	var pkcs8Key pkcs8Info
	pkcs8Key.Version = 0
	pkcs8Key.PrivateKeyAlgorithm = make([]asn1.ObjectIdentifier, 2)
	pkcs8Key.PrivateKeyAlgorithm[0] = oidPublicKeyECDSA
	pkcs8Key.PrivateKeyAlgorithm[1] = oidNamedCurve
	pkcs8Key.PrivateKey = asn1Bytes

	pkcs8Bytes, err := asn1.Marshal(pkcs8Key)
	if err != nil {
		return nil, fmt.Errorf("error marshaling EC key to asn1 [%s]", err)
	}
	return pem.EncodeToMemory(
		&pem.Block{
			Type:  "PRIVATE KEY",
			Bytes: pkcs8Bytes,
		},
	), nil
}

func oidFromNamedCurve(curve elliptic.Curve) (asn1.ObjectIdentifier, bool) {
	switch curve {
	case elliptic.P224():
		return oidNamedCurveP224, true
	case elliptic.P256():
		return oidNamedCurveP256, true
	case elliptic.P384():
		return oidNamedCurveP384, true
	case elliptic.P521():
		return oidNamedCurveP521, true
	}
	return nil, false
}

var (
	oidNamedCurveP224 = asn1.ObjectIdentifier{1, 3, 132, 0, 33}
	oidNamedCurveP256 = asn1.ObjectIdentifier{1, 2, 840, 10045, 3, 1, 7}
	oidNamedCurveP384 = asn1.ObjectIdentifier{1, 3, 132, 0, 34}
	oidNamedCurveP521 = asn1.ObjectIdentifier{1, 3, 132, 0, 35}
)

// struct to hold info required for PKCS#8
type pkcs8Info struct {
	Version             int
	PrivateKeyAlgorithm []asn1.ObjectIdentifier
	PrivateKey          []byte
}

type ecPrivateKey struct {
	Version       int
	PrivateKey    []byte
	NamedCurveOID asn1.ObjectIdentifier `asn1:"optional,explicit,tag:0"`
	PublicKey     asn1.BitString        `asn1:"optional,explicit,tag:1"`
}

var oidPublicKeyECDSA = asn1.ObjectIdentifier{1, 2, 840, 10045, 2, 1}
