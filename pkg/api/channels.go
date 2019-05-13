package api

import (
	"encoding/json"
	"errors"
	"fabric-rest-api-go/pkg/sdk"
	"github.com/Jeffail/gabs"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

func Channels(fsc sdk.AdminProvider, peer fab.Peer) (string, error) {
	response, err := fsc.Admin().QueryChannels(resmgmt.WithTargets(peer))
	if err != nil {
		return "", err
	}

	jsonBytes, err := json.Marshal(response.GetChannels())
	if err != nil {
		return "", err
	}

	return string(jsonBytes), nil
}

func ChannelInfo(fsc sdk.AdminProvider, channelId string) (string, error) {

	// TODO implement
	return "", errors.New("not implemented")
	/*response, err := fsc.admin.QueryChannels(resmgmt.WithTargets(fsc.GetRandomPeer()))
	if err != nil {
		return "", err
	}

	var channelInfo *pb.ChannelInfo
	for _, ch := range response.GetChannels() {
		if channelId == ch.ChannelId {
			channelInfo = ch
		}
	}

	if channelInfo == nil {
		return "", fmt.Errorf(`unable to find channel by id "%s"`, channelId)
	}

	jsonBytes, err := json.Marshal(channelInfo)
	if err != nil {
		return "", fmt.Errorf(`API returned mallformes data: %s`, err.Error())
	}

	return string(jsonBytes), nil*/
}

func ChannelOrgs(fsc sdk.AdminProvider, channelId string) (string, error) {
	// TODO implement
	return "", errors.New("not implemented")
}

func ChannelPeers(fsc *sdk.FabricSdkClient, channelId string) (string, error) {
	chProvider := fsc.Sdk().ChannelContext(channelId, fabsdk.WithUser(fsc.ApiConfig.User.Name))

	chContext, err := chProvider()
	if err != nil {
		return "", err
	}

	discovery, err := chContext.ChannelService().Discovery()
	if err != nil {
		return "", err
	}

	chContext.ChannelService().Membership()

	peers, err := discovery.GetPeers()
	if err != nil {
		return "", err
	}

	return PeersToJsonString(peers), nil
}

func PeersToJsonString(peers []fab.Peer) string {
	jsonObj := gabs.New()

	jsonObj.Array("Peers")
	for _, peer := range peers {
		jsonObj.ArrayAppend(PeerToJsonObjects(peer).Data(), "Peers")
	}

	return jsonObj.String()
}

// PeerToJsonObjects prints the peer
func PeerToJsonObjects(peer fab.Peer) *gabs.Container {
	jsonObj := gabs.New()

	jsonObj.SetP(peer.URL(), "URL")
	jsonObj.SetP(peer.MSPID(), "MSP")

	peerState, ok := peer.(fab.PeerState)
	if ok {
		jsonObj.SetP(peerState.BlockHeight(), "BlockHeight")
	}

	return jsonObj
}

func CheckChannelExist(adminProvider sdk.AdminProvider, peer fab.Peer, channelId string) bool{
	response, _ := adminProvider.Admin().QueryChannels(resmgmt.WithTargets(peer))
	channels := response.GetChannels()

	for _, ch := range channels {
		if ch.ChannelId == channelId {
			return true
		}
	}
	return false
}