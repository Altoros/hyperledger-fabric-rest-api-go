package chaincode_test

import (
	"fabric-rest-api-go/pkg/cc_testing"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

var (
	channelId   = os.Getenv("CHANNEL")
	chaincodeId = os.Getenv("CC_ID")
)

func TestFirst(t *testing.T) {
	testVar := "42"

	cc_testing.EventsInit(t, channelId, chaincodeId)

	cc_testing.Invoke(t, channelId, chaincodeId, "setA", []string{testVar})

	q := cc_testing.Query(t, channelId, chaincodeId, "getA", []string{})

	t.Log("Assert that value is changed")
	assert.Equal(t, testVar, q)

	cc_testing.Invoke(t, channelId, chaincodeId, "event", []string{})

	time.Sleep(time.Second * 2)

	t.Log("Assert that block event is received")
	assert.Contains(t, cc_testing.PopFilteredBlockEvent(), `"event_type":"block"`)

	t.Log(`Assert that chaincode event "testEvent" is received`)
	assert.Contains(t, cc_testing.PopChaincodeEvent(), `"EventName":"testEvent"`)
}
