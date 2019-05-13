package handlers

import (
	"fabric-rest-api-go/pkg/sdk"
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"regexp"
)

type ApiContext struct {
	echo.Context

	fsc *sdk.FabricSdkClient
}

func (c *ApiContext) Fsc() *sdk.FabricSdkClient {
	return c.fsc
}

func (c *ApiContext) SetFsc(fsc *sdk.FabricSdkClient) {
	c.fsc = fsc
}

type PeerParsed struct {
	Peer, Org string
}

func (c *ApiContext) ParsePeers(peersStrings []string) ([]fab.Peer, error) {
	var peers []fab.Peer
	for _, peerString := range peersStrings {
		if peerParsed, success := ParsePeerFormat(peerString); success {
			peer, err := c.Fsc().GetPeerByOrgAndServerName(peerParsed.Org, fmt.Sprintf(`^%s\.`, peerParsed.Peer))

			if err != nil {
				return nil, errors.Wrapf(err, `problem while fetching peer by template "%s"`, peerString)
			}

			peers = append(peers, peer)
		}
	}

	return peers, nil
}

func ParsePeerFormat(peerString string) (*PeerParsed, bool) {
	r, _ := regexp.Compile(`^(?P<ORG>[^/]*)/(?P<PEER>[^/]*)$`)
	if r.MatchString(peerString) {
		sm := r.FindStringSubmatch(peerString)
		return &PeerParsed{Org: sm[1], Peer: sm[2]}, true
	}
	return nil, false
}

func (c *ApiContext) CurrentPeer() (fab.Peer, error) {
	peerString := c.FormValue("peer")

	if peerParsed, success := ParsePeerFormat(peerString); success {
		peer, err := c.Fsc().GetPeerByOrgAndServerName(peerParsed.Org, fmt.Sprintf(`^%s\.`, peerParsed.Peer))

		if err != nil {
			return nil, errors.Wrapf(err, `problem while fetching peer by template "%s"`, peerString)
		}

		return peer, nil
	}

	randomPeer, err := c.fsc.GetRandomPeer()
	if err != nil {
		return nil, err
	}

	return randomPeer, nil
}
