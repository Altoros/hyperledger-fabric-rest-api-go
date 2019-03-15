package handlers

import (
	"fabric-rest-api-go/pkg/api"
	"github.com/Jeffail/gabs"
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

	resultJsonObj := gabs.New()

	// Note: if result is a number, then .ParseJSON will detect it, and .String() it without quotes

	jsonParsed, err := gabs.ParseJSON([]byte(resultString))
	if err == nil {
		// if query result is parsed as JSON, it will be used as one
		resultJsonObj.Set(jsonParsed.Data(), "result")
	} else {
		// in other cases it will be used as a string
		resultJsonObj.Set(resultString, "result")
	}

	return c.JSONBlob(http.StatusOK, []byte(resultJsonObj.String()))
}
