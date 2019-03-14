package handlers

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

// Generic wrapper for GET handlers
func GetHandlerWrapper(c echo.Context, handlerFunc func() (string, error)) error {
	jsonString, err := handlerFunc()
	return GetJsonOutputWrapper(c, jsonString, err)
}

func GetJsonOutputWrapper(c echo.Context, jsonString string, err error) error {
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSONBlob(http.StatusOK, []byte(jsonString))
}

