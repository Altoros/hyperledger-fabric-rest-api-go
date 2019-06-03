package handlers

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"fabric-rest-api-go/pkg/ca"
	"fabric-rest-api-go/pkg/utils"
	"fmt"
	"github.com/labstack/echo/v4"
	"math/big"
	"net/http"
)

func PostCaTbsCsrHandler(ec echo.Context) error {
	c := ec.(*ApiContext)

	caTbsCsrRequest := new(ca.CaTbsCsrRequest)
	if err := c.Bind(caTbsCsrRequest); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(caTbsCsrRequest); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, c.ValidationErrors(err).Error())
	}

	x := new(big.Int)
	x.SetString(caTbsCsrRequest.X, 16)
	y := new(big.Int)
	y.SetString(caTbsCsrRequest.Y, 16)

	publicKey := &ecdsa.PublicKey{X: x, Y: y, Curve: elliptic.P256()}

	tbsCsr, err := ca.TbsCsrBytes(caTbsCsrRequest, publicKey)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	hasher := sha256.New()
	hasher.Write(tbsCsr)
	tbsCsrSha256 := hasher.Sum(nil)

	// TODO change to actual signing material
	return ec.JSON(http.StatusOK, ca.CaTbsCsrResponse{TbsCsrBytes: utils.B64Encode(tbsCsr), TbsCsrHash: fmt.Sprintf("%x", tbsCsrSha256)})
}
