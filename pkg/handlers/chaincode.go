package handlers

import (
	"fabric-rest-api-go/pkg/api"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

func PostChaincodesInstallHandler(ec echo.Context) error {
	c := ec.(*ApiContext)

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

	peer, err := c.CurrentPeer()
	if err != nil {
		return err
	}

	if !api.CheckChannelExist(c.Fsc(), peer, channelId) {
		return c.String(http.StatusInternalServerError, "Channel not exist")
	}

	resultString, err := api.ChaincodeInstall(c.Fsc(), peer, channelId, chaincodeName, chaincodeVersion)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSONBlob(http.StatusOK, []byte( fmt.Sprintf(`{"result": "%s"}`, resultString)))
}

func PostChaincodesInstantiateHandler(ec echo.Context) error {
	c := ec.(*ApiContext)

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

	peer, err := c.CurrentPeer()
	if err != nil {
		return err
	}

	if !api.CheckChannelExist(c.Fsc(), peer, channelId) {
		return c.String(http.StatusInternalServerError, "Channel not exist")
	}

	resultString, err := api.ChaincodeInstantiate(c.Fsc(), peer, channelId, chaincodeName, chaincodeVersion)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSONBlob(http.StatusOK, []byte( fmt.Sprintf(`{"result": "%s"}`, resultString)))
}

func GetChaincodesInstalledHandler(ec echo.Context) error {
	c := ec.(*ApiContext)

	peer, err := c.CurrentPeer()
	if err != nil {
		return err
	}

	jsonString, err := c.Fsc().InstalledChaincodes(peer)
	return GetJsonOutputWrapper(c, jsonString, err)
}

// Get instantiated chaincodes list
func GetChaincodesInstantiatedHandler(ec echo.Context) error {
	c := ec.(*ApiContext)
	jsonString, err := c.Fsc().InstantiatedChaincodes(c.Param("channelId"))
	return GetJsonOutputWrapper(c, jsonString, err)
}

func GetChaincodesInfoHandler(ec echo.Context) error {
	c := ec.(*ApiContext)
	// TODO validate
	jsonString, err := c.Fsc().ChaincodeInfo(c.Param("channelId"), c.Param("chaincodeId"))
	return GetJsonOutputWrapper(c, jsonString, err)
}
