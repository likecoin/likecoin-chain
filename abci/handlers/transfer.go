package handlers

import (
	"reflect"

	"github.com/likecoin/likechain/abci/account"
	"github.com/likecoin/likechain/abci/context"
	"github.com/likecoin/likechain/abci/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func checkTransfer(ctx context.Context, rawTx *types.Transaction) abci.ResponseCheckTx {
	tx := rawTx.GetTransferTx()
	_ = tx.From
	return abci.ResponseCheckTx{} // TODO
}

func deliverTransfer(ctx context.Context, rawTx *types.Transaction) abci.ResponseDeliverTx {
	tx := rawTx.GetTransferTx()

	if !validateTransferSignature(tx.Sig) {
		panic("Invalid signature")
	}

	if !validateTransferTransaction(tx) {
		panic("Invalid TransferTransaction in TransferTx")
	}

	_ = account.FetchBalance(context, tx.From)
	_ = account.FetchNextNonce(context, tx.From)
	// Increment nonce
	// Adjust balance of sender and receiver

	return abci.ResponseDeliverTx{} // TODO
}

func validateTransferSignature(sig *types.Signature) bool {
	return false // TODO
}

func validateTransferTransaction(tx *types.TransferTransaction) bool {
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
