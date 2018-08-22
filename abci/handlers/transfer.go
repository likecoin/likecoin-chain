package handlers

import (
	"reflect"

	"github.com/likecoin/likechain/abci/account"
	"github.com/likecoin/likechain/abci/context"
	"github.com/likecoin/likechain/abci/response"
	"github.com/likecoin/likechain/abci/types"
)

func checkTransfer(state context.IImmutableState, rawTx *types.Transaction) response.R {
	tx := rawTx.GetTransferTx()
	if tx == nil {
		// TODO: log
		panic("Expect TransferTx but got nil")
	}

	if !validateTransferTransactionFormat(tx) {
		return response.TransferCheckTxInvalidFormat
	}

	if !validateTransferSignature(tx.Sig) {
		return response.TransferCheckTxInvalidSignature
	}

	return response.Success // TODO
}

func deliverTransfer(state context.IMutableState, rawTx *types.Transaction) response.R {
	tx := rawTx.GetTransferTx()
	if tx == nil {
		// TODO: log
		panic("Expect TransferTx but got nil")
	}

	if !validateTransferTransactionFormat(tx) {
		return response.TransferDeliverTxInvalidFormat
	}

	if !validateTransferSignature(tx.Sig) {
		return response.TransferDeliverTxInvalidSignature
	}

	fromID, exist := account.GetLikeChainID(state, *tx.From)
	if !exist {
		return response.Success // TODO: error code for sender account does not exist
	}

	_ = account.FetchBalance(state, fromID)
	_ = account.FetchNextNonce(state, fromID)
	// Increment nonce
	// Adjust balance of sender and receiver

	return response.Success // TODO
}

func validateTransferSignature(sig *types.Signature) bool {
	return false // TODO
}

func validateTransferTransactionFormat(tx *types.TransferTransaction) bool {
	return false // TODO
}

func transfer(state context.IMutableState, tx *types.TransferTransaction) {
	// TODO
}

func init() {
	t := reflect.TypeOf((*types.Transaction_TransferTx)(nil))
	registerCheckTxHandler(t, checkTransfer)
	registerDeliverTxHandler(t, deliverTransfer)
}
