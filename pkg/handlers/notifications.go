package handlers

import (
	"encoding/json"
	"fabric-rest-api-go/pkg/notifications"
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
	c := ec.(*ApiContext)

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
	Event     string `json:"event"`
	ChannelId string `json:"channel"`
}

func ProcessMessage(c *ApiContext, ws *websocket.Conn, msg []byte) error {
	EventSubscribeMessage := EventSubscribeMessage{}
	err := json.Unmarshal(msg, &EventSubscribeMessage)

	if err != nil {
		return err
	}

	switch EventSubscribeMessage.Event {
	case notifications.BlockEvent:
		err := notifications.CreateFilteredBlockEventListener(c.Fsc(), ws, EventSubscribeMessage.ChannelId)
		if err != nil {
			return err
		}
		break
	}

	return nil
}
