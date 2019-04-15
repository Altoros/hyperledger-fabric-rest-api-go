package notifications

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/event"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
)

type EventClientProvider interface {
	EventClient(channelId string) (*event.Client, error)
}

type EventNotification struct {
	EventType string `json:"event_type"`
	Payload   *fab.FilteredBlockEvent `json:"payload"`
}

const (
	BlockEvent = "block"
)

func CreateFilteredBlockEventListener(ecp EventClientProvider, ws *websocket.Conn, channelId string) error {
	eventClient, err := ecp.EventClient(channelId)
	if err != nil {
		return err
	}

	reg, filteredBlockEventNotifier, err := eventClient.RegisterFilteredBlockEvent()
	if err != nil {
		return err
	}

	filteredBlockEventListener := func() {
		defer eventClient.Unregister(reg)

		for {
			bEvent := <-filteredBlockEventNotifier

			eventNotification := EventNotification{EventType: BlockEvent, Payload: bEvent}

			jsonEventNotification, err := json.Marshal(eventNotification)
			if err != nil {
				break
			}

			err = ws.WriteMessage(websocket.TextMessage, jsonEventNotification)
			if err != nil {
				break
			}
		}
	}
	go filteredBlockEventListener()

	return nil
}

func createEventsListeners(ecp EventClientProvider, channelId, chaincodeId string) error {
	// TODO split listeners, activate channel listener on invoke

	// Creation of the client which will enables access to our channel events
	eventID := ".*"

	eventClient, err := ecp.EventClient(channelId)
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
			HandleChaincodeEvent(ccEvent)
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
			fmt.Println(bEvent.Block.CcData)
		}
	}
	go blockEventListener()*/

	return nil
}
