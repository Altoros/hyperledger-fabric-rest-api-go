package handlers

import (
	"fabric-rest-api-go/pkg/api"
	"fabric-rest-api-go/pkg/context"
	"github.com/labstack/echo/v4"
	"net/http"
)

// Get channels list
func GetChannelsHandler(ec echo.Context) error {
	c := ec.(*context.ApiContext)

	peer, err := c.CurrentPeer()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	jsonString, err := api.Channels(c.Fsc(), peer)
	return GetJsonOutputWrapper(c, jsonString, err)
}

// Create channel
func PostChannelsHandler(ec echo.Context) error {
	// TODO implement
	return api.ChannelCreate("")
}

func GetChannelsChannelIdHandler(ec echo.Context) error {
	c := ec.(*context.ApiContext)
	jsonString, err := api.ChannelInfo(c.Fsc(), c.Param("channelId"))
	return GetJsonOutputWrapper(c, jsonString, err)
}

func GetChannelsChannelIdOrgsHandler(ec echo.Context) error {
	c := ec.(*context.ApiContext)
	jsonString, err := api.ChannelOrgs(c.Fsc(), c.Param("channelId"))
	return GetJsonOutputWrapper(c, jsonString, err)
}

func GetChannelsChannelIdPeersHandler(ec echo.Context) error {
	c := ec.(*context.ApiContext)
	jsonString, err := api.ChannelPeers(c.Fsc(), c.Param("channelId"))
	return GetJsonOutputWrapper(c, jsonString, err)
}
