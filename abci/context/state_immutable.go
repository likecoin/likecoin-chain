package context

import (
	"encoding/binary"

	"github.com/tendermint/iavl"
	"github.com/tendermint/tendermint/libs/db"
)

// ImmutableState is a struct contains the most recently saved state
type ImmutableState struct {
	appDb        db.DB
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
	value := state.appDb.Get(appBlockHashKey)
	return value
}

// GetHeight returns the most recently block height
func (state *ImmutableState) GetHeight() int64 {
	value := state.appDb.Get(appHeightKey)
	if value == nil {
		return 0
	}
	return int64(binary.BigEndian.Uint64(value))
}

// SetWithdrawVersionAtHeight is used to store the withdraw tree version mapping corresponding to the block height
func (state *ImmutableState) SetWithdrawVersionAtHeight(height int64, version int64) {
	key := heightWithdrawVersionKey(height)
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(height))
	state.appDb.Set(key, buf)
}

// GetWithdrawVersionAtHeight gets the withdraw tree version corresponding to the block height
func (state *ImmutableState) GetWithdrawVersionAtHeight(height int64) int64 {
	key := heightWithdrawVersionKey(height)
	buf := state.appDb.Get(key)
	if buf == nil {
		return -1
	}
	return int64(binary.BigEndian.Uint64(buf))
}

const appHashLength = 40

func generateAppHash(stateHash, withdrawHash []byte) (hash []byte) {
	hash = make([]byte, appHashLength)
	// Indended to put withdraw tree hash first,
	// easier for Relay contract to parse
	if binary.Size(withdrawHash) > 0 {
		copy(hash, withdrawHash[:appHashLength/2])
	}
	if binary.Size(stateHash) > 0 {
		copy(hash[appHashLength/2:], stateHash[:appHashLength/2])
	}
	return hash
}

// GetAppHash returns the app hash of the most recently saved state
func (state *ImmutableState) GetAppHash() []byte {
	if state.GetHeight() == 0 {
		return nil
	}
	return generateAppHash(state.stateTree.Hash(), state.withdrawTree.Hash())
}
