package handlers

import (
	"fabric-rest-api-go/pkg/api"
	"fabric-rest-api-go/pkg/ca"
	"fabric-rest-api-go/pkg/context"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

func PostCaEnrollHandler(ec echo.Context) error {
	c := ec.(*context.ApiContext)

	enrollRequest := new(ca.ApiEnrollRequest)
	if err := c.Bind(enrollRequest); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if enrollRequest.Login == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Login is required")
	}

	if enrollRequest.Password == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Password is required")
	}

	resultString, err := api.CaEnroll(c.Fsc().ApiConfig, enrollRequest)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSONBlob(http.StatusOK, []byte( fmt.Sprintf(`{"result": "%s"}`, resultString)))
}
