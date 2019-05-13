package api

import (
	"fabric-rest-api-go/pkg/sdk"
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/common/cauthdsl"
	"github.com/pkg/errors"
)

func ChaincodeInstantiate(channelClientProvider sdk.AdminProvider, peers []fab.Peer, channelId, ccName, chaincodeVersion, policyString string, args []string) (string, error) {

	// TODO implement more complex policy parameters
	// Set up chaincode policy
	// example: ccPolicy := cauthdsl.SignedByAnyMember([]string{"Org1MSP"})

	ccPolicy, err := cauthdsl.FromString(policyString)
	if err != nil {
		return "", err
	}

	requestArgs := [][]byte{[]byte("init"),}
	// Prepare arguments
	for _, arg := range args {
		requestArgs = append(requestArgs, []byte(arg))
	}

	// TODO find out, seems like Path is redundant
	resp, err := channelClientProvider.Admin().InstantiateCC(
		channelId,
		resmgmt.InstantiateCCRequest{Name: ccName, Path: "chaincode/" + ccName + "/", Version: chaincodeVersion, Args: requestArgs, Policy: ccPolicy},
		resmgmt.WithTargets(peers...),
	)
	if err != nil || resp.TransactionID == "" {
		return "", errors.WithMessage(err, "failed to instantiate the chaincode")
	}
	fmt.Println("Chaincode instantiated")

	return "ok", nil
}
