package context

import (
	"github.com/tendermint/iavl"
	"github.com/tendermint/tendermint/libs/db"
)

type ContextMock struct {
	stateTree    *iavl.MutableTree
	withdrawTree *iavl.MutableTree
}

func NewMock() *ContextMock {
	stateTree := iavl.NewMutableTree(db.NewMemDB(), 128)
	withdrawTree := iavl.NewMutableTree(db.NewMemDB(), 128)
	return &ContextMock{stateTree, withdrawTree}
}

func (ctx *ContextMock) Reset() {
	// TODO: clear versions
	ctx.stateTree.Rollback()
	ctx.withdrawTree.Rollback()
}

func (ctx *ContextMock) StateTree() *iavl.MutableTree {
	return ctx.stateTree
}

func (ctx *ContextMock) WithdrawTree() *iavl.MutableTree {
	return ctx.withdrawTree
}

// GetBlockHash returns current block hash in state
func (ctx *ContextMock) GetBlockHash() []byte {
	return nil // TODO
}
