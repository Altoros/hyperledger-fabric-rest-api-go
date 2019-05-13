package main

import (
	"fabric-rest-api-go/pkg/handlers"
	"fabric-rest-api-go/pkg/sdk"
	"flag"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	var apiConfigPath string
	var sdkConfigPath string
	flag.StringVar(&apiConfigPath, "api-config", "./configs/api.yaml", "Path to API configuration file (example: -api-config=./api.yaml)")
	flag.StringVar(&sdkConfigPath, "sdk-config", "./configs/network.yaml", "Path to SDK configuration file (example: -sdk-config=./network.yaml)")
	flag.Parse()

	apiConfig, err := sdk.LoadConfiguration(apiConfigPath)
	if err != nil {
		panic(err)
	}

	fsc := sdk.FabricSdkClient{
		ConfigFile: sdkConfigPath,
		ApiConfig:  apiConfig,
	}

	err = fsc.Initialize()
	if err != nil {
		panic(err)
	}

	e := echo.New()

	e.Use(func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &handlers.ApiContext{Context: c}
			cc.SetFsc(&fsc)
			return h(cc)
		}
	})

	e.Use(middleware.CORS()) // TODO make this optional

	e.Use(middleware.Recover())

	e.GET("/", handlers.WelcomeHandler)
	e.GET("/health", handlers.HealthCheckHandler)

	e.GET("/config", handlers.GetConfigHandler)

	e.POST("/users", handlers.PostUsersHandler)

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

	// CA and certificates
	e.POST("/ca/enroll", handlers.PostCaEnrollHandler)
	e.POST("/ca/register", handlers.PostCaRegisterHandler)

	e.GET("/notifications", handlers.NotificationsHandler)

	e.Logger.Fatal(e.Start(":8080"))
}
