package handlers

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/likecoin/likechain/abci/context"
	"github.com/likecoin/likechain/abci/types"
)

func TestCheckWithdraw(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	ctx := context.NewMockContext(mockCtrl)
	// TODO: mock ctx calls

	rawTx := &types.Transaction{}
	res := checkWithdraw(ctx, rawTx)
	t.Log(res)
	// TODO
}

func TestDeliverWithdraw(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	ctx := context.NewMockContext(mockCtrl)
	// TODO: mock ctx calls

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
	mockCtrl := gomock.NewController(t)
	ctx := context.NewMockContext(mockCtrl)
	// TODO: mock ctx calls

	tx := &types.WithdrawTransaction{}
	withdraw(ctx, tx)
	// TODO
}
