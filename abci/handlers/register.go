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

	if validateRegisterTransaction(tx) {
		panic("Invalid RegisterTransaction in CheckTx")
	}

	_ = tx.Addr

	return abci.ResponseCheckTx{} // TODO
}

func deliverRegister(ctx context.Context, rawTx *types.Transaction) abci.ResponseDeliverTx {
	tx := rawTx.GetRegisterTx()

	if !validateRegisterSignature(tx.Sig) {
		panic("Invalid signature")
	}

	if !validateRegisterTransaction(tx) {
		panic("Invalid RegisterTransaction in DeliverTx")
	}

	err := register(context, tx)
	if err {
		panic("Register error")
	}

	return abci.ResponseDeliverTx{} // TODO
}

// validateRegisterSignature validates register transaction
func validateRegisterSignature(sig *types.Signature) bool {
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
