package handlers

import (
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/likecoin/likechain/abci/context"
	"github.com/likecoin/likechain/abci/types"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCheckRegister(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	ctx := context.NewMockContext(mockCtrl)
	// TODO: mock ctx calls

	rawTx := &types.Transaction{}
	res := checkRegister(ctx, rawTx)
	t.Log(res)
	// TODO
}

func TestDeliverRegister(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	ctx := context.NewMockContext(mockCtrl)
	// TODO: mock ctx calls

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
	Convey("Given a valid Register transaction", t, func() {
		tx := &types.RegisterTransaction{} // TODO
		Convey("The Register transaction should pass the validation", func() {
			So(validateRegisterTransaction(tx), ShouldBeTrue)
		})
	})
}

func TestRegister(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	ctx := context.NewMockContext(mockCtrl)
	// TODO: mock ctx calls

	tx := &types.RegisterTransaction{}
	register(ctx, tx)
	// TODO
}
