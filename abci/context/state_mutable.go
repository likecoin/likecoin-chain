package context

import (
	"github.com/tendermint/iavl"
)

// MutableState is a struct contains mutable state
type MutableState struct {
	stateTree    *iavl.MutableTree
	withdrawTree *iavl.MutableTree
}

// ImmutableStateTree returns immutable state tree of the current state
func (state *MutableState) ImmutableStateTree() *iavl.ImmutableTree {
	return state.stateTree.ImmutableTree
}

// ImmutableWithdrawTree returns immutable state tree of the current state
func (state *MutableState) ImmutableWithdrawTree() *iavl.ImmutableTree {
	return state.withdrawTree.ImmutableTree
}

// MutableStateTree returns mutable state tree of the current state
func (state *MutableState) MutableStateTree() *iavl.MutableTree {
	return state.stateTree
}

// MutableWithdrawTree returns mutable withdraw tree of the current state
func (state *MutableState) MutableWithdrawTree() *iavl.MutableTree {
	return state.withdrawTree
}

// GetBlockHash returns the block hash of the current state
func (state *MutableState) GetBlockHash() []byte {
	return nil // TODO
}

// SetBlockHash saves a block hash to the current state
func (state *MutableState) SetBlockHash(blockHash []byte) {
	// TODO
}

// Save saves a new state tree version and a new state withdraw tree version,
// based on the current state of those trees.
// Returns a merged hash from those trees
func (state *MutableState) Save() (hash []byte) {
	stateHash, _, err := state.stateTree.SaveVersion()
	if err != nil {
		log.WithError(err).Panic("Cannot save state tree")
	}
	withdrawHash, _, err := state.withdrawTree.SaveVersion()
	if err != nil {
		log.WithError(err).Panic("Cannot save withdraw tree")
	}
	hash = make([]byte, 40) // TODO: remove magic number
	// TODO: After InitChain implementation, the hash will be no longer empty
	// the following if checking can be removed by then
	if len(stateHash) >= 20 && len(withdrawHash) >= 20 {
		// Indended to put withdraw tree hash first,
		// easier for Relay contract to parse
		copy(hash, withdrawHash[:20])
		copy(hash[20:], stateHash[:20])
	}
	return hash
}
