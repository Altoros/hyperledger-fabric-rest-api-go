package api

import (
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/common/cauthdsl"
	"github.com/pkg/errors"
)

// TODO remove all testing fixed data

const (
	testChannelId     = "mychannel"
	testChaincodeName = "mycc"
	channelConfigTx   = "./test/fabric-samples/basic-network/config/channel.tx"
	testOrdererId     = "orderer.example.com"
	chaincodePath     = "./chaincode"
	projectName       = "project"
)

// Create test channel, install and instantiate test chaincode
func (fsc *FabricSdkClient) InitBasicTestFixturesHandler() error {
	ordererEndPoint := resmgmt.WithOrdererEndpoint(testOrdererId)

	response, _ := fsc.admin.QueryChannels(resmgmt.WithTargets(fsc.GetCurrentPeer()))
	channels := response.GetChannels()

	channelExist := false
	for _, ch := range channels {
		if ch.ChannelId == testChannelId {
			channelExist = true
		}
	}

	if !channelExist {
		req := resmgmt.SaveChannelRequest{ChannelID: testChannelId, ChannelConfigPath: channelConfigTx, SigningIdentities: []msp.SigningIdentity{fsc.adminIdentity}}
		txID, err := fsc.admin.SaveChannel(req, ordererEndPoint)
		if err != nil || txID.TransactionID == "" {
			return errors.WithMessage(err, "failed to save channel")
		}
		fmt.Println("Channel created")

		// Make admin user join the previously created channel
		if err = fsc.admin.JoinChannel(testChannelId, resmgmt.WithRetry(retry.DefaultResMgmtOpts), ordererEndPoint); err != nil {
			return errors.WithMessage(err, "failed to make admin join channel")
		}
		fmt.Println("Channel joined")
	}

	queryInstalledChaincodesResponse, _ := fsc.admin.QueryInstalledChaincodes(resmgmt.WithTargets(fsc.GetCurrentPeer()))
	installedChaincodes := queryInstalledChaincodesResponse.GetChaincodes()

	chaincodeExists := false
	for _, cc := range installedChaincodes {
		if string(cc.Name) == testChaincodeName {
			chaincodeExists = true
		}
	}

	if !chaincodeExists {

		// Create the chaincode package that will be sent to the peers
		ccPkg, err := CCPackage(chaincodePath, projectName)
		if err != nil {
			return errors.WithMessage(err, "failed to create chaincode package")
		}
		fmt.Println("ccPkg created")

		// TODO research path usage inside peer, seems like it is only comment of some kind
		// Install example cc to org peers
		installCCReq := resmgmt.InstallCCRequest{Name: testChaincodeName, Path: projectName + "/chaincode/", Version: "0", Package: ccPkg}

		installCcResponses, err := fsc.admin.InstallCC(installCCReq, resmgmt.WithRetry(retry.DefaultResMgmtOpts))
		if err != nil {
			return errors.WithMessage(err, "failed to install chaincode")
		}

		fmt.Println("Peers responses in install request")
		for _, installCcResponse := range installCcResponses {
			fmt.Printf("Info: %s  Target: %s  Status %d\n", installCcResponse.Info, installCcResponse.Target, installCcResponse.Status)
		}
		fmt.Println("Chaincode installed")

		// Set up chaincode policy
		ccPolicy := cauthdsl.SignedByAnyMember([]string{"Org1MSP"})

		// TODO find out, seems like Path is redundant
		resp, err := fsc.admin.InstantiateCC(testChannelId, resmgmt.InstantiateCCRequest{Name: testChaincodeName, Path: "chaincode", Version: "0", Args: [][]byte{[]byte("init")}, Policy: ccPolicy})
		if err != nil || resp.TransactionID == "" {
			return errors.WithMessage(err, "failed to instantiate the chaincode")
		}
		fmt.Println("Chaincode instantiated")
	}

	return nil
}
