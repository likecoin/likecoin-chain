package handlers

import (
	"reflect"

	"github.com/likecoin/likechain/abci/account"
	"github.com/likecoin/likechain/abci/context"
	"github.com/likecoin/likechain/abci/error"
	"github.com/likecoin/likechain/abci/types"
	"github.com/likecoin/likechain/abci/utils"
	abci "github.com/tendermint/tendermint/abci/types"
)

func checkRegister(ctx context.Context, rawTx *types.Transaction) abci.ResponseCheckTx {
	tx := rawTx.GetRegisterTx()
	if tx == nil {
		// TODO: log
		panic("Expect RegisterTx but got nil")
	}

	if !validateRegisterTransaction(tx) {
		code, info := error.RegisterCheckTxInvalidFormat()
		return abci.ResponseCheckTx{
			Code: code,
			Info: info,
		}
	}

	if !validateRegisterSignature(ctx, tx) {
		code, info := error.RegisterCheckTxInvalidSignature()
		return abci.ResponseCheckTx{
			Code: code,
			Info: info,
		}
	}

	return abci.ResponseCheckTx{Code: 0}
}

func deliverRegister(ctx context.Context, rawTx *types.Transaction) abci.ResponseDeliverTx {
	tx := rawTx.GetRegisterTx()
	if tx == nil {
		// TODO: log
		panic("Expect RegisterTx but got nil")
	}

	if !validateRegisterTransaction(tx) {
		code, info := error.RegisterDeliverTxInvalidFormat()
		return abci.ResponseDeliverTx{
			Code: code,
			Info: info,
		}
	}

	if !validateRegisterSignature(ctx, tx) {
		code, info := error.RegisterDeliverTxInvalidSignature()
		return abci.ResponseDeliverTx{
			Code: code,
			Info: info,
		}
	}

	addr := tx.Addr.ToEthereum()
	_ = addr // TODO: check if address is already registered

	err := register(ctx, tx)
	if err {
		panic("Register error")
	}

	return abci.ResponseDeliverTx{Code: 0}
}

// validateRegisterSignature validates register transaction
func validateRegisterSignature(ctx context.Context, tx *types.RegisterTransaction) bool {
	hashedMsg, err := tx.GenerateSigningMessageHash()
	if err != nil {
		// TODO: log
		return false
	}

	sigAddr, err := utils.RecoverSignature(hashedMsg, *tx.Sig)
	if err != nil {
		return false
	}

	if tx.Addr.ToEthereum() != sigAddr {
		return false
	}

	return true
}

// validateRegisterTransaction validates register transaction
func validateRegisterTransaction(tx *types.RegisterTransaction) bool {
	return tx.Sig.IsValidFormat() && tx.Addr.IsValidFormat()
}

// register creates a new LikeChain account
func register(ctx context.Context, tx *types.RegisterTransaction) types.LikeChainID {
	err := true
	ethAddr := tx.Addr.ToEthereum()
	return account.NewAccount(ethAddr)
}

func init() {
	t := reflect.TypeOf((*types.Transaction_RegisterTx)(nil))
	registerCheckTxHandler(t, checkRegister)
	registerDeliverTxHandler(t, deliverRegister)
}
