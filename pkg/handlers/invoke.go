package handlers

import (
	"fabric-rest-api-go/pkg/api"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

type InvokeRequest struct {
	Fcn   string   `json:"fcn"`
	Args  []string `json:"args"`
	Peers []string `json:"peers"`
}

func PostInvokeHandler(ec echo.Context) error {
	c := ec.(*ApiContext)

	invokeRequest := new(InvokeRequest)
	if err := c.Bind(invokeRequest); err != nil {
		return err
	}

	fcn := invokeRequest.Fcn
	args := invokeRequest.Args

	if fcn == "" {
		return c.String(http.StatusBadRequest, "Fcn is required")
	}

	peers, err := c.ParsePeers(invokeRequest.Peers)
	if err != nil {
		return err
	}

	resultString, err := api.Invoke(c.Fsc(), c.Param("channelId"), c.Param("chaincodeId"), fcn, args, peers)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSONBlob(http.StatusOK, []byte(fmt.Sprintf(`{"result": "%s"}`, resultString)))
}
