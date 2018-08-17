package context

import (
	"github.com/tendermint/iavl"
	"github.com/tendermint/tendermint/libs/db"
)

type LikeContextMock struct {
	*DeliverTxContext
}

func NewMock() *LikeContextMock {
	return &LikeContextMock{
		DeliverTxContext: &DeliverTxContext{
			stateTree:    iavl.NewMutableTree(db.NewMemDB(), 0),
			withdrawTree: iavl.NewMutableTree(db.NewMemDB(), 0),
			blockHash:    []byte{1, 3, 3, 7},
		},
	}
}

func (ctx *DeliverTxContext) Reset() {
	ctx.MutableStateTree().Rollback()
	ctx.MutableWithdrawTree().Rollback()
}
