package handlers

import (
	"fmt"
	"reflect"

	"github.com/likecoin/likechain/abci/account"
	"github.com/likecoin/likechain/abci/context"
	"github.com/likecoin/likechain/abci/response"
	"github.com/likecoin/likechain/abci/types"
	"github.com/likecoin/likechain/abci/utils"
)

func checkRegister(state context.IImmutableState, rawTx *types.Transaction) response.R {
	tx := rawTx.GetRegisterTx()
	if tx == nil {
		// TODO: log
		panic("Expect RegisterTx but got nil")
	}

	if !validateRegisterTransaction(tx) {
		return response.RegisterCheckTxInvalidFormat
	}

	if !validateRegisterSignature(state, tx) {
		return response.RegisterCheckTxInvalidSignature
	}

	_, existed := state.ImmutableStateTree().Get(utils.DbAddrKey(tx.Addr.ToEthereum()))
	if existed != nil {
		return response.RegisterCheckTxDuplicated
	}

	return response.Success
}

func deliverRegister(state context.IMutableState, rawTx *types.Transaction) response.R {
	tx := rawTx.GetRegisterTx()
	if tx == nil {
		// TODO: log
		panic("Expect RegisterTx but got nil")
	}

	if !validateRegisterTransaction(tx) {
		return response.RegisterDeliverTxInvalidFormat
	}

	if !validateRegisterSignature(state, tx) {
		return response.RegisterDeliverTxInvalidSignature
	}

	_, existed := state.ImmutableStateTree().Get(utils.DbAddrKey(tx.Addr.ToEthereum()))
	if existed != nil {
		return response.RegisterDeliverTxDuplicated
	}

	id, err := register(state, tx)
	if err != nil {
		panic(fmt.Sprintf("Error occurs during registration, details: %v", err))
	}

	return response.Success.Merge(response.R{
		Data: id.Content,
	})
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
