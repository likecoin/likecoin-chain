package context

import (
	"github.com/tendermint/iavl"
)

type ImmutableContext interface {
	StateTree() *iavl.ImmutableTree
	WithdrawTree() *iavl.ImmutableTree
	GetBlockHash() []byte
}

type MutableContext interface {
	ImmutableContext
	MutableStateTree() *iavl.MutableTree
	MutableWithdrawTree() *iavl.MutableTree
}

var cacheSize = 1048576
