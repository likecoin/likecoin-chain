package context

import "github.com/tendermint/iavl"

type CheckTxContext struct {
	stateTree    *iavl.ImmutableTree
	withdrawTree *iavl.ImmutableTree
	blockHash    []byte
}

func (ctx *CheckTxContext) StateTree() *iavl.ImmutableTree {
	return ctx.stateTree
}

func (ctx *CheckTxContext) WithdrawTree() *iavl.ImmutableTree {
	return ctx.withdrawTree
}

func (ctx *CheckTxContext) GetBlockHash() []byte {
	return ctx.blockHash
}
