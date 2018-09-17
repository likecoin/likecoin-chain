package context

import (
	"encoding/binary"

	"github.com/tendermint/iavl"
	"github.com/tendermint/tendermint/libs/db"
)

// MutableState is a struct contains mutable state
type MutableState struct {
	appDb      db.DB
	stateDb    db.DB
	withdrawDb db.DB

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
	value := state.appDb.Get(appBlockHashKey)
	return value
}

// SetBlockHash saves the block hash to the current state
func (state *MutableState) SetBlockHash(blockHash []byte) {
	state.appDb.Set(appBlockHashKey, blockHash)
}

// GetHeight returns the block height of the current state
func (state *MutableState) GetHeight() int64 {
	value := state.appDb.Get(appHeightKey)
	if value == nil {
		return 0
	}
	return int64(binary.BigEndian.Uint64(value))
}

// SetHeight saves the block height to the current state
func (state *MutableState) SetHeight(height int64) {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(height))
	state.appDb.Set(appHeightKey, buf)
}

// GetAppHash returns the app hash of the current state
func (state *MutableState) GetAppHash() []byte {
	if state.GetHeight() == 0 {
		return nil
	}
	return generateAppHash(state.stateTree.Hash(), state.withdrawTree.Hash())
}

// Save saves a new state tree version and a new state withdraw tree version,
// based on the current state of those trees.
// Returns a merged hash from those trees
func (state *MutableState) Save() []byte {
	stateHash, _, err := state.stateTree.SaveVersion()
	if err != nil {
		log.WithError(err).Panic("Cannot save state tree")
	}
	withdrawHash, _, err := state.withdrawTree.SaveVersion()
	if err != nil {
		log.WithError(err).Panic("Cannot save withdraw tree")
	}
	return generateAppHash(stateHash, withdrawHash)
}

// Init initializes states
func (state *MutableState) Init() {
	log.Info("Init states")
	state.SetHeight(0)
	state.withdrawTree.Set(initKey, []byte{})
}
