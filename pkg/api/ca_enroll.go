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

func CaEnroll(apiConfig *sdk.Config, enrollRequest *ca.ApiEnrollRequest) (string, error) {

	// Private key generation
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return "", errors.Wrap(err, "failed to generate private key")
	}

	// Storing private key in PEM format
	privatePem, err := ca.PrivateKeyToPEM(privateKey)
	if err != nil {
		return "", errors.Wrap(err, "failed to create private key PEM")
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", errors.Wrap(err, "cannot obtain user home dir")
	}
	// TODO move this to config
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

	signCertPemB64, err := CaEnrollWithCsrPem(apiConfig, enrollRequest.Login, enrollRequest.Password, csrPEM)
	if err != nil {
		return "", errors.Wrap(err, "failed to enroll with CSR")
	}

	signCertPem, err := utils.B64Decode(signCertPemB64)
	if err != nil {
		return "", errors.Wrap(err, "failed to decode PEM from CA")
	}

	err = ioutil.WriteFile(restPath+"/admin_signcert.pem", signCertPem, 0644)
	if err != nil {
		return "", errors.Wrap(err, "key saving error")
	}

	return signCertPemB64, nil
}

func CaEnrollWithCsrPem(apiConfig *sdk.Config, login, password string, csrPEM string) (string, error) {
	enrollCaRequest := ca.CaEnrollRequest{
		CertificateRequest: csrPEM,
	}

	jsonEnrollRequest, err := json.Marshal(enrollCaRequest)
	if err != nil {
		return "", err
	}

	caEnrollUrl := fmt.Sprintf("%s/enroll", apiConfig.Ca.Address)

	// sending CSR to CA /enroll
	req, err := http.NewRequest("POST", caEnrollUrl, bytes.NewBuffer(jsonEnrollRequest))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(login, password)

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

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Wrap(err, "failed to read response")
	}

	caEnrollResponse := ca.CaEnrollResponse{}
	err = json.Unmarshal(body, &caEnrollResponse)
	if err != nil {
		return "", errors.Wrap(err, "CA response unmarshal error")
	}

	if !caEnrollResponse.Success {
		return "", errors.Errorf("CA response with errors: %s", strings.Join(caEnrollResponse.Errors, ", "))
	}

	return caEnrollResponse.Result.Cert, nil
}
