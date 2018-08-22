package handlers

import (
	"reflect"

	"github.com/likecoin/likechain/abci/context"
	"github.com/likecoin/likechain/abci/response"
	"github.com/likecoin/likechain/abci/types"
)

func checkDeposit(state context.IImmutableState, rawTx *types.Transaction) response.R {
	tx := rawTx.GetDepositTx()
	if tx == nil {
		// TODO: log
		panic("Expect DepositTx but got nil")
	}

	_ = tx.BlockNumber

	if !validateDepositTransactionFormat(tx) {
		return response.DepositCheckTxInvalidFormat
	}

	return response.Success // TODO
}

func deliverDeposit(state context.IMutableState, rawTx *types.Transaction) response.R {
	tx := rawTx.GetDepositTx()
	if tx == nil {
		// TODO: log
		panic("Expect DepositTx but got nil")
	}

	_ = tx.BlockNumber

	if !validateDepositTransactionFormat(tx) {
		return response.DepositDeliverTxInvalidFormat
	}

	return response.Success // TODO
}

func validateDepositTransactionFormat(tx *types.DepositTransaction) bool {
	return false // TODO
}

func deposit(state context.IMutableState, tx *types.DepositTransaction) {
	// TODO
}

func init() {
	t := reflect.TypeOf((*types.Transaction_DepositTx)(nil))
	registerCheckTxHandler(t, checkDeposit)
	registerDeliverTxHandler(t, deliverDeposit)
}
