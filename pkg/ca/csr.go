package ca

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/pem"
	"github.com/pkg/errors"
)

var oidEmailAddress = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 9, 1}

func CsrPEM(privateKey *ecdsa.PrivateKey, enrollRequest *ApiEnrollRequest) (string, error) {
	// TODO: Fill CSR details from enrollRequest
	subj := pkix.Name{
		CommonName:         enrollRequest.Login,
		Country:            []string{"US"},
		Province:           []string{"North Carolina"},
		Organization:       []string{"Hyperledger"},
		OrganizationalUnit: []string{"Fabric"},
	}
	rawSubj := subj.ToRDNSequence()
	var emailAddresses []string
	if enrollRequest.Email != "" {
		rawSubj = append(rawSubj, []pkix.AttributeTypeAndValue{
			{Type: oidEmailAddress, Value: enrollRequest.Email},
		})

		emailAddresses = []string{enrollRequest.Email}
	}

	asn1Subj, _ := asn1.Marshal(rawSubj)
	template := x509.CertificateRequest{
		RawSubject:         asn1Subj,
		EmailAddresses:     emailAddresses,
		SignatureAlgorithm: x509.ECDSAWithSHA256,
		DNSNames:           []string{},
	}

	csrBytes, err := x509.CreateCertificateRequest(rand.Reader, &template, privateKey)
	if err != nil {
		return "", errors.Wrap(err, "x509 CSR creating error")
	}

	buf := new(bytes.Buffer)

	err = pem.Encode(buf, &pem.Block{Type: "CERTIFICATE REQUEST", Bytes: csrBytes})
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func TbsCsrBytes(caTbsCsrRequest *CaTbsCsrRequest, publicKey *ecdsa.PublicKey) ([]byte, error){
	// TODO: Fill CSR details from enrollRequest
	subj := pkix.Name{
		CommonName:         caTbsCsrRequest.Login,
		Country:            []string{"Belarus"},
		Province:           []string{"Minsk"},
		Organization:       []string{"Altoros"},
		OrganizationalUnit: []string{"Altoros"},
	}
	rawSubj := subj.ToRDNSequence()
	var emailAddresses []string
	if caTbsCsrRequest.Email != "" {
		rawSubj = append(rawSubj, []pkix.AttributeTypeAndValue{
			{Type: oidEmailAddress, Value: caTbsCsrRequest.Email},
		})

		emailAddresses = []string{caTbsCsrRequest.Email}
	}

	asn1Subj, _ := asn1.Marshal(rawSubj)
	template := x509.CertificateRequest{
		RawSubject:         asn1Subj,
		EmailAddresses:     emailAddresses,
		SignatureAlgorithm: x509.ECDSAWithSHA256,
		DNSNames:           []string{},
	}

	tbsCsrBytes, err := CreateTbsCsr(&template, publicKey)
	if err != nil {
		return nil, errors.Wrap(err, "x509 CSR creating error")
	}

	return tbsCsrBytes, nil
}