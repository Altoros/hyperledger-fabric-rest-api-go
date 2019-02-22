package api

import (
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
)

func (fsc *FabricSdkClient) Query(channelId, chaincodeId, fcn string, args []string) (string, error) {

	// Prepare arguments
	requestArgs := [][]byte{[]byte(fcn)}
	for _, arg := range args {
		requestArgs = append(requestArgs, []byte(arg))
	}

	client, err := fsc.channelClient(channelId)
	if err != nil {
		return "", fmt.Errorf("failed to create channel client")
	}

	response, err := client.Query(
		channel.Request{ChaincodeID: chaincodeId, Fcn: "invoke", Args: requestArgs},
		channel.WithTargets(fsc.getFirstPeer()),
	)
	if err != nil {
		return "", fmt.Errorf("failed to query: %v", err)
	}

	return string(response.Payload), nil
}
