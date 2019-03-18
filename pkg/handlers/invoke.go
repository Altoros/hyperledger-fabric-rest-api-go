package handlers

import (
	"fabric-rest-api-go/pkg/api"
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"net/http"
)

type invokeRequest struct {
	Fcn   string   `json:"fcn"`
	Args  []string `json:"args"`
	Peers []string `json:"peers"`
}

func PostInvokeHandler(ec echo.Context) error {
	c := ec.(*ApiContext)

	m := new(invokeRequest)
	if err := c.Bind(m); err != nil {
		return err
	}

	fcn := m.Fcn
	args := m.Args

	if fcn == "" {
		return c.String(http.StatusBadRequest, "Fcn is required")
	}

	var peers []fab.Peer
	for _, peerString := range m.Peers {
		if peerParsed, success := ParsePeerFormat(peerString); success {
			peer, err := c.Fsc().GetPeerByOrgAndServerName(peerParsed.Org, fmt.Sprintf(`^%s\.`, peerParsed.Peer))

			if err != nil {
				return errors.Wrapf(err, `problem while fetching peer by template "%s"`, peerString)
			}

			peers = append(peers, peer)
		}
	}

	resultString, err := api.Invoke(c.Fsc(), c.Param("channelId"), c.Param("chaincodeId"), fcn, args, peers)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSONBlob(http.StatusOK, []byte(fmt.Sprintf(`{"result": "%s"}`, resultString)))
}
