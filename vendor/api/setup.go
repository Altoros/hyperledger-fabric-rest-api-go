package api

import (
	"notifications"
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
)

var FscInstance FabricSdkClient

// FabricSdkClient implementation
type FabricSdkClient struct {
	ConfigFile string
	OrgID      string

	initialized   bool
	// ChannelConfig string
	GOPATH        string
	ChaincodePath string // TODO remove
	OrgAdmin      string
	OrgName       string
	UserName      string
	//channelClient *channel.Client
	admin *resmgmt.Client
	sdk   *fabsdk.FabricSDK
	// event *event.Client

	adminIdentity msp.SigningIdentity

	allPeersByOrg map[string][]fab.Peer
	allPeers      []fab.Peer
	//installedChaincodes []*peer.ChaincodeInfo
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
	resourceManagerClientContext := fsc.sdk.Context(fabsdk.WithUser(fsc.OrgAdmin), fabsdk.WithOrg(fsc.OrgName))
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
	mspClient, err := mspclient.New(sdk.Context(), mspclient.WithOrg(fsc.OrgName))
	if err != nil {
		return errors.WithMessage(err, "failed to create MSP client")
	}

	adminIdentity, err := mspClient.GetSigningIdentity(fsc.OrgAdmin)
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

	for orgID := range networkConfig.Organizations {
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
			//action.orgIDByPeer[endorser.URL()] = orgID
		}
		fsc.allPeersByOrg[orgID] = peers
		fsc.allPeers = append(fsc.allPeers, peers...)
	}

	/*
	level := levelFromName(cliconfig.Config().LoggingLevel())
	logging.SetLevel("", level)


	for orgID := range networkConfig.Organizations {
		peersConfig, ok := action.endpointConfig.PeersConfig(orgID)
		if !ok {
			return errors.Errorf("failed to get peer configs for org [%s]", orgID)
		}

		cliconfig.Config().Logger().Debugf("Peers for org [%s]: %v\n", orgID, peersConfig)

		var peers []fab.Peer
		for _, p := range peersConfig {
			endorser, err := ctx.InfraProvider().CreatePeerFromConfig(&fab.NetworkPeer{PeerConfig: p})
			if err != nil {
				return errors.Wrapf(err, "failed to create peer from config")
			}
			peers = append(peers, endorser)
			action.orgIDByPeer[endorser.URL()] = orgID
		}
		allPeersByOrg[orgID] = peers
		allPeers = append(allPeers, peers...)
	}
	*/

	// TODO remove, required only for testing
	fsc.InitTestFixturesHandler()

	err = fsc.createEventsListeners("chainhero", "heroes-service") // TODO remove test data
	if err != nil {
		return errors.WithMessage(err, "failed to register and create events listeners")
	}

	fmt.Println("Initialization Successful")
	fsc.initialized = true
	return nil
}

func (fsc *FabricSdkClient) createEventsListeners(channelId, chaincodeId string) error {
	// TODO split listeners, activate channel listener on invoke

	// Creation of the client which will enables access to our channel events
	eventID := ".*"

	eventClient, err := fsc.eventClient(channelId)
	if err != nil {
		return err
	}

	_, chaincodeEventNotifier, err := eventClient.RegisterChaincodeEvent(chaincodeId, eventID)
	if err != nil {
		return err
	}
	// defer eventClient.Unregister(reg)

	chaincodeEventListener := func() {
		for {
			ccEvent := <-chaincodeEventNotifier
			notifications.HandleChaincodeEvent(ccEvent)
		}
	}
	go chaincodeEventListener()

	// TODO find out about permissions
	/*_, blockEventNotifier, err := eventClient.RegisterBlockEvent()
	if err != nil {
		return err
	}

	blockEventListener := func() {
		for {
			bEvent := <-blockEventNotifier
			fmt.Println(bEvent.Block.Data)
		}
	}
	go blockEventListener()*/

	_, filteredBlockEventNotifier, err := eventClient.RegisterFilteredBlockEvent()
	if err != nil {
		return err
	}

	filteredBlockEventListener := func() {
		for {
			bEvent := <-filteredBlockEventNotifier
			notifications.HandleFilteredBlockEvent(bEvent)
		}
	}
	go filteredBlockEventListener()

	return nil
}

func (fsc *FabricSdkClient) channelClient(channelId string) (*channel.Client, error) {
	chProvider := fsc.sdk.ChannelContext(channelId, fabsdk.WithUser(fsc.UserName))

	client, err := channel.New(chProvider)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to create new channel client")
	}
	fmt.Println("Channel client created")

	return client, nil
}

func (fsc *FabricSdkClient) eventClient(channelId string) (*event.Client, error) {
	clientContext := fsc.sdk.ChannelContext(channelId, fabsdk.WithUser(fsc.UserName))

	eventClient, err := event.New(clientContext)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to create new event client")
	}
	fmt.Println("Event client created")

	return eventClient, nil
}

func (fsc *FabricSdkClient) getFirstPeer() fab.Peer {
	return fsc.allPeers[0]
}
