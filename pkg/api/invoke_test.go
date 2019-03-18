package api

import (
	"fabric-rest-api-go/pkg/api/mock_api"
	"github.com/golang/mock/gomock"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	fcmocks "github.com/hyperledger/fabric-sdk-go/pkg/fab/mocks"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	invokeChannelID = "testChannel"
	invokeCcID      = "ccId"
	invokeFcn       = "invokeFcn"
)

func TestInvoke(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testPeer1 := fcmocks.NewMockPeer("Peer1", "http://peer1.com")
	args := []string{"arg1", "arg2"}
	peers := []fab.Peer{testPeer1}

	m := mock_api.NewMockChannelClientProvider(ctrl)
	m.EXPECT().ChannelClient(invokeChannelID).Return(nil, errors.New("error")).Times(1)

	_, err := Invoke(m, invokeChannelID, invokeCcID, invokeFcn, args, peers)
	assert.Error(t, err)

	chClient := setupChannelClient(peers, t)

	m.EXPECT().ChannelClient(invokeChannelID).Return(chClient, nil).Times(1)

	response, err := Invoke(m, invokeChannelID, invokeCcID, invokeFcn, args, peers)

	assert.NotEmpty(t, response)
	assert.Equal(t, nil, err)
}
