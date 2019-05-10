package notifications

import (
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"sync"
)

var mutex = &sync.Mutex{}

func WriteSocket(ws *websocket.Conn, msg []byte) error {

	mutex.Lock()
	err := ws.WriteMessage(websocket.TextMessage, msg)
	mutex.Unlock()

	if err != nil {
		return errors.Wrap(err, "socket write error")
	}

	return nil
}
