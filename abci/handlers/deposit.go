package handlers

import (
	"reflect"

	"github.com/likecoin/likechain/abci/context"
	"github.com/likecoin/likechain/abci/errcode"
	"github.com/likecoin/likechain/abci/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func checkDeposit(ctx context.ImmutableContext, rawTx *types.Transaction) abci.ResponseCheckTx {
	tx := rawTx.GetDepositTx()
	if tx == nil {
		// TODO: log
		panic("Expect DepositTx but got nil")
	}

	_ = tx.BlockNumber

	if !validateDepositTransactionFormat(tx) {
		code, info := errcode.DepositCheckTxInvalidFormat()
		return abci.ResponseCheckTx{
			Code: code,
			Info: info,
		}
	}

	return abci.ResponseCheckTx{} // TODO
}

func deliverDeposit(ctx context.MutableContext, rawTx *types.Transaction) abci.ResponseDeliverTx {
	tx := rawTx.GetDepositTx()
	if tx == nil {
		// TODO: log
		panic("Expect DepositTx but got nil")
	}

	_ = tx.BlockNumber

	if !validateDepositTransactionFormat(tx) {
		code, info := errcode.DepositDeliverTxInvalidFormat()
		return abci.ResponseDeliverTx{
			Code: code,
			Info: info,
		}
	}

	return abci.ResponseDeliverTx{} // TODO
}

func validateDepositTransactionFormat(tx *types.DepositTransaction) bool {
	return false // TODO
}

func deposit(ctx context.MutableContext, tx *types.DepositTransaction) {
	// TODO
}

func init() {
	t := reflect.TypeOf((*types.Transaction_DepositTx)(nil))
	registerCheckTxHandler(t, checkDeposit)
	registerDeliverTxHandler(t, deliverDeposit)
}
