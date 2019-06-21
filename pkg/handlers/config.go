package handlers

import (
	"encoding/json"
	"fabric-rest-api-go/pkg/context"
	"github.com/labstack/echo/v4"
	"net/http"
)

type ConfigResponse struct {
	Org string `json:"org"`
}

func GetConfigHandler(ec echo.Context) error {
	c := ec.(*context.ApiContext)

	configResponse := ConfigResponse{Org: c.Fsc().ApiConfig.Org.Name}

	jsonString, _ := json.Marshal(configResponse)

	ec.JSONBlob(http.StatusOK, []byte(jsonString))

	return nil
}
