package handlers

import (
	"bytes"
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/asn1"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"io/ioutil"
	"math/big"
	"net/http"
	"os"
	"strings"
)

type RegisterRequest struct {
	Login string `json:"login"`
}

type CaRegisterResponseResult struct {
	Secret string `json:"secret"`
}

type CaRegisterResponse struct {
	Result   CaRegisterResponseResult `json:"result"`
	Errors   []string                 `json:"errors"` // TODO add code and message. Example: "errors":[{"code":0,"message":"Registration of 'abc' failed: Identity 'abc' is already registered"}]
	Messages []string                 `json:"messages"`
	Success  bool                     `json:"success"`
}

func PostCaRegisterHandler(ec echo.Context) error {
	c := ec.(*ApiContext)

	registerRequest := new(RegisterRequest)
	if err := c.Bind(registerRequest); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if registerRequest.Login == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "New user login is required")
	}

	jsonRegisterRequest := fmt.Sprintf(`{"id":"%s","type":"client","affiliation":""}`, registerRequest.Login)

	// load private key pem
	home, err := os.UserHomeDir()
	if err != nil {
		return errors.Wrap(err, "Cannot obtain user home dir")
	}
	restPath := home + "/.fabric-rest-api-go/"

	privatePem, err := ioutil.ReadFile(restPath + "/admin_private.pem")
	if err != nil {
		return errors.Wrap(err, "Key reading error")
	}

	// load signcert pem
	signCertPem, err := ioutil.ReadFile(restPath + "/admin_signcert.pem")
	if err != nil {
		return errors.Wrap(err, "Signcert reading error")
	}

	// generate body
	requestUri := "/register"
	requestMethod := "POST"
	requestBody := []byte (jsonRegisterRequest)

	// generate payload

	b64body := B64Encode(requestBody)
	b64signCert := B64Encode(signCertPem)
	b64uri := B64Encode([]byte(requestUri))
	payload := requestMethod + "." + b64uri + "." + b64body + "." + b64signCert

	hasher := sha256.New()
	hasher.Write([]byte(payload))
	payloadSha256 := hasher.Sum(nil)

	// decode private key from PEM
	block, _ := pem.Decode(privatePem)

	privateKey, err := parsePrivateKey(block.Bytes)

	privateKeyEC, _ := privateKey.(*ecdsa.PrivateKey)

	// sign payload hash
	ecSignature, err := signECDSA(privateKeyEC, payloadSha256)
	if err != nil {
		return errors.WithMessage(err, "Signature generation failure")
	}
	if len(ecSignature) == 0 {
		return errors.New("Signature creation failed. Signature must be different than nil")
	}

	b64sig := B64Encode(ecSignature)
	// Authorization token
	token := b64signCert + "." + b64sig

	caRegisterUrl := fmt.Sprintf("http://%s/register", c.Fsc().ApiConfig.Ca.Host)

	req, err := http.NewRequest("POST", caRegisterUrl, bytes.NewBufferString(jsonRegisterRequest))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	caRegisterResponse := CaRegisterResponse{}
	err = json.Unmarshal(body, &caRegisterResponse)
	if err != nil {
		return errors.Wrap(err, "CA response unmarshal error")
	}

	if !caRegisterResponse.Success {
		return errors.Errorf("CA response with errors: %s", strings.Join(caRegisterResponse.Errors, ", "))
	}

	if resp.StatusCode != http.StatusCreated {
		return errors.Errorf("CA returned status %s, but 201 expected. With message: %s", resp.Status, string(body))
	}

	return c.JSONBlob(http.StatusOK, []byte( fmt.Sprintf(`{"result": "%s"}`, caRegisterResponse.Result.Secret)))
}

// B64Encode base64 encodes bytes
func B64Encode(buf []byte) string {
	return base64.StdEncoding.EncodeToString(buf)
}

// B64Decode base64 decodes a string
func B64Decode(str string) (buf []byte, err error) {
	return base64.StdEncoding.DecodeString(str)
}

func signECDSA(k *ecdsa.PrivateKey, digest []byte) ([]byte, error) {
	r, s, err := ecdsa.Sign(rand.Reader, k, digest)

	if err != nil {
		return nil, err
	}

	s, _, err = ToLowS(&k.PublicKey, s)
	if err != nil {
		return nil, err
	}

	return MarshalECDSASignature(r, s)
}

var (
	// curveHalfOrders contains the precomputed curve group orders halved.
	// It is used to ensure that signature' S value is lower or equal to the
	// curve group order halved. We accept only low-S signatures.
	// They are precomputed for efficiency reasons.
	curveHalfOrders = map[elliptic.Curve]*big.Int{
		elliptic.P224(): new(big.Int).Rsh(elliptic.P224().Params().N, 1),
		elliptic.P256(): new(big.Int).Rsh(elliptic.P256().Params().N, 1),
		elliptic.P384(): new(big.Int).Rsh(elliptic.P384().Params().N, 1),
		elliptic.P521(): new(big.Int).Rsh(elliptic.P521().Params().N, 1),
	}
)

type ECDSASignature struct {
	R, S *big.Int
}

func MarshalECDSASignature(r, s *big.Int) ([]byte, error) {
	return asn1.Marshal(ECDSASignature{r, s})
}

// IsLow checks that s is a low-S
func IsLowS(k *ecdsa.PublicKey, s *big.Int) (bool, error) {
	halfOrder, ok := curveHalfOrders[k.Curve]
	if !ok {
		return false, fmt.Errorf("curve not recognized [%s]", k.Curve)
	}

	return s.Cmp(halfOrder) != 1, nil

}

func ToLowS(k *ecdsa.PublicKey, s *big.Int) (*big.Int, bool, error) {
	lowS, err := IsLowS(k, s)
	if err != nil {
		return nil, false, err
	}

	if !lowS {
		// Set s to N - s that will be then in the lower part of signature space
		// less or equal to half order
		s.Sub(k.Params().N, s)

		return s, true, nil
	}

	return s, false, nil
}

func parsePrivateKey(der []byte) (crypto.PrivateKey, error) {
	if key, err := x509.ParsePKCS1PrivateKey(der); err == nil {
		return key, nil
	}
	if key, err := x509.ParsePKCS8PrivateKey(der); err == nil {
		switch key := key.(type) {
		case *ecdsa.PrivateKey:
			return key, nil
		default:
			return nil, fmt.Errorf("Found unknown private key type in PKCS#8 wrapping")
		}
	}
	if key, err := x509.ParseECPrivateKey(der); err == nil {
		return key, nil
	}
	return nil, fmt.Errorf("Failed to parse private key")
}
