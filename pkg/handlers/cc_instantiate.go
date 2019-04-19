package handlers

import (
	"fabric-rest-api-go/pkg/api"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

type InstantiateCcRequest struct {
	CcName    string   `json:"name"`
	CcVersion string   `json:"version"`
	Channel   string   `json:"channel"`
	Policy    string   `json:"policy"`
	Args      []string `json:"args"`
	Peers     []string `json:"peers"`
}

func PostChaincodesInstantiateHandler(ec echo.Context) error {
	c := ec.(*ApiContext)

	instantiateCcRequest := new(InstantiateCcRequest)
	if err := c.Bind(instantiateCcRequest); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if instantiateCcRequest.CcName == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Chaincode name is required")
	}

	if instantiateCcRequest.CcVersion == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Chaincode version is required")
	}

	if instantiateCcRequest.Channel == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Channel name is required")
	}

	if instantiateCcRequest.Policy == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Chaincode policy is required")
	}

	peers, err := c.ParsePeers(instantiateCcRequest.Peers)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	for _, peer := range peers {
		if !api.CheckChannelExist(c.Fsc(), peer, instantiateCcRequest.Channel) {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Channel %s not exist on peer %s", instantiateCcRequest.Channel, peer.URL()))
		}
	}

	resultString, err := api.ChaincodeInstantiate(c.Fsc(), peers, instantiateCcRequest.Channel, instantiateCcRequest.CcName, instantiateCcRequest.CcVersion, instantiateCcRequest.Policy, instantiateCcRequest.Args)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSONBlob(http.StatusOK, []byte( fmt.Sprintf(`{"result": "%s"}`, resultString)))
}
