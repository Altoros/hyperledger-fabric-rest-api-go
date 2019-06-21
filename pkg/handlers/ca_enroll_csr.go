package handlers

import (
	"bytes"
	"encoding/pem"
	"fabric-rest-api-go/pkg/api"
	"fabric-rest-api-go/pkg/ca"
	"fabric-rest-api-go/pkg/context"
	"fabric-rest-api-go/pkg/utils"
	"github.com/labstack/echo/v4"
	"math/big"
	"net/http"
)

type CaEnrollCsrRequest struct {
	Login       string `json:"login" validate:"required"`
	Password    string `json:"password" validate:"required"`
	TbsCsrBytes string `json:"tbs_csr_bytes" validate:"required"`
	R           string `json:"r" validate:"required"`
	S           string `json:"s" validate:"required"`
}

type CaEnrollCsrResponse struct {
	UserCert string `json:"user_cert"`
}

func PostCaEnrollCsrHandler(ec echo.Context) error {
	c := ec.(*context.ApiContext)

	caEnrollCsrRequest := new(CaEnrollCsrRequest)
	if err := c.Bind(caEnrollCsrRequest); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(caEnrollCsrRequest); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, c.ValidationErrors(err).Error())
	}

	r := new(big.Int)
	r.SetString(caEnrollCsrRequest.R, 16)
	s := new(big.Int)
	s.SetString(caEnrollCsrRequest.S, 16)

	//signature, err := ca.MarshalECDSASignature(r, ca.ToLowS_P256(s))
	signature, err := ca.MarshalECDSASignature(r, s)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// create tbsCsr from b64
	tbsCsrBytes, err := utils.B64Decode(caEnrollCsrRequest.TbsCsrBytes)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	csrBytes, err := ca.CsrFromTbsAndSignature(tbsCsrBytes, signature)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// TODO refactor
	buf := new(bytes.Buffer)

	err = pem.Encode(buf, &pem.Block{Type: "CERTIFICATE REQUEST", Bytes: csrBytes})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	csrPEM := buf.String()

	cert, err := api.CaEnrollWithCsrPem(c.Fsc().ApiConfig, caEnrollCsrRequest.Login, caEnrollCsrRequest.Password, csrPEM)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ec.JSON(http.StatusOK, CaEnrollCsrResponse{UserCert: cert})
}
