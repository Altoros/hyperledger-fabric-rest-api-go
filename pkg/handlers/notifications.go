package handlers

import (
	"encoding/json"
	"fabric-rest-api-go/pkg/context"
	"fabric-rest-api-go/pkg/notifications"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"net/http"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(_ *http.Request) bool {
		return true // TODO remove, only for testing purposes
	},
}

var wsList []*websocket.Conn

func NotificationsHandler(ec echo.Context) error {
	c := ec.(*context.ApiContext)

	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	defer ws.Close()

	wsList = append(wsList, ws)

LOOP:
	for {
		messageType, msg, err := ws.ReadMessage()
		if err != nil {

			if ce, ok := err.(*websocket.CloseError); ok {
				switch ce.Code {
				case websocket.CloseNormalClosure,
					websocket.CloseGoingAway,
					websocket.CloseNoStatusReceived:
					c.Logger().Info("Web socket closed by client: %s", err)
					break LOOP
				}
			}

			c.Logger().Error(err)
			break LOOP
		}

		if messageType == websocket.TextMessage {
			err = ProcessMessage(c, ws, msg)

			if err != nil {
				c.Logger().Error(err)
			}
		}

		c.Logger().Info("Message received: %s\n", msg)
	}

	return nil
}

type EventSubscribeMessage struct {
	Event       string `json:"event"`
	ChannelId   string `json:"channel"`
	ChaincodeId string `json:"chaincode"`
	CcEventId   string `json:"cc_event"`
}

type NotificationsServiceMessage struct {
	Message string `json:"message"`
}

func ProcessMessage(c *context.ApiContext, ws *websocket.Conn, msg []byte) error {
	EventSubscribeMessage := EventSubscribeMessage{}
	err := json.Unmarshal(msg, &EventSubscribeMessage)

	if err != nil {
		return err
	}

	switch EventSubscribeMessage.Event {

	case notifications.FilteredBlockEvent:
		err := notifications.CreateFilteredBlockEventListener(
			c.Fsc(),
			ws,
			EventSubscribeMessage.ChannelId,
		)

		if err != nil {
			return err
		}

		nsm := NotificationsServiceMessage{
			Message: fmt.Sprintf("Successfully subscribed to channel %s block events",
				EventSubscribeMessage.ChannelId,
			),
		}
		nsmJson, _ := json.Marshal(nsm)

		//err = ws.WriteMessage(websocket.TextMessage, nsmJson)
		err = notifications.WriteSocket(ws, nsmJson)
		if err != nil {
			return err
		}
		break

	case notifications.CcEvent:
		err := notifications.CreateChaincodeEventListener(
			c.Fsc(),
			ws,
			EventSubscribeMessage.ChannelId,
			EventSubscribeMessage.ChaincodeId,
			EventSubscribeMessage.CcEventId,
		)

		if err != nil {
			return err
		}

		nsm := NotificationsServiceMessage{
			Message: fmt.Sprintf("Successfully subscribed to channel %s, chaincode %s events %s",
				EventSubscribeMessage.ChannelId,
				EventSubscribeMessage.ChaincodeId,
				EventSubscribeMessage.CcEventId,
			),
		}
		nsmJson, _ := json.Marshal(nsm)

		err = notifications.WriteSocket(ws, nsmJson)
		//err =  ws.WriteMessage(websocket.TextMessage, nsmJson)
		if err != nil {
			return err
		}
		break
	}

	return nil
}
