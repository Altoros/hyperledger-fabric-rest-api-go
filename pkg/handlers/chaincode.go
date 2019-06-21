package handlers

import (
	"fabric-rest-api-go/pkg/api"
	"fabric-rest-api-go/pkg/context"
	"github.com/labstack/echo/v4"
	"net/http"
)

func GetChaincodesInstalledHandler(ec echo.Context) error {
	c := ec.(*context.ApiContext)

	peer, err := c.CurrentPeer()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	jsonString, err := api.InstalledChaincodes(c.Fsc(), peer)
	return GetJsonOutputWrapper(c, jsonString, err)
}

// Get instantiated chaincodes list
func GetChaincodesInstantiatedHandler(ec echo.Context) error {
	c := ec.(*context.ApiContext)

	peer, err := c.CurrentPeer()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	jsonString, err := api.InstantiatedChaincodes(c.Fsc(),c.Param("channelId"), peer)
	return GetJsonOutputWrapper(c, jsonString, err)
}

func GetChaincodesInfoHandler(ec echo.Context) error {
	c := ec.(*context.ApiContext)

	peer, err := c.CurrentPeer()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// TODO validate
	jsonString, err := api.ChaincodeInfo(c.Fsc(), c.Param("channelId"), c.Param("chaincodeId"), peer)
	return GetJsonOutputWrapper(c, jsonString, err)
}
