package context

import (
	"fabric-rest-api-go/pkg/sdk"
	"fmt"
	"github.com/go-playground/universal-translator"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"gopkg.in/go-playground/validator.v9"
	"regexp"
	"strings"
)

type ApiContext struct {
	echo.Context

	fsc        *sdk.FabricSdkClient
	Translator ut.Translator
}

func (c *ApiContext) Fsc() *sdk.FabricSdkClient {
	return c.fsc
}

func (c *ApiContext) SetFsc(fsc *sdk.FabricSdkClient) {
	c.fsc = fsc
}

type ValidationProvider interface {
	Validate(i interface{}) error
	ValidationErrors(err error) error
}

// converts go-playground/validator ValidationErrors to one simple combined english message error
func (c *ApiContext) ValidationErrors(err error) error {
	errs := err.(validator.ValidationErrors)

	errorsTranslations := errs.Translate(c.Translator)

	var combinedErrors []string
	for _, errorTranslation := range errorsTranslations {
		combinedErrors = append(combinedErrors, fmt.Sprintf("%s.", errorTranslation))
	}

	return errors.New(strings.Join(combinedErrors, " "))
}

type PeerParsed struct {
	Peer, Org string
}

func (c *ApiContext) ParseApiPeers(apiPeersStrings []string) ([]sdk.ApiPeer, error) {
	var peers []sdk.ApiPeer
	for _, peerString := range apiPeersStrings {
		peer, ok := c.Fsc().ApiConfig.Peers[peerString]
		if ok {
			peers = append(peers, peer)
		} else {
			return nil, errors.Errorf(`failed to parse peer string "%s"`, peerString)
		}
	}

	return peers, nil
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
		} else {
			return nil, errors.Errorf(`failed to parse peer string "%s"`, peerString)
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
