package notifications

import (
	"container/list"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/event"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
)

type EventClientProvider interface {
	EventClient(channelId string) (*event.Client, error)
}

type FilteredBlockEventNotification struct {
	EventType string                  `json:"event_type"`
	Payload   *fab.FilteredBlockEvent `json:"payload"`
}

type BlockEventNotification struct {
	EventType string          `json:"event_type"`
	Payload   *fab.BlockEvent `json:"payload"`
}

type ChaincodeEventNotification struct {
	EventType string       `json:"event_type"`
	Payload   *fab.CCEvent `json:"payload"`
}

const (
	FilteredBlockEvent = "block"
	BlockEvent         = "full_block"
	CcEvent            = "cc_event"
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

			eventNotification := FilteredBlockEventNotification{EventType: FilteredBlockEvent, Payload: bEvent}

			jsonEventNotification, err := json.Marshal(eventNotification)
			if err != nil {
				break
			}

			err = WriteSocket(ws, jsonEventNotification)
			if err != nil {
				break
			}
		}
	}
	go filteredBlockEventListener()

	return nil
}

var WsConnPool = make(map[string]*list.List)

func CreateChaincodeEventListener(ecp EventClientProvider, ws *websocket.Conn, channelId, chaincodeId, eventIdTemplate string) error {
	eventClient, err := ecp.EventClient(channelId)
	if err != nil {
		return err
	}

	if eventIdTemplate == "" {
		// if eventIdTemplate is empty, will listen all events
		eventIdTemplate = ".*"
	}

	EventKey := fmt.Sprintf("{%s}{%s}{%s}", channelId, chaincodeId, eventIdTemplate)

	// add connection to pool, if event already registered
	if WsConnPool[EventKey] != nil {
		WsConnPool[EventKey].PushFront(ws)
		return nil
	}
	WsConnPool[EventKey] = list.New()
	WsConnPool[EventKey].PushFront(ws)

	reg, chaincodeEventNotifier, err := eventClient.RegisterChaincodeEvent(chaincodeId, eventIdTemplate)
	if err != nil {
		return err
	}

	chaincodeEventListener := func() {
		defer eventClient.Unregister(reg)

		for {
			ccEvent := <-chaincodeEventNotifier

			eventNotification := ChaincodeEventNotification{EventType: CcEvent, Payload: ccEvent}

			jsonEventNotification, err := json.Marshal(eventNotification)
			if err != nil {
				break
			}

			l := WsConnPool[EventKey]
			for e := l.Front(); e != nil; e = e.Next() {
				conn := e.Value.(*websocket.Conn)
				err = WriteSocket(conn, jsonEventNotification)
				if err != nil {
					l.Remove(e)
				}
			}
		}
	}
	go chaincodeEventListener()

	return nil
}

// TODO split listeners, activate channel listener on invoke

// TODO add handler
// TODO find out about permissions
func CreateBlockEventListener(ecp EventClientProvider, ws *websocket.Conn, channelId string) error {
	eventClient, err := ecp.EventClient(channelId)
	if err != nil {
		return err
	}

	reg, blockEventNotifier, err := eventClient.RegisterBlockEvent()
	if err != nil {
		return err
	}

	blockEventListener := func() {
		defer eventClient.Unregister(reg)

		for {
			fbEvent := <-blockEventNotifier

			eventNotification := BlockEventNotification{EventType: FilteredBlockEvent, Payload: fbEvent}

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
	go blockEventListener()

	return nil
}
