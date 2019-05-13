package ca

import (
	"fmt"
	"strconv"
	"strings"
)

type ApiEnrollRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Email    string `json:"email"`
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

type ApiRegisterRequest struct {
	Login string `json:"login"`
}

type CaRegisterResponseResult struct {
	Secret string `json:"secret"`
}

type CaRegisterResponseError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type CaRegisterResponse struct {
	Success bool `json:"success"`
}

type CaRegisterResponseWithResult struct {
	Result CaRegisterResponseResult `json:"result"`
	Errors   []CaRegisterResponseError `json:"errors"`
	Messages []string                  `json:"messages"`
	Success bool `json:"success"`
}

type CaRegisterResponseWithErrors struct {
	Errors   []CaRegisterResponseError `json:"errors"`
	Messages []string                  `json:"messages"`
}

func (caRegisterResponseWithErrors *CaRegisterResponseWithErrors) ErrorsString() string {
	var messages []string
	for _, caRegisterResponseError := range caRegisterResponseWithErrors.Errors {
		messages = append(messages, fmt.Sprintf("%s (code %s)", caRegisterResponseError.Message, strconv.Itoa(caRegisterResponseError.Code)))
	}
	return strings.Join(messages, ", ")
}
