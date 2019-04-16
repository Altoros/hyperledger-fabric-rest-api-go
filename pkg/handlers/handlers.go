package handlers

import (
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
