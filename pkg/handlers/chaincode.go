package handlers

import (
	"fabric-rest-api-go/pkg/api"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"net/http"
)

func PostChaincodesInstallHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		io.WriteString(w, err.Error())
		return
	}

	chaincodeName := r.FormValue("name")
	chaincodeVersion := r.FormValue("version")
	channelId := r.FormValue("channel")

	ccFile, ccHeader, err := r.FormFile("cc")
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		io.WriteString(w, "Problem with chaincode file upload: "+err.Error())
		return
	}
	_ = ccFile
	_ = ccHeader
	// TODO handle chaincode upload


	if chaincodeName == "" {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "Chaincode name is required")
		return
	}

	if chaincodeVersion == "" {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "Chaincode version is required")
		return
	}

	if channelId == "" {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "Channel version is required")
		return
	}

	if !api.CheckChannelExist(&api.FscInstance, api.FscInstance.GetCurrentPeer(), channelId) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, "Channel not exist")
		return
	}

	resultString, err := api.ChaincodeInstall(&api.FscInstance, api.FscInstance.GetCurrentPeer(), channelId, chaincodeName, chaincodeVersion)
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

func PostChaincodesInstantiateHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		io.WriteString(w, err.Error())
		return
	}

	chaincodeName := r.FormValue("name")
	chaincodeVersion := r.FormValue("version")
	channelId := r.FormValue("channel")

	if chaincodeName == "" {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "Chaincode name is required")
		return
	}

	if chaincodeVersion == "" {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "Chaincode version is required")
		return
	}

	if channelId == "" {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "Channel is required")
		return
	}

	if !api.CheckChannelExist(&api.FscInstance, api.FscInstance.GetCurrentPeer(), channelId) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, "Channel not exist")
		return
	}

	resultString, err := api.ChaincodeInstantiate(&api.FscInstance, api.FscInstance.GetCurrentPeer(), channelId, chaincodeName, chaincodeVersion)
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

func GetChaincodesInstalledHandler(w http.ResponseWriter, r *http.Request) {
	GetHandlerWrapper(w, r, api.FscInstance.InstalledChaincodes)
}

// Get instantiated chaincodes list
func GetChaincodesInstantiatedHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	jsonString, err := api.FscInstance.InstantiatedChaincodes(vars["channelId"])
	GetJsonOutputWrapper(w, jsonString, err)
}

func GetChaincodesInfoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// TODO validate
	jsonString, err := api.FscInstance.ChaincodeInfo(vars["channelId"], vars["chaincodeId"])
	GetJsonOutputWrapper(w, jsonString, err)
}
