package context

import (
	"github.com/tendermint/iavl"
	"github.com/tendermint/tendermint/libs/db"
)

type Context interface {
	StateTree() *iavl.MutableTree
	WithdrawTree() *iavl.MutableTree
	// TODO
}

type ContextImpl struct {
	stateTree    *iavl.MutableTree
	widhdrawTree *iavl.MutableTree
}

func New(db db.DB) {
}
