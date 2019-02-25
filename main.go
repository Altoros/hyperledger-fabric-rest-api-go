package main

import (
	"fabric-rest-api-go/api"
	"fmt"
	"github.com/pkg/errors"
	"log"
	"net/http"

	"encoding/json"
	"github.com/gorilla/mux"
	"os"
)

type ApiConfig struct {
	Org struct {
		Admin string `json:"admin"`
		Name  string `json:"name"`
	} `json:"org"`
	User struct {
		Name string `json:"name"`
	} `json:"user"`
	ConfigPath string `json:"configPath"`
}

func LoadConfiguration(file string) (*ApiConfig, error) {
	var config *ApiConfig
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		return nil, errors.WithMessage(err, "Unable to open configuration file")
	}
	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&config)
	if err != nil {
		return nil, errors.WithMessage(err, "Unable to parse configuration file JSON")
	}
	return config, nil
}

func main() {

	// TODO support set config path by command line option
	config, err := LoadConfiguration("./config.json")
	if err != nil {
		panic(err)
	}

	api.FscInstance = api.FabricSdkClient{
		// Org parameters
		OrgAdmin: config.Org.Admin,
		OrgName:  config.Org.Name,

		ConfigFile: config.ConfigPath,

		// User parameters
		UserName: config.User.Name,
	}

	err = api.FscInstance.Initialize()
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
