package sdk

import (
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/event"
	mspclient "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/pkg/errors"
	"math/rand"
	"reflect"
	"regexp"
	"time"
)

// FabricSdkClient implementation
type FabricSdkClient struct {
	ConfigFile string

	ApiConfig *Config

	initialized bool

	admin *resmgmt.Client
	sdk   *fabsdk.FabricSDK

	adminIdentity msp.SigningIdentity

	allPeersByOrg              map[string][]fab.Peer
	allPeers                   []fab.Peer
	allPeersByOrgAndServerName map[string]map[string]fab.Peer
}

// Initialize reads the configuration file and sets up the client, chain and event hub
func (fsc *FabricSdkClient) Initialize() error {

	// Add parameters for the initialization
	if fsc.initialized {
		return errors.New("sdk already initialized")
	}

	// Initialize the SDK with the configuration file
	sdk, err := fabsdk.New(config.FromFile(fsc.ConfigFile))
	if err != nil {
		return errors.WithMessage(err, "failed to create SDK")
	}
	fsc.sdk = sdk
	fmt.Println("SDK created")

	// The resource management client is responsible for managing Channels (create/update channel)
	resourceManagerClientContext := fsc.sdk.Context(fabsdk.WithUser(fsc.ApiConfig.Org.Admin), fabsdk.WithOrg(fsc.ApiConfig.Org.Name))
	if err != nil {
		return errors.WithMessage(err, "failed to load Admin identity")
	}
	resMgmtClient, err := resmgmt.New(resourceManagerClientContext)
	if err != nil {
		return errors.WithMessage(err, "failed to create channel management client from Admin identity")
	}
	fsc.admin = resMgmtClient
	fmt.Println("Resource management client created")

	// The MSP client allow us to retrieve user information from their identity, like its signing identity which we will need to save the channel
	mspClient, err := mspclient.New(sdk.Context(), mspclient.WithOrg(fsc.ApiConfig.Org.Name))
	if err != nil {
		return errors.WithMessage(err, "failed to create MSP client")
	}

	adminIdentity, err := mspClient.GetSigningIdentity(fsc.ApiConfig.Org.Admin)
	if err != nil {
		return errors.WithMessage(err, "failed to get admin signing identity")
	}
	fsc.adminIdentity = adminIdentity

	clientContext, err := fsc.sdk.Context()()
	if err != nil {
		return errors.WithMessage(err, "failed to create client context")
	}
	endpointConfig := clientContext.EndpointConfig()
	networkConfig := endpointConfig.NetworkConfig()

	fsc.allPeersByOrg = make(map[string][]fab.Peer)

	fsc.allPeersByOrgAndServerName = make(map[string]map[string]fab.Peer)

	for orgID := range networkConfig.Organizations {
		fsc.allPeersByOrgAndServerName[orgID] = make(map[string]fab.Peer)

		peersConfig, ok := endpointConfig.PeersConfig(orgID)
		if !ok {
			return errors.Errorf("failed to get peer configs for org [%s]", orgID)
		}

		var peers []fab.Peer
		for _, p := range peersConfig {
			endorser, err := clientContext.InfraProvider().CreatePeerFromConfig(&fab.NetworkPeer{PeerConfig: p})
			if err != nil {
				return errors.Wrapf(err, "failed to create peer from config")
			}
			peers = append(peers, endorser)

			serverName := reflect.ValueOf(endorser).Elem().FieldByName("serverName").String()

			fsc.allPeersByOrgAndServerName[orgID][serverName] = endorser
		}
		fsc.allPeersByOrg[orgID] = peers
		fsc.allPeers = append(fsc.allPeers, peers...)
	}

	fmt.Println("Initialization Successful")
	fsc.initialized = true
	return nil
}

func (fsc *FabricSdkClient) ChannelClient(channelId string) (*channel.Client, error) {
	chProvider := fsc.sdk.ChannelContext(channelId, fabsdk.WithUser(fsc.ApiConfig.User.Name))

	client, err := channel.New(chProvider)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to create new channel client")
	}
	fmt.Println("Channel client created")

	return client, nil
}

func (fsc *FabricSdkClient) EventClient(channelId string) (*event.Client, error) {
	channelProvider := fsc.sdk.ChannelContext(channelId, fabsdk.WithUser(fsc.ApiConfig.User.Name))

	eventClient, err := event.New(channelProvider)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to create new event client")
	}
	fmt.Println("Event client created")

	return eventClient, nil
}

func (fsc *FabricSdkClient) GetRandomPeer() (fab.Peer, error) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	if len(fsc.allPeers) == 0 {
		return nil, errors.New("no peers was loaded in client from configuration file")
	}

	randomPeerN := r.Intn(len(fsc.allPeers))
	return fsc.allPeers[randomPeerN], nil
}

func (fsc *FabricSdkClient) GetPeerByOrgAndServerName(org, serverNameTemplate string) (fab.Peer, error) {
	if len(fsc.allPeersByOrg[org]) == 0 {
		return nil, errors.Errorf(`could't find any peers in "%s" organisation`, org)
	}

	for serverName, peer := range fsc.allPeersByOrgAndServerName[org] {
		match, err := regexp.MatchString(serverNameTemplate, serverName)
		if err != nil {
			return nil, err
		}

		if match {
			return peer, nil
		}
	}

	return nil, errors.Errorf(`could't find peer "%s" in "%s" organisation`, serverNameTemplate, org)
}

func (fsc *FabricSdkClient) Admin() *resmgmt.Client {
	return fsc.admin
}

func (fsc *FabricSdkClient) Sdk() *fabsdk.FabricSDK {
	return fsc.sdk
}

type ChannelClientProvider interface {
	ChannelClient(string) (*channel.Client, error)
}

type AdminProvider interface {
	Admin() *resmgmt.Client
}
