package handlers

import (
	"testing"

	"github.com/likecoin/likechain/abci/context"
	"github.com/likecoin/likechain/abci/types"
)

func TestCheckRegister(t *testing.T) {
	ctx := context.Context{}
	rawTx := &types.Transaction{}
	res := checkRegister(ctx, rawTx)
	t.Log(res)
	// TODO
}

func TestDeliverRegister(t *testing.T) {
	ctx := context.Context{}
	rawTx := &types.Transaction{}
	res := deliverRegister(ctx, rawTx)
	t.Log(res)
	// TODO
}

func TestValidateRegisterSignature(t *testing.T) {
	tx := &types.RegisterTransaction{}
	if !validateRegisterSignature(tx.Sig) {
		t.Error("Validate RegisterSignature failed")
	}
}

func TestValidateRegisterTransaction(t *testing.T) {
	tx := &types.RegisterTransaction{}
	if !validateRegisterTransaction(tx) {
		t.Error("Validate RegisterTransaction failed")
	}
	// TODO
}

func TestRegister(t *testing.T) {
	ctx := context.Context{}
	tx := &types.RegisterTransaction{}
	register(ctx, tx)
	// TODO
}
