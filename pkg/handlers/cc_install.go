package handlers

import (
	"bytes"
	"encoding/base64"
	"fabric-rest-api-go/pkg/api"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

type InstallCcRequest struct {
	CcName    string   `json:"name"`
	CcVersion string   `json:"version"`
	Channel   string   `json:"channel"`
	CcData    string   `json:"data"`
	Peers     []string `json:"peers"`
}

func PostChaincodesInstallHandler(ec echo.Context) error {
	c := ec.(*ApiContext)

	installRequest := new(InstallCcRequest)
	if err := c.Bind(installRequest); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// TODO validation enhancement
	if installRequest.CcName == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Chaincode name is required")
	}

	if installRequest.CcVersion == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Chaincode version is required")
	}

	if installRequest.Channel == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Channel name is required")
	}

	if installRequest.CcData == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Chaincode data is required")
	}

	if installRequest.Peers == nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Peers are required")
	}

	peers, err := c.ParsePeers(installRequest.Peers)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())

	}

	if peers == nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Peers not found")
	}

	for _, peer := range peers {
		if !api.CheckChannelExist(c.Fsc(), peer, installRequest.Channel) {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Channel %s not exist on peer %s", installRequest.Channel, peer.URL()))
		}
	}

	// preparing chaincode tarball
	ccTarBytes, err := base64.StdEncoding.DecodeString(installRequest.CcData)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	ccTarBytes, err = PrependPathToTar(bytes.NewReader(ccTarBytes), "src/chaincode/"+installRequest.CcName+"/")
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	resultString, err := api.ChaincodeInstall(c.Fsc(), peers, installRequest.Channel, installRequest.CcName, installRequest.CcVersion, ccTarBytes)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSONBlob(http.StatusOK, []byte( fmt.Sprintf(`{"result": "%s"}`, resultString)))
}
