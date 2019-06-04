package cc_testing

import (
	"bytes"
	"fmt"
	"github.com/Jeffail/gabs"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

var filteredBlockEventsQueue []string
var chaincodeEventsQueue []string

func PopFilteredBlockEvent() string {
	for len(filteredBlockEventsQueue) > 0 {
		first := filteredBlockEventsQueue[0]
		filteredBlockEventsQueue = filteredBlockEventsQueue[1:]
		return first
	}
	return ""
}

func PopChaincodeEvent() string {
	for len(chaincodeEventsQueue) > 0 {
		first := chaincodeEventsQueue[0]
		chaincodeEventsQueue = chaincodeEventsQueue[1:]
		return first
	}
	return ""
}

func EventsInit(t *testing.T, channelId, chaincodeId string) {
	FilteredBlockEventsInit(t, channelId)
	ChaincodeEventsInit(t, channelId, chaincodeId)
}

func FilteredBlockEventsInit(t *testing.T, channelId string) {
	conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/notifications", nil)
	if err != nil {
		t.Log("dial:", err)
	}

	err = conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf(`{"event": "block", "channel": "%s"}`, channelId)))
	if err != nil {
		t.Log("send:", err)
	}

	go func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				t.Log("ws err:", err)
				return
			}
			if strings.Contains(string(message), `"event_type"`) {
				filteredBlockEventsQueue = append(filteredBlockEventsQueue, string(message))
			}
			t.Logf("Received from socket: %s", message)
		}
	}()
}

func ChaincodeEventsInit(t *testing.T, channelId, chaincodeId string) {
	conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/notifications", nil)
	if err != nil {
		t.Log("dial:", err)
	}

	err = conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf(`{"event":"cc_event", "channel": "%s", "chaincode": "%s"}`, channelId, chaincodeId)))
	if err != nil {
		t.Log("send:", err)
	}

	go func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				t.Log("ws err:", err)
				return
			}
			if strings.Contains(string(message), `"event_type"`) {
				chaincodeEventsQueue = append(chaincodeEventsQueue, string(message))
			}
			t.Logf("Received from socket: %s", message)
		}
	}()
}

func Invoke(t *testing.T, channelId, chaincodeId, fcn string, args []string) {
	url := fmt.Sprintf("http://localhost:8080/channels/%s/chaincodes/%s", channelId, chaincodeId)

	var quotedArgs []string
	for _, arg := range args {
		quotedArgs = append(quotedArgs, fmt.Sprintf(`"%s"`, arg))
	}

	var jsonStr = []byte(fmt.Sprintf(`{"fcn": "%s", "args": [%s], "peers": ["org1/peer0"]}`, fcn, strings.Join(quotedArgs, ",")))
	t.Logf("Invoke CC: %s", jsonStr)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	t.Logf("response %s: %s", resp.Status, string(body))
}

func Query(t *testing.T, channelId, chaincodeId, fcn string, args []string) string {
	url := fmt.Sprintf("http://localhost:8080/channels/%s/chaincodes/%s", channelId, chaincodeId)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Log(err)
		return ""
	}

	q := req.URL.Query()
	q.Add("fcn", fcn)
	q.Add("args", strings.Join(args, ","))
	q.Add("peer", "org1/peer0")
	req.URL.RawQuery = q.Encode()

	t.Logf("Query CC: %s", req.URL.String())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	t.Logf("API Response %s: %s", resp.Status, string(body))

	jsonParsed, _ := gabs.ParseJSON(body)
	fVal, isFloat := jsonParsed.Path("result").Data().(float64)
	if isFloat {
		return fmt.Sprintf("%g", fVal)
	}

	result, _ := jsonParsed.Path("result").Data().(string)
	return result
}
