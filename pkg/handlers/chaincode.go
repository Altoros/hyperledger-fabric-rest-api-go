package handlers

import (
	"fabric-rest-api-go/pkg/api"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

func PostChaincodesInstallHandler(c echo.Context) error {
	chaincodeName := c.FormValue("name")
	chaincodeVersion := c.FormValue("version")
	channelId := c.FormValue("channel")

	ccHeader, err := c.FormFile("cc")
	if err != nil {
		return c.String(http.StatusUnprocessableEntity, "Problem with chaincode file upload: "+err.Error())
	}
	_ = ccHeader
	// TODO handle chaincode upload

	if chaincodeName == "" {
		return c.String(http.StatusBadRequest, "Chaincode name is required")
	}

	if chaincodeVersion == "" {
		return c.String(http.StatusBadRequest, "Chaincode version is required")
	}

	if channelId == "" {
		return c.String(http.StatusBadRequest, "Channel name is required")
	}

	if !api.CheckChannelExist(&api.FscInstance, api.FscInstance.GetCurrentPeer(), channelId) {
		return c.String(http.StatusInternalServerError, "Channel not exist")
	}

	resultString, err := api.ChaincodeInstall(&api.FscInstance, api.FscInstance.GetCurrentPeer(), channelId, chaincodeName, chaincodeVersion)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSONBlob(http.StatusOK, []byte( fmt.Sprintf(`{"result": "%s"}`, resultString)))
}

func PostChaincodesInstantiateHandler(c echo.Context) error {
	chaincodeName := c.FormValue("name")
	chaincodeVersion := c.FormValue("version")
	channelId := c.FormValue("channel")

	if chaincodeName == "" {
		return c.String(http.StatusBadRequest, "Chaincode name is required")
	}

	if chaincodeVersion == "" {
		return c.String(http.StatusBadRequest, "Chaincode version is required")
	}

	if channelId == "" {
		return c.String(http.StatusBadRequest, "Channel name is required")
	}

	if !api.CheckChannelExist(&api.FscInstance, api.FscInstance.GetCurrentPeer(), channelId) {
		return c.String(http.StatusInternalServerError, "Channel not exist")
	}

	resultString, err := api.ChaincodeInstantiate(&api.FscInstance, api.FscInstance.GetCurrentPeer(), channelId, chaincodeName, chaincodeVersion)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSONBlob(http.StatusOK, []byte( fmt.Sprintf(`{"result": "%s"}`, resultString)))
}

func GetChaincodesInstalledHandler(c echo.Context) error {
	return GetHandlerWrapper(c, api.FscInstance.InstalledChaincodes)
}

// Get instantiated chaincodes list
func GetChaincodesInstantiatedHandler(c echo.Context) error {
	jsonString, err := api.FscInstance.InstantiatedChaincodes(c.Param("channelId"))
	return GetJsonOutputWrapper(c, jsonString, err)
}

func GetChaincodesInfoHandler(c echo.Context) error {
	// TODO validate
	jsonString, err := api.FscInstance.ChaincodeInfo(c.Param("channelId"), c.Param("chaincodeId"))
	return GetJsonOutputWrapper(c, jsonString, err)
}
