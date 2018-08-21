package handlers

import (
	"reflect"

	"github.com/likecoin/likechain/abci/account"
	"github.com/likecoin/likechain/abci/context"
	"github.com/likecoin/likechain/abci/errcode"
	"github.com/likecoin/likechain/abci/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func checkTransfer(state context.IImmutableState, rawTx *types.Transaction) abci.ResponseCheckTx {
	tx := rawTx.GetTransferTx()
	if tx == nil {
		// TODO: log
		panic("Expect TransferTx but got nil")
	}

	if !validateTransferTransactionFormat(tx) {
		code, info := errcode.TransferCheckTxInvalidFormat()
		return abci.ResponseCheckTx{
			Code: code,
			Info: info,
		}
	}

	if !validateTransferSignature(tx.Sig) {
		code, info := errcode.TransferCheckTxInvalidSignature()
		return abci.ResponseCheckTx{
			Code: code,
			Info: info,
		}
	}

	return abci.ResponseCheckTx{} // TODO
}

func deliverTransfer(state context.IMutableState, rawTx *types.Transaction) abci.ResponseDeliverTx {
	tx := rawTx.GetTransferTx()
	if tx == nil {
		// TODO: log
		panic("Expect TransferTx but got nil")
	}

	if !validateTransferTransactionFormat(tx) {
		code, info := errcode.TransferDeliverTxInvalidFormat()
		return abci.ResponseDeliverTx{
			Code: code,
			Info: info,
		}
	}

	if !validateTransferSignature(tx.Sig) {
		code, info := errcode.TransferDeliverTxInvalidSignature()
		return abci.ResponseDeliverTx{
			Code: code,
			Info: info,
		}
	}

	fromID, exist := account.GetLikeChainID(state, *tx.From)
	if !exist {
		return abci.ResponseDeliverTx{} // TODO: error code for sender account does not exist
	}

	_ = account.FetchBalance(state, fromID)
	_ = account.FetchNextNonce(state, fromID)
	// Increment nonce
	// Adjust balance of sender and receiver

	return abci.ResponseDeliverTx{} // TODO
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
