package main

import (
	"fabric-rest-api-go/api"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	api.FscInstance = api.FabricSdkClient{
		// Org parameters
		OrgAdmin: "Admin",
		OrgName:  "org1",

		ConfigFile: "test/config.yaml",

		// User parameters
		UserName: "User1",
	}

	err := api.FscInstance.Initialize()
	if err != nil {
		panic(err)
	}

	fmt.Println("Start listening to localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", Router()))
}

func Router() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", api.WelcomeHandler)
	r.HandleFunc("/health", api.HealthCheckHandler)
	r.HandleFunc("/chaincodes/installed", api.GetChaincodesInstalledHandler).Methods("GET")
	r.HandleFunc("/channels/{channelId}/chaincodes/instantiated", api.GetChaincodesInstantiatedHandler).Methods("GET") // TODO
	r.HandleFunc("/channels/{channelId}/chaincodes/{chaincodeId}/info", api.GetChaincodesInfoHandler).Methods("GET")
	r.HandleFunc("/channels", api.GetChannelsHandler).Methods("GET")
	r.HandleFunc("/channels", api.PostChannelsHandler).Methods("POST") // TODO

	r.HandleFunc("/channels/{channelId}", api.GetChannelsChannelIdHandler).Methods("GET")
	r.HandleFunc("/channels/{channelId}/orgs", api.GetChannelsChannelIdOrgsHandler).Methods("GET") // TODO
	r.HandleFunc("/channels/{channelId}/peers", api.GetChannelsChannelIdPeersHandler).Methods("GET")

	r.HandleFunc("/channels/{channelId}/chaincodes/{chaincodeId}/query", api.GetQueryHandler).Methods("GET")
	r.HandleFunc("/channels/{channelId}/chaincodes/{chaincodeId}/invoke", api.PostInvokeHandler).Methods("POST")

	r.HandleFunc("/init_test_fixtures", api.InitTestFixturesHandler).Methods("POST") // for test purposes

	return r
}
