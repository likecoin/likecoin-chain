package context

import "github.com/tendermint/iavl"

// ImmutableState is a struct contains the most recently saved state
type ImmutableState struct {
	stateTree    *iavl.ImmutableTree
	withdrawTree *iavl.ImmutableTree
}

// ImmutableStateTree returns the most recently saved state tree
func (state *ImmutableState) ImmutableStateTree() *iavl.ImmutableTree {
	return state.stateTree
}

// ImmutableWithdrawTree returns the most recently saved withdraw tree
func (state *ImmutableState) ImmutableWithdrawTree() *iavl.ImmutableTree {
	return state.withdrawTree
}

// GetBlockHash returns the most recently saved block hash
func (state *ImmutableState) GetBlockHash() []byte {
	return nil // TODO
}
