package handlers

import (
	"reflect"

	"github.com/likecoin/likechain/abci/account"
	"github.com/likecoin/likechain/abci/context"
	"github.com/likecoin/likechain/abci/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func checkRegister(ctx context.Context, rawTx *types.Transaction) abci.ResponseCheckTx {
	tx := rawTx.GetRegisterTx()

	if !validateRegisterTransaction(tx) {
		return abci.ResponseCheckTx{
			Code: 1001,
			Info: "Invalid RegisterTransaction format",
		}
	}

	if !validateRegisterSignature(ctx, tx) {
		return abci.ResponseCheckTx{
			Code: 1002,
			Info: "Duplicated registration",
		}
	}

	return abci.ResponseCheckTx{Code: 0}
}

func deliverRegister(ctx context.Context, rawTx *types.Transaction) abci.ResponseDeliverTx {
	tx := rawTx.GetRegisterTx()

	if !validateRegisterTransaction(tx) {
		return abci.ResponseDeliverTx{
			Code: 1001,
			Info: "Invalid RegisterTransaction format",
		}
	}

	if !validateRegisterSignature(ctx, tx) {
		return abci.ResponseDeliverTx{
			Code: 1002,
			Info: "Duplicated registration",
		}
	}

	err := register(ctx, tx)
	if err {
		panic("Register error")
	}

	return abci.ResponseDeliverTx{Code: 0}
}

// validateRegisterSignature validates register transaction
func validateRegisterSignature(ctx context.Context, tx *types.RegisterTransaction) bool {
	return false // TODO
}

// validateRegisterTransaction validates register transaction
func validateRegisterTransaction(tx *types.RegisterTransaction) bool {
	return false // TODO
}

// register creates a new LikeChain account
func register(ctx context.Context, tx *types.RegisterTransaction) bool {
	err := true
	account.NewAccount(tx.Addr)
	return err // TODO
}

func init() {
	t := reflect.TypeOf((*types.Transaction_RegisterTx)(nil))
	registerCheckTxHandler(t, checkRegister)
	registerDeliverTxHandler(t, deliverRegister)
}
