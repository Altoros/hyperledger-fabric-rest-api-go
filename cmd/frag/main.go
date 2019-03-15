package main

import (
	"fabric-rest-api-go/pkg/api"
	"fabric-rest-api-go/pkg/handlers"
	"flag"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type ApiConfig struct {
	Org struct {
		Admin string `yaml:"admin"`
		Name  string `yaml:"name"`
	} `yaml:"org"`
	User struct {
		Name string `yaml:"name"`
	} `yaml:"user"`
}

func LoadConfiguration(file string) (*ApiConfig, error) {
	var config *ApiConfig
	configFile, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, errors.WithMessage(err, "Unable to open configuration file")
	}
	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		return nil, errors.WithMessage(err, "Unable to parse configuration file JSON")
	}
	return config, nil
}

func main() {
	var apiConfigPath string
	var sdkConfigPath string
	flag.StringVar(&apiConfigPath, "api-config", "./configs/api.yaml", "Path to API configuration file (example: -api-config=./api.yaml)")
	flag.StringVar(&sdkConfigPath, "sdk-config", "./configs/network.yaml", "Path to SDK configuration file (example: -sdk-config=./network.yaml)")
	flag.Parse()

	config, err := LoadConfiguration(apiConfigPath)
	if err != nil {
		panic(err)
	}

	fsc := api.FabricSdkClient{
		ConfigFile: sdkConfigPath,

		// Org parameters
		OrgAdmin: config.Org.Admin,
		OrgName:  config.Org.Name,

		// User parameters
		UserName: config.User.Name,
	}

	err = fsc.Initialize()
	if err != nil {
		panic(err)
	}

	fmt.Println("Start listening to localhost:8080")

	e := echo.New()

	e.Use(func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &handlers.ApiContext{Context: c}
			cc.SetFsc(&fsc)
			return h(cc)
		}
	})

	e.Use(middleware.CORS()) // TODO make this optional

	e.GET("/", handlers.WelcomeHandler)
	e.GET("/health", handlers.HealthCheckHandler)

	e.POST("/chaincodes/install", handlers.PostChaincodesInstallHandler)
	e.POST("/chaincodes/instantiate", handlers.PostChaincodesInstantiateHandler)

	e.GET("/chaincodes/installed", handlers.GetChaincodesInstalledHandler)

	e.GET("/channels/:channelId/chaincodes/instantiated", handlers.GetChaincodesInstantiatedHandler) // TODO
	e.GET("/channels/:channelId/chaincodes/:chaincodeId/info", handlers.GetChaincodesInfoHandler)
	e.GET("/channels", handlers.GetChannelsHandler)
	e.POST("/channels", handlers.PostChannelsHandler) // TODO

	e.GET("/channels/:channelId", handlers.GetChannelsChannelIdHandler)
	e.GET("/channels/:channelId/orgs", handlers.GetChannelsChannelIdOrgsHandler) // TODO
	e.GET("/channels/:channelId/peers", handlers.GetChannelsChannelIdPeersHandler)

	e.GET("/channels/:channelId/chaincodes/:chaincodeId", handlers.GetQueryHandler)
	e.POST("/channels/:channelId/chaincodes/:chaincodeId", handlers.PostInvokeHandler)

	e.POST("/init_test_fixtures", handlers.InitTestFixturesHandler) // TODO remove, for test purposes only

	e.Logger.Fatal(e.Start(":8080"))
}
