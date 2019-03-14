package handlers

import (
	"fabric-rest-api-go/pkg/api"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

func GetQueryHandler(c echo.Context) error {
	fcn := c.FormValue("fcn")
	args := strings.Split(c.FormValue("args"), ",")

	if fcn == "" {
		return c.String(http.StatusBadRequest, "Fcn is required")
	}

	resultString, err := api.Query(&api.FscInstance, api.FscInstance.GetCurrentPeer(), c.Param("channelId"), c.Param("chaincodeId"), fcn, args)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSONBlob(http.StatusOK, []byte(fmt.Sprintf(`{"result": "%s"}`, resultString)))
}
