package handlers

import (
	"fabric-rest-api-go/pkg/api"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

type invokeRequest struct {
	Fcn string `json:"fcn"`
	Args []string `json:"args"`
}

func PostInvokeHandler(ec echo.Context) error {
	c := ec.(*ApiContext)

	m := new(invokeRequest)
	if err := c.Bind(m); err != nil {
		return err
	}

	fcn := m.Fcn
	args := m.Args

	if fcn == "" {
		return c.String(http.StatusBadRequest, "Fcn is required")
	}

	resultString, err := api.Invoke(c.Fsc(), c.Param("channelId"), c.Param("chaincodeId"), fcn, args)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSONBlob(http.StatusOK, []byte(fmt.Sprintf(`{"result": "%s"}`, resultString)))
}
