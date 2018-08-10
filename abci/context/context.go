package context

import (
	"github.com/tendermint/iavl"
)

type Context interface {
	StateTree() *iavl.VersionedTree
	WithdrawTree() *iavl.VersionedTree
	// TODO
}
