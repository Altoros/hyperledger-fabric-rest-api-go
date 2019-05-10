package handlers

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"net/http"
)

type ConfigResponse struct {
	Org string `json:"org"`
}

func GetConfigHandler(ec echo.Context) error {
	c := ec.(*ApiContext)

	configResponse := ConfigResponse{Org: c.Fsc().ApiConfig.Org.Name}

	jsonString, _ := json.Marshal(configResponse)

	ec.JSONBlob(http.StatusOK, []byte(jsonString))

	return nil
}
