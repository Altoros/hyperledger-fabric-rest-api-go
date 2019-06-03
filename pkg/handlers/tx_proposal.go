package handlers

import (
	"fabric-rest-api-go/pkg/tx"
	"fabric-rest-api-go/pkg/utils"
	"fmt"
	protoutils "github.com/hyperledger/fabric/protos/utils"
	"github.com/labstack/echo/v4"
	"net/http"
)

type ProposalRequest struct {
	ChaincodeId string   `json:"chaincode_id" validate:"required"`
	ChannelId   string   `json:"channel_id" validate:"required"`
	Fcn         string   `json:"fcn" validate:"required"`
	Args        []string `json:"args"`
	MspId       string   `json:"msp_id" validate:"required"`
	UserCert    string   `json:"user_cert" validate:"required"`
}

type ProposalResponse struct {
	ProposalBytes string `json:"proposal_bytes"`
	ProposalHash  string `json:"proposal_hash"`
}

func PostTxProposalHandler(ec echo.Context) error {
	c := ec.(*ApiContext)

	proposalRequest := new(ProposalRequest)
	if err := c.Bind(proposalRequest); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(proposalRequest); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, c.ValidationErrors(err).Error())
	}

	argsArray := [][]byte{[]byte(proposalRequest.Fcn)}
	for _, arg := range proposalRequest.Args {
		argsArray = append(argsArray, []byte(arg))
	}

	proposal, err := tx.CreateProposal(
		proposalRequest.ChannelId,
		proposalRequest.ChaincodeId,
		argsArray,
		proposalRequest.MspId,
		[]byte(proposalRequest.UserCert),
	)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	proposalBytes, err := protoutils.GetBytesProposal(proposal)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	proposalHash, err := tx.ProposalHash(proposal)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ec.JSON(http.StatusOK, ProposalResponse{
		ProposalBytes: utils.B64Encode(proposalBytes),
		ProposalHash:  fmt.Sprintf("%x", proposalHash),
	})

}
