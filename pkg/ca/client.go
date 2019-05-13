package ca

import (
	"crypto/tls"
	"crypto/x509"
	"fabric-rest-api-go/pkg/sdk"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
)

func HttpClient(apiConfig *sdk.Config) (*http.Client, error) {
	client := &http.Client{}

	if apiConfig.Ca.Tls {

		// Read in the cert file
		certs, err := ioutil.ReadFile(apiConfig.Ca.TlsCertFile)
		if err != nil {
			return nil, errors.Errorf("failed to load CA TLS certificate")
		}

		// Add cert to the certificates pool
		rootCAs := x509.NewCertPool()
		if ok := rootCAs.AppendCertsFromPEM(certs); !ok {
			return nil, errors.Errorf("failed to add CA TLS certificate")
		}

		// Trust the augmented cert pool in our client
		config := &tls.Config{
			RootCAs:            rootCAs,
			InsecureSkipVerify: true, // TODO get rid of
		}
		tr := &http.Transport{TLSClientConfig: config}
		client = &http.Client{Transport: tr}
	}

	return client, nil
}
