package handlers

import (
	"testing"

	"github.com/likecoin/likechain/abci/context"
	"github.com/likecoin/likechain/abci/types"
)

func TestCheckWithdraw(t *testing.T) {
	ctx := context.Context{}
	rawTx := &types.Transaction{}
	res := checkWithdraw(ctx, rawTx)
	t.Log(res)
	// TODO
}

func TestDeliverWithdraw(t *testing.T) {
	ctx := context.Context{}
	rawTx := &types.Transaction{}
	res := deliverWithdraw(ctx, rawTx)
	t.Log(res)
	// TODO
}

func TestValidateWithdrawTransaction(t *testing.T) {
	tx := &types.WithdrawTransaction{}
	if !validateWithdrawTransaction(tx) {
		t.Error("Validate WithdrawTransaction failed")
	}
}

func TestWithdraw(t *testing.T) {
	ctx := context.Context{}
	tx := &types.WithdrawTransaction{}
	withdraw(ctx, tx)
	// TODO
}
