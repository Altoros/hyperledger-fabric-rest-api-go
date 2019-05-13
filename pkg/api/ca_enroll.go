package api

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/json"
	"fabric-rest-api-go/pkg/ca"
	"fabric-rest-api-go/pkg/sdk"
	"fabric-rest-api-go/pkg/utils"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func CaEnroll(apiConfig *sdk.Config,  enrollRequest *ca.ApiEnrollRequest) (string, error) {

	// Private key generation
	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	// Storing private key in PEM format
	privatePem, _ := ca.PrivateKeyToPEM(privateKey)

	home, err := os.UserHomeDir()
	if err != nil {
		return "", errors.Wrap(err, "cannot obtain user home dir")
	}
	restPath := home + "/.fabric-rest-api-go/"

	err = os.MkdirAll(restPath, os.ModePerm)
	if err != nil {
		return "", errors.Wrap(err, "key saving error")
	}
	err = ioutil.WriteFile(restPath+"/admin_private.pem", privatePem, 0644)
	if err != nil {
		return "", errors.Wrap(err, "key saving error")
	}

	csrPEM, err := ca.CsrPEM(privateKey, enrollRequest)
	if err != nil {
		return "", err
	}

	// send CSR to CA /enroll
	enrollCaRequest := ca.CaEnrollRequest{
		CertificateRequest: csrPEM,
	}

	jsonEnrollRequest, err := json.Marshal(enrollCaRequest)
	if err != nil {
		return "", err
	}

	caEnrollUrl := fmt.Sprintf("%s/enroll", apiConfig.Ca.Address)

	req, err := http.NewRequest("POST", caEnrollUrl, bytes.NewBuffer(jsonEnrollRequest))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(enrollRequest.Login, enrollRequest.Password)

	client, err := ca.HttpClient(apiConfig)
	if err != nil {
		return "", errors.Wrap(err, "failed to create HTTP client")
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "failed to connect to CA")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return "", errors.Errorf("CA returned status %s, but 201 expected", resp.Status)
	}

	body, _ := ioutil.ReadAll(resp.Body)

	caEnrollResponse := ca.CaEnrollResponse{}
	err = json.Unmarshal(body, &caEnrollResponse)
	if err != nil {
		return "", errors.Wrap(err, "CA response unmarshal error")
	}

	if !caEnrollResponse.Success {
		return "", errors.Errorf("CA response with errors: %s", strings.Join(caEnrollResponse.Errors, ", "))
	}

	signCertPem, _ := utils.B64Decode(caEnrollResponse.Result.Cert)

	err = ioutil.WriteFile(restPath+"/admin_signcert.pem", signCertPem, 0644)
	if err != nil {
		return "", errors.Wrap(err, "key saving error")
	}

	return caEnrollResponse.Result.Cert, nil
}
