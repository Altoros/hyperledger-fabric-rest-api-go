package handlers

import (
	"fabric-rest-api-go/pkg/ca"
	"fabric-rest-api-go/pkg/sdk"
	"fabric-rest-api-go/pkg/tx"
	"fabric-rest-api-go/pkg/utils"
	"github.com/hyperledger/fabric/protos/peer"
	"github.com/labstack/echo/v4"
	"math/big"
	"net/http"
)

type TxQueryRequest struct {
	ProposalBytes string `json:"proposal_bytes" validate:"required"`
	R             string `json:"r" validate:"required"`
	S             string `json:"s" validate:"required"`
	Peer          string `json:"peer" validate:"required"`
}

type TxQueryResponse struct {
	Result string `json:"result"`
}

func PostTxQueryHandler(ec echo.Context) error {
	c := ec.(*ApiContext)

	txQueryRequest := new(TxQueryRequest)
	if err := c.Bind(txQueryRequest); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(txQueryRequest); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, c.ValidationErrors(err).Error())
	}

	peers, err := c.ParseApiPeers([]string{txQueryRequest.Peer})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	targets := []sdk.ApiPeer{peers[0]}

	// peer.proposal bytes from b64
	proposalBytes, err := utils.B64Decode(txQueryRequest.ProposalBytes)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// TODO wrap every signature marshaling in project to one method
	r := new(big.Int)
	r.SetString(txQueryRequest.R, 16)
	s := new(big.Int)
	s.SetString(txQueryRequest.S, 16)

	proposalSignature, err := ca.MarshalECDSASignature(r, ca.ToLowS_P256(s))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	signedProposal := &peer.SignedProposal{ProposalBytes: proposalBytes, Signature: proposalSignature}

	proposalResponses, err := tx.SendProposalToPeersEndorsment(signedProposal, targets)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ec.JSON(http.StatusOK, TxQueryResponse{Result: string(proposalResponses[0].Response.Payload)})
}
