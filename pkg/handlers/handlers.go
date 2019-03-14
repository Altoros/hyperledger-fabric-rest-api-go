package handlers

import (
	"fabric-rest-api-go/pkg/api"
	"github.com/labstack/echo/v4"
	"net/http"
)

func WelcomeHandler(c echo.Context) error {
	return c.String(http.StatusOK, "This is a Fabric REST Api welcome page.")
}

// A very simple health check
func HealthCheckHandler(c echo.Context) error {
	return c.JSONBlob(http.StatusOK, []byte(`{"alive": true}`))
}

// TODO remove, replace with chaincode install & instantiate calls
// Create test channel, install and instantiate test chaincode
func InitTestFixturesHandler(c echo.Context) error {
	err := api.FscInstance.InitBasicTestFixturesHandler()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSONBlob(http.StatusOK, []byte(`{"message": "Test channel created, chaincode installed and instantiated"}`))
}
