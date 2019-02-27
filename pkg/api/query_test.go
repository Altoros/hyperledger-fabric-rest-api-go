package api

import (
	"fabric-rest-api-go/pkg/api/mock_api"
	"github.com/golang/mock/gomock"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	txnmocks "github.com/hyperledger/fabric-sdk-go/pkg/client/common/mocks"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/context"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	contextImpl "github.com/hyperledger/fabric-sdk-go/pkg/context"
	fcmocks "github.com/hyperledger/fabric-sdk-go/pkg/fab/mocks"
	mspmocks "github.com/hyperledger/fabric-sdk-go/pkg/msp/test/mockmsp"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	queryChannelID = "testChannel"
	queryCcID      = "ccId"
	queryFcn       = "fcn"
)

func TestQuery(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testPeer1 := fcmocks.NewMockPeer("Peer1", "http://peer1.com")
	testPeer1.Payload = []byte("mock payload")
	args := []string{"arg1", "arg2"}

	m := mock_api.NewMockChannelClientProvider(ctrl)
	m.EXPECT().ChannelClient(queryChannelID).Return(nil, errors.New("error")).Times(1)

	_, err := Query(m, testPeer1, queryChannelID, queryCcID, queryFcn, args)
	assert.Error(t, err)

	// Make query against mock peer, expect payload

	peers := []fab.Peer{testPeer1}
	chClient := setupChannelClient(peers, t)

	m.EXPECT().ChannelClient(queryChannelID).Return(chClient, nil).Times(1)
	response, err := Query(m, testPeer1, queryChannelID, queryCcID, queryFcn, args)

	assert.Equal(t, "mock payload", response)
	assert.Equal(t, nil, err)
}

func setupChannelClient(peers []fab.Peer, t *testing.T) *channel.Client {

	return setupChannelClientWithError(nil, nil, peers, t)
}

func setupChannelClientWithError(discErr error, selectionErr error, peers []fab.Peer, t *testing.T) *channel.Client {
	fabCtx := setupCustomTestContext(t, txnmocks.NewMockSelectionService(selectionErr, peers...), txnmocks.NewMockDiscoveryService(discErr), nil)

	ctx := createChannelContext(fabCtx, queryChannelID)

	ch, err := channel.New(ctx)
	if err != nil {
		t.Fatalf("Failed to create new channel client: %s", err)
	}

	return ch
}

func setupCustomTestContext(t *testing.T, selectionService fab.SelectionService, discoveryService fab.DiscoveryService, orderers []fab.Orderer) context.ClientProvider {
	user := mspmocks.NewMockSigningIdentity("test", "test")
	ctx := fcmocks.NewMockContext(user)

	if orderers == nil {
		orderer := fcmocks.NewMockOrderer("", nil)
		orderers = []fab.Orderer{orderer}
	}

	transactor := txnmocks.MockTransactor{
		Ctx:       ctx,
		ChannelID: queryChannelID,
		Orderers:  orderers,
	}

	testChannelSvc, err := setupTestChannelService(ctx)
	assert.Nil(t, err, "Got error %s", err)

	mockChService := testChannelSvc.(*fcmocks.MockChannelService)
	mockChService.SetTransactor(&transactor)
	mockChService.SetDiscovery(discoveryService)
	mockChService.SetSelection(selectionService)

	channelProvider := ctx.MockProviderContext.ChannelProvider()
	channelProvider.(*fcmocks.MockChannelProvider).SetCustomChannelService(testChannelSvc)

	return createClientContext(ctx)
}

func setupTestChannelService(ctx context.Client) (fab.ChannelService, error) {
	chProvider, err := fcmocks.NewMockChannelProvider(ctx)
	if err != nil {
		return nil, errors.WithMessage(err, "mock channel provider creation failed")
	}

	chService, err := chProvider.ChannelService(ctx, queryChannelID)
	if err != nil {
		return nil, errors.WithMessage(err, "mock channel service creation failed")
	}

	return chService, nil
}

func createChannelContext(clientContext context.ClientProvider, channelID string) context.ChannelProvider {

	channelProvider := func() (context.Channel, error) {
		return contextImpl.NewChannel(clientContext, channelID)
	}

	return channelProvider
}

func createClientContext(client context.Client) context.ClientProvider {
	return func() (context.Client, error) {
		return client, nil
	}
}
