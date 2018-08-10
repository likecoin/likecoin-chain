package handlers

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/likecoin/likechain/abci/context"
	"github.com/likecoin/likechain/abci/types"
)

func TestCheckDeposit(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	ctx := context.NewMockContext(mockCtrl)
	// TODO: mock ctx calls

	rawTx := &types.Transaction{}
	res := checkDeposit(ctx, rawTx)
	t.Log(res)
	// TODO
}

func TestDeliverDeposit(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	ctx := context.NewMockContext(mockCtrl)
	// TODO: mock ctx calls

	rawTx := &types.Transaction{}
	res := deliverDeposit(ctx, rawTx)
	t.Log(res)
	// TODO
}

func TestValidateDepositTransaction(t *testing.T) {
	tx := &types.DepositTransaction{}
	if !validateDepositTransaction(tx) {
		t.Error("Validate DepositTransaction failed")
	}
}

func TestDeposit(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	ctx := context.NewMockContext(mockCtrl)
	// TODO: mock ctx calls

	tx := &types.DepositTransaction{}
	deposit(ctx, tx)
	// TODO
}
