package tx

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fabric-rest-api-go/pkg/sdk"
	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric/common/crypto"
	"github.com/hyperledger/fabric/protos/common"
	"github.com/hyperledger/fabric/protos/msp"
	"github.com/hyperledger/fabric/protos/peer"
	"github.com/hyperledger/fabric/protos/utils"
	"github.com/pkg/errors"
)

func CreateProposal(ChannelID, ChaincodeID string, argsArray [][]byte, mspid string, userCert []byte) (*peer.Proposal, error) {
	ccis := &peer.ChaincodeInvocationSpec{ChaincodeSpec: &peer.ChaincodeSpec{
		Type: peer.ChaincodeSpec_GOLANG, ChaincodeId: &peer.ChaincodeID{Name: ChaincodeID},
		Input: &peer.ChaincodeInput{Args: argsArray}}}

	nonce, err := crypto.GetRandomNonce()
	if err != nil {
		return nil, errors.Wrap(err, "unable to get random nonce")
	}

	transientMap := make(map[string][]byte)

	creator, err := CombineIdentity(mspid, userCert)
	if err != nil {
		return nil, errors.Wrap(err, "unable to compose identity")
	}

	transactionID, err := СomputeTxnID(nonce, creator)
	if err != nil {
		return nil, errors.Wrap(err, "unable to compute transactionID")
	}

	// TODO switch to CreateChaincodeProposal
	proposal, _, err := utils.CreateChaincodeProposalWithTxIDNonceAndTransient(transactionID, common.HeaderType_ENDORSER_TRANSACTION, ChannelID, ccis, nonce, creator, transientMap)
	return proposal, err
}

func CombineIdentity(mspid string, userCert []byte) ([]byte, error) {
	serializedIdentity := &msp.SerializedIdentity{
		Mspid:   mspid,
		IdBytes: userCert,
	}
	identity, err := proto.Marshal(serializedIdentity)
	if err != nil {
		return nil, errors.Wrap(err, "marshal serializedIdentity failed")
	}
	return identity, nil
}

func СomputeTxnID(nonce, creator []byte) (string, error) {
	h := sha256.New()

	b := append(nonce, creator...)

	_, err := h.Write(b)
	if err != nil {
		return "", err
	}
	digest := h.Sum(nil)
	id := hex.EncodeToString(digest)

	return id, nil
}

func ProposalHash(proposal *peer.Proposal) ([]byte, error) {
	proposalBytes, err := utils.GetBytesProposal(proposal)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get proposalBytes")
	}

	hasher := sha256.New()
	hasher.Write(proposalBytes)
	proposalSha256 := hasher.Sum(nil)
	return proposalSha256, nil
}

func SendProposalToPeersEndorsement(signedProposal *peer.SignedProposal, targets []sdk.ApiPeer) ([]*peer.ProposalResponse, error) {
	var proposalResponses []*peer.ProposalResponse

	for _, target := range targets {
		ctx := context.TODO()
		conn, err := target.GrpcConn(ctx)
		if err != nil {
			return nil, errors.Wrapf(err, "grpc connection error to peer, %v", target)
		}

		endorserClient := peer.NewEndorserClient(conn)
		proposalResponse, err := endorserClient.ProcessProposal(ctx, signedProposal)
		if err != nil {
			return nil, errors.Wrap(err, "error with grpc ProcessProposal function")
		}

		proposalResponses = append(proposalResponses, proposalResponse)
	}
	return proposalResponses, nil
}
