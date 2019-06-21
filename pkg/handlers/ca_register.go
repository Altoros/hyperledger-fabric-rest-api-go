package handlers

import (
	"fabric-rest-api-go/pkg/api"
	"fabric-rest-api-go/pkg/ca"
	"fabric-rest-api-go/pkg/context"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

func PostCaRegisterHandler(ec echo.Context) error {
	c := ec.(*context.ApiContext)

	registerRequest := new(ca.ApiCaRegisterRequest)
	if err := c.Bind(registerRequest); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(registerRequest); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, c.ValidationErrors(err).Error())
	}

	resultString, err := api.CaRegister(c, c.Fsc().ApiConfig, registerRequest)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSONBlob(http.StatusOK, []byte( fmt.Sprintf(`{"result": "%s"}`, resultString)))
}
