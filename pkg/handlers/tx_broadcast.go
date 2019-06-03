package handlers

import (
	"fabric-rest-api-go/pkg/ca"
	"fabric-rest-api-go/pkg/tx"
	"fabric-rest-api-go/pkg/utils"
	protoutils "github.com/hyperledger/fabric/protos/utils"
	"github.com/labstack/echo/v4"
	"math/big"
	"net/http"
)

type TxBroadcastRequest struct {
	PayloadBytes string `json:"payload_bytes" validate:"required"`
	R            string `json:"r" validate:"required"`
	S            string `json:"s" validate:"required"`
}

type TxBroadcastResponse struct {
	PayloadBytes string `json:"payload_bytes"`
	PayloadHash  string `json:"payload_hash"`
}

func PostTxBroadcastHandler(ec echo.Context) error {
	c := ec.(*ApiContext)

	txBroadcastRequest := new(TxBroadcastRequest)
	if err := c.Bind(txBroadcastRequest); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(txBroadcastRequest); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, c.ValidationErrors(err).Error())
	}

	// create proposal from b64
	payloadBytes, err := utils.B64Decode(txBroadcastRequest.PayloadBytes)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	r := new(big.Int)
	r.SetString(txBroadcastRequest.R, 16)
	s := new(big.Int)
	s.SetString(txBroadcastRequest.S, 16)

	payloadSignature, err := ca.MarshalECDSASignature(r, ca.ToLowS_P256(s))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	payload, err := protoutils.UnmarshalPayload(payloadBytes)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	envelope, err := tx.CreateBroadcastEnvelope(payload, payloadSignature)
	if err != nil {
		return err
	}

	// TODO replace with ordered url
	target := "localhost:7050"
	// send to orderer
	err = tx.SendBroadcastToOrderer(envelope, target)
	if err != nil {
		return err
	}

	// TODO replace with TxId
	return ec.JSONBlob(http.StatusOK, []byte(`{"result": "ok"}`))
}
