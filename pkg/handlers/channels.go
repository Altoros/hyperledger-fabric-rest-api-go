package handlers

import (
	"errors"
	"fabric-rest-api-go/pkg/api"
	"github.com/labstack/echo/v4"
)

// Get channels list
func GetChannelsHandler(c echo.Context) error {
	return GetHandlerWrapper(c, api.FscInstance.Channels)
}

// Create channel
func PostChannelsHandler(c echo.Context) error {
	return GetJsonOutputWrapper(c, "", errors.New("not implemented"))
}

func GetChannelsChannelIdHandler(c echo.Context) error {
	jsonString, err := api.FscInstance.ChannelInfo(c.Param("channelId"))
	return GetJsonOutputWrapper(c, jsonString, err)
}

func GetChannelsChannelIdOrgsHandler(c echo.Context) error {
	jsonString, err := api.FscInstance.ChannelOrgs(c.Param("channelId"))
	return GetJsonOutputWrapper(c, jsonString, err)
}

func GetChannelsChannelIdPeersHandler(c echo.Context) error {
	jsonString, err := api.FscInstance.ChannelPeers(c.Param("channelId"))
	return GetJsonOutputWrapper(c, jsonString, err)
}
