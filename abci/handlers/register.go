package handlers

import (
	"fmt"
	"reflect"

	"github.com/likecoin/likechain/abci/account"
	"github.com/likecoin/likechain/abci/context"
	"github.com/likecoin/likechain/abci/errcode"
	"github.com/likecoin/likechain/abci/types"
	"github.com/likecoin/likechain/abci/utils"
	abci "github.com/tendermint/tendermint/abci/types"
)

func checkRegister(state context.IImmutableState, rawTx *types.Transaction) abci.ResponseCheckTx {
	tx := rawTx.GetRegisterTx()
	if tx == nil {
		// TODO: log
		panic("Expect RegisterTx but got nil")
	}

	if !validateRegisterTransaction(tx) {
		code, info := errcode.RegisterCheckTxInvalidFormat()
		return abci.ResponseCheckTx{
			Code: code,
			Info: info,
		}
	}

	if !validateRegisterSignature(state, tx) {
		code, info := errcode.RegisterCheckTxInvalidSignature()
		return abci.ResponseCheckTx{
			Code: code,
			Info: info,
		}
	}

	_, existed := state.ImmutableStateTree().Get(utils.DbAddrKey(tx.Addr.ToEthereum()))
	if existed != nil {
		code, info := errcode.RegisterCheckTxDuplicated()
		return abci.ResponseCheckTx{
			Code: code,
			Info: info,
		}
	}

	return abci.ResponseCheckTx{Code: 0}
}

func deliverRegister(state context.IMutableState, rawTx *types.Transaction) abci.ResponseDeliverTx {
	tx := rawTx.GetRegisterTx()
	if tx == nil {
		// TODO: log
		panic("Expect RegisterTx but got nil")
	}

	if !validateRegisterTransaction(tx) {
		code, info := errcode.RegisterDeliverTxInvalidFormat()
		return abci.ResponseDeliverTx{
			Code: code,
			Info: info,
		}
	}

	if !validateRegisterSignature(state, tx) {
		code, info := errcode.RegisterDeliverTxInvalidSignature()
		return abci.ResponseDeliverTx{
			Code: code,
			Info: info,
		}
	}

	_, existed := state.ImmutableStateTree().Get(utils.DbAddrKey(tx.Addr.ToEthereum()))
	if existed != nil {
		code, info := errcode.RegisterDeliverTxDuplicated()
		return abci.ResponseDeliverTx{
			Code: code,
			Info: info,
		}
	}

	id, err := register(state, tx)
	if err != nil {
		panic(fmt.Sprintf("Error occurs during registration, details: %v", err))
	}

	return abci.ResponseDeliverTx{
		Code: 0,
		Data: id.Content,
	}
}

// validateRegisterSignature validates register transaction
func validateRegisterSignature(state context.IImmutableState, tx *types.RegisterTransaction) bool {
	hashedMsg, err := tx.GenerateSigningMessageHash()
	if err != nil {
		// TODO: log
		return false
	}

	sigAddr, err := utils.RecoverSignature(hashedMsg, tx.Sig)
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
	return tx.Addr.IsValidFormat() && tx.Sig.IsValidFormat()
}

// register creates a new LikeChain account
func register(state context.IMutableState, tx *types.RegisterTransaction) (types.LikeChainID, error) {
	ethAddr := tx.Addr.ToEthereum()
	return account.NewAccount(state, ethAddr)
}

func init() {
	t := reflect.TypeOf((*types.Transaction_RegisterTx)(nil))
	registerCheckTxHandler(t, checkRegister)
	registerDeliverTxHandler(t, deliverRegister)
}
