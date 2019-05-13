package api

import (
	"fabric-rest-api-go/pkg/sdk"
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
)

func Invoke(fsc sdk.ChannelClientProvider, channelId, chaincodeId, fcn string, args []string, peers []fab.Peer) (string, error) {

	// Prepare arguments
	var requestArgs [][]byte
	for _, arg := range args {
		requestArgs = append(requestArgs, []byte(arg))
	}

	// Add data that will be visible in the proposal, like a description of the invoke request
	transientDataMap := make(map[string][]byte)
	transientDataMap["result"] = []byte("Transient data in hello invoke")

	// TODO txStatus event support
	/*_, txStatusEventNotifier, err := eventClient.RegisterTxStatusEvent()
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		txStatusEventListener := func() {
			for {
				txEvent := <-txStatusEventNotifier
				notifications.HandleTxStatusEvent(txEvent)
			}
		}
		go txStatusEventListener()*/

	client, err := fsc.ChannelClient(channelId)
	if err != nil {
		return "", err
	}

	// Create a request (proposal) and send it
	response, err := client.Execute(
		channel.Request{ChaincodeID: chaincodeId, Fcn: fcn, Args: requestArgs, TransientMap: transientDataMap},
		channel.WithTargets(peers...),
	)
	if err != nil {
		return "", fmt.Errorf("failed to move funds: %v", err)
	}

	return string(response.TransactionID), nil
}
