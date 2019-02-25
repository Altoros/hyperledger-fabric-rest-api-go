package main

import (
	"fabric-rest-api-go/api"
	"fmt"
	"log"
	"net/http"

	"os"
	"github.com/alexdnn11/fabric-rest-go/api"
	"github.com/gorilla/mux"
	"encoding/json"
)

type Config struct {
	Org struct {
		Admin string `json:"admin"`
		Name  string `json:"name"`
	} `json:"org"`
	User struct {
		Name string `json:"name"`
	} `json:"user"`
	Gopath        string `json:"gopath"`
	ChaincodePath string `json:"chaincodePath"`
	ConfigFile    string `json:"configPath"`
}

func LoadConfiguration(file string) Config {
	var config Config
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)
	return config
}

func main() {

	config := LoadConfiguration("./config.json")

	api.FscInstance = api.FabricSdkClient{
		// Chaincode parameters
		GOPATH:        os.Getenv(config.Gopath),
		ChaincodePath: config.ChaincodePath,

		// Org parameters
		OrgAdmin: config.Org.Admin,
		OrgName:  config.Org.Name,

		ConfigFile: config.ConfigFile,

		// User parameters
		UserName: config.User.Name,
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
