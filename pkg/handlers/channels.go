package handlers

import (
	"errors"
	"github.com/labstack/echo/v4"
	"net/http"
)

// Get channels list
func GetChannelsHandler(ec echo.Context) error {
	c := ec.(*ApiContext)

	peer, err := c.CurrentPeer()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	jsonString, err := c.Fsc().Channels(peer)
	return GetJsonOutputWrapper(c, jsonString, err)
}

// Create channel
func PostChannelsHandler(ec echo.Context) error {
	c := ec.(*ApiContext)
	return GetJsonOutputWrapper(c, "", errors.New("not implemented"))
}

func GetChannelsChannelIdHandler(ec echo.Context) error {
	c := ec.(*ApiContext)
	jsonString, err := c.Fsc().ChannelInfo(c.Param("channelId"))
	return GetJsonOutputWrapper(c, jsonString, err)
}

func GetChannelsChannelIdOrgsHandler(ec echo.Context) error {
	c := ec.(*ApiContext)
	jsonString, err := c.Fsc().ChannelOrgs(c.Param("channelId"))
	return GetJsonOutputWrapper(c, jsonString, err)
}

func GetChannelsChannelIdPeersHandler(ec echo.Context) error {
	c := ec.(*ApiContext)
	jsonString, err := c.Fsc().ChannelPeers(c.Param("channelId"))
	return GetJsonOutputWrapper(c, jsonString, err)
}
