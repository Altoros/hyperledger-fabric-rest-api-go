package handlers

import (
	"fabric-rest-api-go/pkg/ca"
	"fabric-rest-api-go/pkg/tx"
	"fabric-rest-api-go/pkg/utils"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric/protos/peer"
	protoutils "github.com/hyperledger/fabric/protos/utils"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"math/big"
	"net/http"
)

type TxBroadcastPayloadRequest struct {
	ProposalBytes string `json:"proposal_bytes" validate:"required"`
	R             string `json:"r" validate:"required"`
	S             string `json:"s" validate:"required"`
}

type TxBroadcastPayloadResponse struct {
	PayloadBytes string `json:"payload_bytes"`
	PayloadHash  string `json:"payload_hash"`
}

func PostTxBroadcastPayloadHandler(ec echo.Context) error {
	c := ec.(*ApiContext)

	txBroadcastPayloadRequest := new(TxBroadcastPayloadRequest)
	if err := c.Bind(txBroadcastPayloadRequest); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(txBroadcastPayloadRequest); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, c.ValidationErrors(err).Error())
	}

	// TODO change to one peer (one for tx/query, multiple for tx/invoke)
	targets := []string{"localhost:7051"}

	proposalBytes, err := utils.B64Decode(txBroadcastPayloadRequest.ProposalBytes)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	proposal, err := protoutils.GetProposal(proposalBytes)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	r := new(big.Int)
	r.SetString(txBroadcastPayloadRequest.R, 16)
	s := new(big.Int)
	s.SetString(txBroadcastPayloadRequest.S, 16)

	proposalSignature, err := ca.MarshalECDSASignature(r, ca.ToLowS_P256(s))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	signedProposal := &peer.SignedProposal{ProposalBytes: proposalBytes, Signature: proposalSignature}

	proposalResponses, err := tx.SendProposalToPeersEndorsment(signedProposal, targets)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	payload, err := tx.CreateBroadcastPayload(proposal, proposalResponses)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, errors.Wrap(err, "creating of broadcast payload failed").Error())
	}

	payloadBytes, err := proto.Marshal(payload)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	payloadHash, err := tx.BroadcastPayloadHash(payload)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ec.JSON(http.StatusOK, TxBroadcastPayloadResponse{
		PayloadBytes: utils.B64Encode(payloadBytes),
		PayloadHash:  fmt.Sprintf("%x", payloadHash),
	})
}
