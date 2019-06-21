package ca

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"github.com/pkg/errors"
	"net"
	"net/url"
	"strconv"
)

func marshalPublicKey(pub interface{}) (publicKeyBytes []byte, publicKeyAlgorithm pkix.AlgorithmIdentifier, err error) {
	switch pub := pub.(type) {
	case *ecdsa.PublicKey:
		publicKeyBytes = elliptic.Marshal(pub.Curve, pub.X, pub.Y)
		oid, ok := oidFromNamedCurve(pub.Curve)
		if !ok {
			return nil, pkix.AlgorithmIdentifier{}, errors.New("x509: unsupported elliptic curve")
		}
		publicKeyAlgorithm.Algorithm = oidPublicKeyECDSA
		var paramBytes []byte
		paramBytes, err = asn1.Marshal(oid)
		if err != nil {
			return
		}
		publicKeyAlgorithm.Parameters.FullBytes = paramBytes
	default:
		return nil, pkix.AlgorithmIdentifier{}, errors.New("x509: only RSA and ECDSA public keys supported")
	}

	return publicKeyBytes, publicKeyAlgorithm, nil
}

// These structures reflect the ASN.1 structure of X.509 certificates.:

type publicKeyInfo struct {
	Raw       asn1.RawContent
	Algorithm pkix.AlgorithmIdentifier
	PublicKey asn1.BitString
}


type PublicKeyAlgorithm int

const (
	UnknownPublicKeyAlgorithm PublicKeyAlgorithm = iota
	RSA
	DSA
	ECDSA
)

var publicKeyAlgoName = [...]string{
	RSA:   "RSA",
	DSA:   "DSA",
	ECDSA: "ECDSA",
}

func (algo PublicKeyAlgorithm) String() string {
	if 0 < algo && int(algo) < len(publicKeyAlgoName) {
		return publicKeyAlgoName[algo]
	}
	return strconv.Itoa(int(algo))
}

var (
	oidSignatureECDSAWithSHA256 = asn1.ObjectIdentifier{1, 2, 840, 10045, 4, 3, 2}
)



// RFC 3279, 2.3 Public Key Algorithms
//
// pkcs-1 OBJECT IDENTIFIER ::== { iso(1) member-body(2) us(840)
//    rsadsi(113549) pkcs(1) 1 }
//
// rsaEncryption OBJECT IDENTIFIER ::== { pkcs1-1 1 }
//
// id-dsa OBJECT IDENTIFIER ::== { iso(1) member-body(2) us(840)
//    x9-57(10040) x9cm(4) 1 }
//
// RFC 5480, 2.1.1 Unrestricted Algorithm Identifier and Parameters
//
// id-ecPublicKey OBJECT IDENTIFIER ::= {
//       iso(1) member-body(2) us(840) ansi-X9-62(10045) keyType(2) 1 }
var (
	oidPublicKeyECDSA = asn1.ObjectIdentifier{1, 2, 840, 10045, 2, 1}
)

// RFC 5480, 2.1.1.1. Named Curve
//
// secp224r1 OBJECT IDENTIFIER ::= {
//   iso(1) identified-organization(3) certicom(132) curve(0) 33 }
//
// secp256r1 OBJECT IDENTIFIER ::= {
//   iso(1) member-body(2) us(840) ansi-X9-62(10045) curves(3)
//   prime(1) 7 }
//
// secp384r1 OBJECT IDENTIFIER ::= {
//   iso(1) identified-organization(3) certicom(132) curve(0) 34 }
//
// secp521r1 OBJECT IDENTIFIER ::= {
//   iso(1) identified-organization(3) certicom(132) curve(0) 35 }
//
// NB: secp256r1 is equivalent to prime256v1
var (
	oidNamedCurveP224 = asn1.ObjectIdentifier{1, 3, 132, 0, 33}
	oidNamedCurveP256 = asn1.ObjectIdentifier{1, 2, 840, 10045, 3, 1, 7}
	oidNamedCurveP384 = asn1.ObjectIdentifier{1, 3, 132, 0, 34}
	oidNamedCurveP521 = asn1.ObjectIdentifier{1, 3, 132, 0, 35}
)

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

// KeyUsage represents the set of actions that are valid for a given key. It's
// a bitmap of the KeyUsage* constants.
type KeyUsage int

const (
	KeyUsageDigitalSignature KeyUsage = 1 << iota
	KeyUsageContentCommitment
	KeyUsageKeyEncipherment
	KeyUsageDataEncipherment
	KeyUsageKeyAgreement
	KeyUsageCertSign
	KeyUsageCRLSign
	KeyUsageEncipherOnly
	KeyUsageDecipherOnly
)

// ExtKeyUsage represents an extended set of actions that are valid for a given key.
// Each of the ExtKeyUsage* constants define a unique action.
type ExtKeyUsage int

const (
	ExtKeyUsageAny ExtKeyUsage = iota
	ExtKeyUsageServerAuth
	ExtKeyUsageClientAuth
	ExtKeyUsageCodeSigning
	ExtKeyUsageEmailProtection
	ExtKeyUsageIPSECEndSystem
	ExtKeyUsageIPSECTunnel
	ExtKeyUsageIPSECUser
	ExtKeyUsageTimeStamping
	ExtKeyUsageOCSPSigning
	ExtKeyUsageMicrosoftServerGatedCrypto
	ExtKeyUsageNetscapeServerGatedCrypto
	ExtKeyUsageMicrosoftCommercialCodeSigning
	ExtKeyUsageMicrosoftKernelCodeSigning
)

// ErrUnsupportedAlgorithm results from attempting to perform an operation that
// involves algorithms that are not currently implemented.
var ErrUnsupportedAlgorithm = errors.New("x509: cannot verify signature: algorithm unimplemented")


// ConstraintViolationError results when a requested usage is not permitted by
// a certificate. For example: checking a signature when the public key isn't a
// certificate signing key.
type ConstraintViolationError struct{}

func (ConstraintViolationError) Error() string {
	return "x509: invalid signature: parent certificate cannot sign this kind of certificate"
}


type UnhandledCriticalExtension struct{}

func (h UnhandledCriticalExtension) Error() string {
	return "x509: unhandled critical extension"
}

const (
	nameTypeEmail = 1
	nameTypeDNS   = 2
	nameTypeURI   = 6
	nameTypeIP    = 7
)

var (
	oidExtensionSubjectKeyId          = []int{2, 5, 29, 14}
	oidExtensionKeyUsage              = []int{2, 5, 29, 15}
	oidExtensionExtendedKeyUsage      = []int{2, 5, 29, 37}
	oidExtensionAuthorityKeyId        = []int{2, 5, 29, 35}
	oidExtensionBasicConstraints      = []int{2, 5, 29, 19}
	oidExtensionSubjectAltName        = []int{2, 5, 29, 17}
	oidExtensionCertificatePolicies   = []int{2, 5, 29, 32}
	oidExtensionNameConstraints       = []int{2, 5, 29, 30}
	oidExtensionCRLDistributionPoints = []int{2, 5, 29, 31}
	oidExtensionAuthorityInfoAccess   = []int{1, 3, 6, 1, 5, 5, 7, 1, 1}
)

// oidNotInExtensions reports whether an extension with the given oid exists in
// extensions.
func oidInExtensions(oid asn1.ObjectIdentifier, extensions []pkix.Extension) bool {
	for _, e := range extensions {
		if e.Id.Equal(oid) {
			return true
		}
	}
	return false
}

// marshalSANs marshals a list of addresses into a the contents of an X.509
// SubjectAlternativeName extension.
func marshalSANs(dnsNames, emailAddresses []string, ipAddresses []net.IP, uris []*url.URL) (derBytes []byte, err error) {
	var rawValues []asn1.RawValue
	for _, name := range dnsNames {
		rawValues = append(rawValues, asn1.RawValue{Tag: nameTypeDNS, Class: 2, Bytes: []byte(name)})
	}
	for _, email := range emailAddresses {
		rawValues = append(rawValues, asn1.RawValue{Tag: nameTypeEmail, Class: 2, Bytes: []byte(email)})
	}
	for _, rawIP := range ipAddresses {
		// If possible, we always want to encode IPv4 addresses in 4 bytes.
		ip := rawIP.To4()
		if ip == nil {
			ip = rawIP
		}
		rawValues = append(rawValues, asn1.RawValue{Tag: nameTypeIP, Class: 2, Bytes: ip})
	}
	for _, uri := range uris {
		rawValues = append(rawValues, asn1.RawValue{Tag: nameTypeURI, Class: 2, Bytes: []byte(uri.String())})
	}
	return asn1.Marshal(rawValues)
}


type tbsCertificateRequest struct {
	Raw           asn1.RawContent
	Version       int
	Subject       asn1.RawValue
	PublicKey     publicKeyInfo
	RawAttributes []asn1.RawValue `asn1:"tag:0"`
}

type certificateRequest struct {
	Raw                asn1.RawContent
	TBSCSR             tbsCertificateRequest
	SignatureAlgorithm pkix.AlgorithmIdentifier
	SignatureValue     asn1.BitString
}

// oidExtensionRequest is a PKCS#9 OBJECT IDENTIFIER that indicates requested
// extensions in a CSR.
var oidExtensionRequest = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 9, 14}

// newRawAttributes converts AttributeTypeAndValueSETs from a template
// CertificateRequest's Attributes into tbsCertificateRequest RawAttributes.
func newRawAttributes(attributes []pkix.AttributeTypeAndValueSET) ([]asn1.RawValue, error) {
	var rawAttributes []asn1.RawValue
	b, err := asn1.Marshal(attributes)
	if err != nil {
		return nil, err
	}
	rest, err := asn1.Unmarshal(b, &rawAttributes)
	if err != nil {
		return nil, err
	}
	if len(rest) != 0 {
		return nil, errors.New("x509: failed to unmarshal raw CSR Attributes")
	}
	return rawAttributes, nil
}



func CreateTbsCsr(template *x509.CertificateRequest, public *ecdsa.PublicKey) (tbsCsr []byte, err error) {
	var publicKeyBytes []byte
	var publicKeyAlgorithm pkix.AlgorithmIdentifier
	publicKeyBytes, publicKeyAlgorithm, err = marshalPublicKey(public)
	if err != nil {
		return nil, err
	}

	var extensions []pkix.Extension

	if (len(template.DNSNames) > 0 || len(template.EmailAddresses) > 0 || len(template.IPAddresses) > 0 || len(template.URIs) > 0) &&
		!oidInExtensions(oidExtensionSubjectAltName, template.ExtraExtensions) {
		sanBytes, err := marshalSANs(template.DNSNames, template.EmailAddresses, template.IPAddresses, template.URIs)
		if err != nil {
			return nil, err
		}

		extensions = append(extensions, pkix.Extension{
			Id:    oidExtensionSubjectAltName,
			Value: sanBytes,
		})
	}

	extensions = append(extensions, template.ExtraExtensions...)

	// Make a copy of template.Attributes because we may alter it below.
	attributes := make([]pkix.AttributeTypeAndValueSET, 0, len(template.Attributes))
	for _, attr := range template.Attributes {
		values := make([][]pkix.AttributeTypeAndValue, len(attr.Value))
		copy(values, attr.Value)
		attributes = append(attributes, pkix.AttributeTypeAndValueSET{
			Type:  attr.Type,
			Value: values,
		})
	}

	extensionsAppended := false
	if len(extensions) > 0 {
		// Append the extensions to an existing attribute if possible.
		for _, atvSet := range attributes {
			if !atvSet.Type.Equal(oidExtensionRequest) || len(atvSet.Value) == 0 {
				continue
			}

			// specifiedExtensions contains all the extensions that we
			// found specified via template.Attributes.
			specifiedExtensions := make(map[string]bool)

			for _, atvs := range atvSet.Value {
				for _, atv := range atvs {
					specifiedExtensions[atv.Type.String()] = true
				}
			}

			newValue := make([]pkix.AttributeTypeAndValue, 0, len(atvSet.Value[0])+len(extensions))
			newValue = append(newValue, atvSet.Value[0]...)

			for _, e := range extensions {
				if specifiedExtensions[e.Id.String()] {
					// Attributes already contained a value for
					// this extension and it takes priority.
					continue
				}

				newValue = append(newValue, pkix.AttributeTypeAndValue{
					// There is no place for the critical
					// flag in an AttributeTypeAndValue.
					Type:  e.Id,
					Value: e.Value,
				})
			}

			atvSet.Value[0] = newValue
			extensionsAppended = true
			break
		}
	}

	rawAttributes, err := newRawAttributes(attributes)
	if err != nil {
		return
	}

	// If not included in attributes, add a new attribute for the
	// extensions.
	if len(extensions) > 0 && !extensionsAppended {
		attr := struct {
			Type  asn1.ObjectIdentifier
			Value [][]pkix.Extension `asn1:"set"`
		}{
			Type:  oidExtensionRequest,
			Value: [][]pkix.Extension{extensions},
		}

		b, err := asn1.Marshal(attr)
		if err != nil {
			return nil, errors.New("x509: failed to serialise extensions attribute: " + err.Error())
		}

		var rawValue asn1.RawValue
		if _, err := asn1.Unmarshal(b, &rawValue); err != nil {
			return nil, err
		}

		rawAttributes = append(rawAttributes, rawValue)
	}

	asn1Subject := template.RawSubject
	if len(asn1Subject) == 0 {
		asn1Subject, err = asn1.Marshal(template.Subject.ToRDNSequence())
		if err != nil {
			return nil, err
		}
	}

	tbsCSR := tbsCertificateRequest{
		Version: 0, // PKCS #10, RFC 2986
		Subject: asn1.RawValue{FullBytes: asn1Subject},
		PublicKey: publicKeyInfo{
			Algorithm: publicKeyAlgorithm,
			PublicKey: asn1.BitString{
				Bytes:     publicKeyBytes,
				BitLength: len(publicKeyBytes) * 8,
			},
		},
		RawAttributes: rawAttributes,
	}

	tbsCSRContents, err := asn1.Marshal(tbsCSR)
	if err != nil {
		return
	}

	return tbsCSRContents, nil
}

func CsrFromTbsAndSignature(tbsCsrContents []byte, signature []byte) ([]byte, error){
	tbsCSR := tbsCertificateRequest{}
	_, err := asn1.Unmarshal(tbsCsrContents, &tbsCSR)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshal TBSCSR failed")
	}

	tbsCSR.Raw = tbsCsrContents

	sigAlgo := pkix.AlgorithmIdentifier{Algorithm:oidSignatureECDSAWithSHA256}

	return asn1.Marshal(certificateRequest{
		TBSCSR:             tbsCSR,
		SignatureAlgorithm: sigAlgo,
		SignatureValue: asn1.BitString{
			Bytes:     signature,
			BitLength: len(signature) * 8,
		},
	})
}
