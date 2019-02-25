package api

import (
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
)

func Query(channelClientProvider ChannelClientProvider, peer fab.Peer, channelId, chaincodeId, fcn string, args []string) (string, error) {

	// Prepare arguments
	var requestArgs [][]byte
	for _, arg := range args {
		requestArgs = append(requestArgs, []byte(arg))
	}

	client, err := channelClientProvider.ChannelClient(channelId)
	if err != nil {
		return "", fmt.Errorf("failed to create channel client")
	}

	response, err := client.Query(
		channel.Request{ChaincodeID: chaincodeId, Fcn: fcn, Args: requestArgs},
		channel.WithTargets(peer),
	)
	if err != nil {
		return "", fmt.Errorf("failed to query: %v", err)
	}

	return string(response.Payload), nil
}
