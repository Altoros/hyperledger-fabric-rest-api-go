package handlers

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func GetJsonOutputWrapper(c echo.Context, jsonString string, err error) error {
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSONBlob(http.StatusOK, []byte(jsonString))
}

