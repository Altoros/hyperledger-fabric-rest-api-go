package api

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fabric-rest-api-go/pkg/ca"
	"fabric-rest-api-go/pkg/context"
	"fabric-rest-api-go/pkg/sdk"
	"fabric-rest-api-go/pkg/utils"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"os"
)

type CaRegisterAttr struct {
	// Attribute name
	Name string `json:"name" validate:"required"`
	// Value of attribute
	Value string `json:"value" validate:"required"`
	// A value of true indicates that this attribute should be included in an enrollment certificate by default
	Ecert bool `json:"ecert"`
}

type CaRegisterRequest struct {
	// The enrollment ID which uniquely identifies an identity
	Id string `json:"id" validate:"required"`
	// The type of the identity (e.g. user, app, peer, orderer, etc)
	Type string `json:"type,omitempty"`
	// The enrollment secret. If not provided, a random secret is generated.
	Secret string `json:"secret,omitempty"`
	// The maximum number of times that the secret can be used to enroll.
	// If 0, use the configured max_enrollments of the fabric-ca-server;
	// If > 0 and <= configured max enrollments of the fabric-ca-server, use max_enrollments;
	// If > configured max enrollments of the fabric-ca-server, error.
	MaxEnrollments int `json:"max_enrollments,omitempty"`
	// The affiliation of the new identity.
	// If no affliation is provided, the affiliation of the registrar is used.
	Affiliation string `json:"affiliation,omitempty"`
	// An array of attribute names and values to give to the registered identity.
	Attrs []CaRegisterAttr `json:"attrs,omitempty"`
	// Name of the CA to direct traffic to within server.
	Caname string `json:"caname"`
}

func CaRegister(c context.ValidationProvider, apiConfig *sdk.Config, registerRequest *ca.ApiCaRegisterRequest) (string, error) {

	caRegisterRequest := CaRegisterRequest{
		Id:     registerRequest.Login,
		Secret: registerRequest.Password,
		Type:   "client",
	}
	if err := c.Validate(caRegisterRequest); err != nil {
		return "", errors.Wrap(c.ValidationErrors(err), "failed to construct correct CaRegisterRequest")
	}

	jsonRegisterRequest, err := json.Marshal(caRegisterRequest)
	if err != nil {
		return "", errors.Wrap(err, "failed to marshal caRegisterRequest json")
	}

	// load private key pem
	home, err := os.UserHomeDir()
	if err != nil {
		return "", errors.Wrap(err, "cannot obtain user home dir")
	}
	restPath := home + "/.fabric-rest-api-go/"

	privatePem, err := ioutil.ReadFile(restPath + "/admin_private.pem")
	if err != nil {
		return "", errors.Wrap(err, "key reading error")
	}

	// load signcert pem
	signCertPem, err := ioutil.ReadFile(restPath + "/admin_signcert.pem")
	if err != nil {
		return "", errors.Wrap(err, "signcert reading error")
	}

	// generate body
	requestUri := "/register"
	requestMethod := "POST"
	requestBody := []byte (jsonRegisterRequest)

	// generate payload

	b64body := utils.B64Encode(requestBody)
	b64signCert := utils.B64Encode(signCertPem)
	b64uri := utils.B64Encode([]byte(requestUri))
	payload := requestMethod + "." + b64uri + "." + b64body + "." + b64signCert

	hasher := sha256.New()
	hasher.Write([]byte(payload))
	payloadSha256 := hasher.Sum(nil)

	// decode private key from PEM
	privateKeyEC, err := ca.PEMtoPrivateKey(privatePem)
	if err != nil {
		return "", errors.WithMessage(err, "private key conversion failure")
	}

	// sign payload hash
	ecSignature, err := ca.SignECDSA(privateKeyEC, payloadSha256)
	if err != nil {
		return "", errors.WithMessage(err, "signature generation failure")
	}
	if len(ecSignature) == 0 {
		return "", errors.New("signature creation failed, must be different than nil")
	}

	b64sig := utils.B64Encode(ecSignature)
	// Authorization token
	token := b64signCert + "." + b64sig

	caRegisterUrl := fmt.Sprintf("%s/register", apiConfig.Ca.Address)

	req, err := http.NewRequest("POST", caRegisterUrl, bytes.NewBuffer(jsonRegisterRequest))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)

	client, err := ca.HttpClient(apiConfig)
	if err != nil {
		return "", errors.Wrap(err, "failed to create HTTP client")
	}

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Wrap(err, "failed to read response")
	}

	caRegisterResponse := ca.CaRegisterResponse{}
	err = json.Unmarshal(body, &caRegisterResponse)
	if err != nil {
		return "", errors.Wrap(err, "CA response unmarshal error")
	}

	if !caRegisterResponse.Success {
		caRegisterResponseWithErrors := ca.CaRegisterResponseWithErrors{}
		err = json.Unmarshal(body, &caRegisterResponseWithErrors)
		if err != nil {
			return "", errors.Wrap(err, "CA response errors unmarshal error")
		}

		return "", errors.Errorf("CA response with errors: %s", caRegisterResponseWithErrors.ErrorsString())
	}

	caRegisterResponseWithResult := ca.CaRegisterResponseWithResult{}
	err = json.Unmarshal(body, &caRegisterResponseWithResult)
	if err != nil {
		return "", errors.Wrap(err, "CA response result unmarshal error")
	}

	return caRegisterResponseWithResult.Result.Secret, nil
}
