package handlers

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func GetChaincodesInstalledHandler(ec echo.Context) error {
	c := ec.(*ApiContext)

	peer, err := c.CurrentPeer()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	jsonString, err := c.Fsc().InstalledChaincodes(peer)
	return GetJsonOutputWrapper(c, jsonString, err)
}

// Get instantiated chaincodes list
func GetChaincodesInstantiatedHandler(ec echo.Context) error {
	c := ec.(*ApiContext)

	peer, err := c.CurrentPeer()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	jsonString, err := c.Fsc().InstantiatedChaincodes(c.Param("channelId"), peer)
	return GetJsonOutputWrapper(c, jsonString, err)
}

func GetChaincodesInfoHandler(ec echo.Context) error {
	c := ec.(*ApiContext)

	peer, err := c.CurrentPeer()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// TODO validate
	jsonString, err := c.Fsc().ChaincodeInfo(c.Param("channelId"), c.Param("chaincodeId"), peer)
	return GetJsonOutputWrapper(c, jsonString, err)
}
