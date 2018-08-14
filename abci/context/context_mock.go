package context

import (
	"github.com/tendermint/iavl"
	"github.com/tendermint/tmlibs/db" // TODO: change to "github.com/tendermint/tendermint/libs/db" after iavl update
)

type ContextMock struct {
	stateTree    *iavl.VersionedTree
	withdrawTree *iavl.VersionedTree
}

func NewMock() *ContextMock {
	memdb := db.NewMemDB()
	stateTree := iavl.NewVersionedTree(memdb, 128)
	withdrawTree := iavl.NewVersionedTree(memdb, 128)
	return &ContextMock{stateTree, withdrawTree}
}

func (ctx *ContextMock) Reset() {
	// TODO: clear versions
	ctx.stateTree.Rollback()
	ctx.withdrawTree.Rollback()
}

func (ctx *ContextMock) StateTree() *iavl.VersionedTree {
	return ctx.stateTree
}

func (ctx *ContextMock) WithdrawTree() *iavl.VersionedTree {
	return ctx.withdrawTree
}
