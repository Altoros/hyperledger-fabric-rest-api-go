package api

import (
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/common/cauthdsl"
	"github.com/pkg/errors"
)

func ChaincodeInstantiate(channelClientProvider AdminProvider, _ fab.Peer, channelId, chaincodeId, chaincodeVersion string) (string, error) {

	// TODO implement policy parameters
	// Set up chaincode policy
	ccPolicy := cauthdsl.SignedByAnyMember([]string{"Org1MSP"})

	// TODO find out, seems like Path is redundant
	resp, err := channelClientProvider.Admin().InstantiateCC(channelId, resmgmt.InstantiateCCRequest{Name: chaincodeId, Path: "chaincode", Version: chaincodeVersion, Args: [][]byte{[]byte("init")}, Policy: ccPolicy})
	if err != nil || resp.TransactionID == "" {
		return "", errors.WithMessage(err, "failed to instantiate the chaincode")
	}
	fmt.Println("Chaincode instantiated")

	return "ok", nil
}
