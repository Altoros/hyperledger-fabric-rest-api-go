package handlers

import (
	"fabric-rest-api-go/pkg/api"
	"fabric-rest-api-go/pkg/ca"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

func PostCaRegisterHandler(ec echo.Context) error {
	c := ec.(*ApiContext)

	registerRequest := new(ca.ApiRegisterRequest)
	if err := c.Bind(registerRequest); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if registerRequest.Login == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "New user login is required")
	}

	resultString, err := api.CaRegister(c.Fsc().ApiConfig, registerRequest)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSONBlob(http.StatusOK, []byte( fmt.Sprintf(`{"result": "%s"}`, resultString)))
}
