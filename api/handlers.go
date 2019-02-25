package api

import (
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"strings"
)

func WelcomeHandler(w http.ResponseWriter, _ *http.Request) {
	_, err := io.WriteString(w, "This is a Fabric REST Api welcome page.")
	if err != nil {
		panic(err)
	}
}

// A very simple health check
func HealthCheckHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, err := io.WriteString(w, `{"alive": true}`)
	if err != nil {
		panic(err)
	}
}

// Generic wrapper for GET handlers
func GetHandlerWrapper(w http.ResponseWriter, _ *http.Request, handlerFunc func() (string, error)) {
	jsonString, err := handlerFunc()
	GetJsonOutputWrapper(w, jsonString, err)
}

func GetJsonOutputWrapper(w http.ResponseWriter, jsonString string, err error) {
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, err = io.WriteString(w, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, err = io.WriteString(w, jsonString)
	if err != nil {
		panic(err)
	}
}

func GetChaincodesInstalledHandler(w http.ResponseWriter, r *http.Request) {
	GetHandlerWrapper(w, r, FscInstance.InstalledChaincodes)
}

// Get instantiated chaincodes list
func GetChaincodesInstantiatedHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	jsonString, err := FscInstance.InstantiatedChaincodes(vars["channelId"])
	GetJsonOutputWrapper(w, jsonString, err)
}

func GetChaincodesInfoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// TODO validate
	jsonString, err := FscInstance.ChaincodeInfo(vars["channelId"], vars["chaincodeId"])
	GetJsonOutputWrapper(w, jsonString, err)
}

// Get channels list
func GetChannelsHandler(w http.ResponseWriter, r *http.Request) {
	GetHandlerWrapper(w, r, FscInstance.Channels)
}

// Create channel
func PostChannelsHandler(w http.ResponseWriter, r *http.Request) {
	panic("not implemented")
}

func GetChannelsChannelIdHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	jsonString, err := FscInstance.ChannelInfo(vars["channelId"])
	GetJsonOutputWrapper(w, jsonString, err)
}

func GetChannelsChannelIdOrgsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	jsonString, err := FscInstance.ChannelOrgs(vars["channelId"])
	GetJsonOutputWrapper(w, jsonString, err)
}

func GetChannelsChannelIdPeersHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	jsonString, err := FscInstance.ChannelPeers(vars["channelId"])
	GetJsonOutputWrapper(w, jsonString, err)
}

func GetQueryHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		io.WriteString(w, err.Error())
		return
	}

	fcn := r.FormValue("fcn")
	args := strings.Split(r.FormValue("args"), ",")

	if fcn == "" {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "Fcn is required")
		return
	}

	resultString, err := Query(&FscInstance, FscInstance.GetCurrentPeer(), vars["channelId"], vars["chaincodeId"], fcn, args)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, err = io.WriteString(w, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, err = io.WriteString(w, fmt.Sprintf(`{"result": "%s"}`, resultString))
	if err != nil {
		panic(err)
	}
}

func PostInvokeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		io.WriteString(w, err.Error())
		return
	}

	fcn := r.FormValue("fcn")
	args := strings.Split(r.FormValue("args"), ",")

	if fcn == "" {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "Fcn is required")
		return
	}

	resultString, err := Invoke(&FscInstance, vars["channelId"], vars["chaincodeId"], fcn, args)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, err = io.WriteString(w, fmt.Sprintf(`{"result": "%s"}`, resultString))
	if err != nil {
		panic(err)
	}
}

// Create test channel, install and instantiate test chaincode
func InitTestFixturesHandler(w http.ResponseWriter, _ *http.Request) {
	err := FscInstance.InitTestFixturesHandler()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, err.Error())
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, `{"message": "Test channel created, chaincode installed and instantiated"}`)
	}
}
