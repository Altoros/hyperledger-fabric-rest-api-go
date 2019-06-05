package tx

import (
	"context"
	"crypto/sha256"
	"fabric-rest-api-go/pkg/sdk"
	"github.com/hyperledger/fabric/protos/common"
	"github.com/hyperledger/fabric/protos/orderer"
	"github.com/hyperledger/fabric/protos/peer"
	"github.com/hyperledger/fabric/protos/utils"
	"github.com/pkg/errors"
	"io"
	"strings"
)

func CreateBroadcastPayload(proposal *peer.Proposal, proposalResponses []*peer.ProposalResponse) (*common.Payload, error) {
	// form transaction with endorsements
	transaction, err := NewTransaction(proposal, proposalResponses)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create transaction")
	}

	// the original header
	hdr, err := utils.GetHeader(proposal.Header)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshal proposal header failed")
	}

	// serialize the tx
	txBytes, err := utils.GetBytesTransaction(transaction)
	if err != nil {
		return nil, err
	}

	// create the payload
	payload := &common.Payload{Header: hdr, Data: txBytes}

	return payload, nil
}

func BroadcastPayloadHash(payload *common.Payload) ([]byte, error) {
	payloadBytes, err := utils.GetBytesPayload(payload)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get payloadBytes")
	}

	hasher := sha256.New()
	hasher.Write(payloadBytes)
	proposalSha256 := hasher.Sum(nil)
	return proposalSha256, nil
}

// New create a transaction with proposal response, following the endorsement policy.
func NewTransaction(proposal *peer.Proposal, proposalResponses []*peer.ProposalResponse) (*peer.Transaction, error) {
	if len(proposalResponses) == 0 {
		return nil, errors.New("at least one proposal response is necessary")
	}

	// the original header
	hdr, err := utils.GetHeader(proposal.Header)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshal proposal header failed")
	}

	// the original payload
	pPayl, err := utils.GetChaincodeProposalPayload(proposal.Payload)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshal proposal payload failed")
	}

	// get header extensions so we have the visibility field
	hdrExt, err := utils.GetChaincodeHeaderExtension(hdr)
	if err != nil {
		return nil, err
	}

	responsePayload := proposalResponses[0].Payload
	//if err := validateProposalResponses(request.ProposalResponses); err != nil {
	//	return nil, err
	//}

	// fill endorsements
	endorsements := make([]*peer.Endorsement, len(proposalResponses))
	for n, r := range proposalResponses {
		endorsements[n] = r.Endorsement
	}

	// create ChaincodeEndorsedAction
	cea := &peer.ChaincodeEndorsedAction{ProposalResponsePayload: responsePayload, Endorsements: endorsements}

	// obtain the bytes of the proposal payload that will go to the transaction
	proposalPayloadBytes, err := utils.GetBytesProposalPayloadForTx(pPayl, hdrExt.PayloadVisibility)
	if err != nil {
		return nil, err
	}

	// serialize the chaincode action payload
	chaincodeActionPayload := &peer.ChaincodeActionPayload{ChaincodeProposalPayload: proposalPayloadBytes, Action: cea}
	chaincodeActionPayloadBytes, err := utils.GetBytesChaincodeActionPayload(chaincodeActionPayload)
	if err != nil {
		return nil, err
	}

	// create a transaction
	taa := &peer.TransactionAction{Header: hdr.SignatureHeader, Payload: chaincodeActionPayloadBytes}
	taas := make([]*peer.TransactionAction, 1)
	taas[0] = taa

	return &peer.Transaction{Actions: taas}, nil
}

func CreateBroadcastEnvelope(payload *common.Payload, signature []byte) (*common.Envelope, error) {
	payloadBytes, err := utils.GetBytesPayload(payload)
	if err != nil {
		return nil, err
	}

	envelope := &common.Envelope{
		Payload:   payloadBytes,
		Signature: signature,
	}
	return envelope, nil
}

// SendBroadcastToOrderer Send the created transaction to Orderer.
func SendBroadcastToOrderer(envelope *common.Envelope, apiOrderer sdk.ApiOrderer) error {
	ctx := context.TODO()
	conn, err := apiOrderer.GrpcConn(ctx)
	if err != nil {
		return err
	}

	broadcastClient, err := orderer.NewAtomicBroadcastClient(conn).Broadcast(ctx)
	if err != nil {
		return errors.Wrap(err, "NewAtomicBroadcastClient failed")
	}

	responses := make(chan common.Status)
	errs := make(chan error, 1)

	// listener of broadcast responses
	go broadcastStream(broadcastClient, responses, errs)

	err = broadcastClient.Send(envelope)
	if err != nil {
		return errors.Wrap(err, "failed to send envelope to orderer")
	}
	if err = broadcastClient.CloseSend(); err != nil {
		return errors.Wrap(err, "unable to close broadcast client [%s]")
	}

	status, err := wrapStreamStatusRPC(responses, errs)

	if err != nil {
		return errors.Wrapf(err, "broadcast error, status %s", status)
	}

	return nil
}

// wrapStreamStatusRPC returns the last response and err and blocks until the chan is closed.
func wrapStreamStatusRPC(responses chan common.Status, errs chan error) (common.Status, error) {
	var status common.Status
	var errorsAr []string

read:
	for {
		select {
		case s, ok := <-responses:
			if !ok {
				break read
			}
			status = s
		case e := <-errs:
			errorsAr = append(errorsAr, e.Error())
		}
	}

	for i := 0; i < len(errs); i++ {
		e := <-errs
		errorsAr = append(errorsAr, e.Error())
	}

	if len(errorsAr) > 0 {
		return status, errors.New(strings.Join(errorsAr, ", "))
	} else {
		return status, nil
	}
}

func broadcastStream(broadcastClient orderer.AtomicBroadcast_BroadcastClient, responses chan common.Status, errs chan error) {
	for {
		broadcastResponse, err := broadcastClient.Recv()
		if err == io.EOF {
			// done
			close(responses)
			return
		}

		if err != nil {
			errs <- errors.Wrap(err, "broadcast recv failed")
			close(responses)
			return
		}

		if broadcastResponse.Status == common.Status_SUCCESS {
			responses <- broadcastResponse.Status
		} else {
			errs <- errors.Errorf("broadcast error %s, %s", broadcastResponse.Status, broadcastResponse.Info)
		}
	}
}
