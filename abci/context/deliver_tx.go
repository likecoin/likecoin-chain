package context

import (
	"github.com/tendermint/iavl"
	"github.com/tendermint/tendermint/libs/db"
)

type DeliverTxContext struct {
	stateTree    *iavl.MutableTree
	withdrawTree *iavl.MutableTree
	blockHash    []byte
}

func newTree(dbPath, dir string) *iavl.MutableTree {
	db, err := db.NewGoLevelDB(dbPath, dir)
	if err != nil {
		panic(err)
	}
	tree := iavl.NewMutableTree(db, cacheSize)
	return tree
}

func New(dbPath string) *DeliverTxContext {
	stateTree := newTree(dbPath, "state")
	withdrawTree := newTree(dbPath, "withdraw")
	return &DeliverTxContext{
		stateTree:    stateTree,
		withdrawTree: withdrawTree,
		blockHash:    nil,
	}
}

func (ctx *DeliverTxContext) StateTree() *iavl.ImmutableTree {
	return ctx.stateTree.ImmutableTree
}

func (ctx *DeliverTxContext) WithdrawTree() *iavl.ImmutableTree {
	return ctx.withdrawTree.ImmutableTree
}

func (ctx *DeliverTxContext) MutableStateTree() *iavl.MutableTree {
	return ctx.stateTree
}

func (ctx *DeliverTxContext) MutableWithdrawTree() *iavl.MutableTree {
	return ctx.withdrawTree
}

func (ctx *DeliverTxContext) GetBlockHash() []byte {
	return ctx.blockHash
}

func (ctx *DeliverTxContext) SetBlockHash(blockHash []byte) {
	ctx.blockHash = blockHash
}

func (ctx *DeliverTxContext) Save() (hash []byte) {
	hash = make([]byte, 40) // TODO: remove magic number
	stateHash, _, err := ctx.stateTree.SaveVersion()
	if err != nil {
		panic("Cannot save state tree")
	}
	withdrawHash, _, err := ctx.withdrawTree.SaveVersion()
	if err != nil {
		panic("Cannot save withdraw tree")
	}
	copy(hash, withdrawHash[:20]) // Indended to put withdraw tree hash first, easier for Relay contract to parse
	copy(hash[20:], stateHash[:20])
	return hash
}

func (ctx *DeliverTxContext) ToCheckTxContext() *CheckTxContext {
	stateTree, err := ctx.stateTree.GetImmutable(ctx.stateTree.Version64())
	if err != nil {
		panic("Cannot get immutable state tree")
	}
	withdrawTree, err := ctx.withdrawTree.GetImmutable(ctx.withdrawTree.Version64())
	if err != nil {
		panic("Cannot get withdraw tree")
	}
	return &CheckTxContext{
		stateTree:    stateTree,
		withdrawTree: withdrawTree,
		blockHash:    ctx.blockHash,
	}
}
