package handlers

import (
	"reflect"

	"github.com/likecoin/likechain/abci/account"
	"github.com/likecoin/likechain/abci/context"
	"github.com/likecoin/likechain/abci/error"
	"github.com/likecoin/likechain/abci/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func checkTransfer(ctx context.Context, rawTx *types.Transaction) abci.ResponseCheckTx {
	tx := rawTx.GetTransferTx()

	if !validateTransferTransactionFormat(tx) {
		code, info := error.TransferCheckTxInvalidFormat()
		return abci.ResponseCheckTx{
			Code: code,
			Info: info,
		}
	}

	if !validateTransferSignature(tx.Sig) {
		code, info := error.TransferCheckTxInvalidSignature()
		return abci.ResponseCheckTx{
			Code: code,
			Info: info,
		}
	}

	return abci.ResponseCheckTx{} // TODO
}

func deliverTransfer(ctx context.Context, rawTx *types.Transaction) abci.ResponseDeliverTx {
	tx := rawTx.GetTransferTx()

	if !validateTransferTransactionFormat(tx) {
		code, info := error.TransferDeliverTxInvalidFormat()
		return abci.ResponseDeliverTx{
			Code: code,
			Info: info,
		}
	}

	if !validateTransferSignature(tx.Sig) {
		code, info := error.TransferDeliverTxInvalidSignature()
		return abci.ResponseDeliverTx{
			Code: code,
			Info: info,
		}
	}

	_ = account.FetchBalance(ctx, tx.From)
	_ = account.FetchNextNonce(ctx, tx.From)
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

func transfer(ctx context.Context, tx *types.TransferTransaction) {
	// TODO
}

func init() {
	t := reflect.TypeOf((*types.Transaction_TransferTx)(nil))
	registerCheckTxHandler(t, checkTransfer)
	registerDeliverTxHandler(t, deliverTransfer)
}
