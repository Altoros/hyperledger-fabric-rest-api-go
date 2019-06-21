package chaincode_test

import (
	"fabric-rest-api-go/pkg/cc_testing"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var (
	channelId   = os.Getenv("CHANNEL")
	chaincodeId = os.Getenv("CC_ID")
)

func TestBasicSetGetValue(t *testing.T) {
	testVar := "a"
	testValue := "42"

	cc_testing.EventsInit(t, channelId, chaincodeId)

	cc_testing.Invoke(t, channelId, chaincodeId, "update", []string{testVar, testValue})

	q := cc_testing.Query(t, channelId, chaincodeId, "query", []string{testVar})

	t.Log("Assert that value is changed")
	assert.Equal(t, testValue, q)
}
