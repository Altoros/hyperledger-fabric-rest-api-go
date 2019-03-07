package handlers

import (
	"fabric-rest-api-go/pkg/api"
	"io"
	"net/http"
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

// TODO remove, replace with chaincode install & instantiate calls
// Create test channel, install and instantiate test chaincode
func InitTestFixturesHandler(w http.ResponseWriter, _ *http.Request) {
	err := api.FscInstance.InitBasicTestFixturesHandler()
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
