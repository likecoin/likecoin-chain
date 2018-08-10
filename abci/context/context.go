package context

import (
	"github.com/tendermint/iavl"
)

type Context struct {
	StateTree    *iavl.VersionedTree
	WithdrawTree *iavl.VersionedTree
	// TODO
}
