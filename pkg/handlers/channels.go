package handlers

import (
	"errors"
	"fabric-rest-api-go/pkg/api"
	"github.com/gorilla/mux"
	"net/http"
)

// Get channels list
func GetChannelsHandler(w http.ResponseWriter, r *http.Request) {
	GetHandlerWrapper(w, r, api.FscInstance.Channels)
}

// Create channel
func PostChannelsHandler(w http.ResponseWriter, r *http.Request) {
	GetJsonOutputWrapper(w, "", errors.New("not implemented"))
}

func GetChannelsChannelIdHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	jsonString, err := api.FscInstance.ChannelInfo(vars["channelId"])
	GetJsonOutputWrapper(w, jsonString, err)
}

func GetChannelsChannelIdOrgsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	jsonString, err := api.FscInstance.ChannelOrgs(vars["channelId"])
	GetJsonOutputWrapper(w, jsonString, err)
}

func GetChannelsChannelIdPeersHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	jsonString, err := api.FscInstance.ChannelPeers(vars["channelId"])
	GetJsonOutputWrapper(w, jsonString, err)
}
