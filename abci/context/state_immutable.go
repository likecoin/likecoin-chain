package context

import (
	"encoding/binary"
	"encoding/json"

	"github.com/tendermint/iavl"
	"github.com/tendermint/tendermint/libs/db"
)

// IImmutableState is an interface for accessing mutable context
type IImmutableState interface {
	ImmutableStateTree() *iavl.ImmutableTree
	ImmutableWithdrawTree() *iavl.ImmutableTree
	GetBlockHash() []byte
	GetHeight() int64
	GetMetadataAtHeight(height int64) *TreeMetadata
}

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

// GetMetadataAtHeight gets the metadata of the trees corresponding to given block height
func (state *ImmutableState) GetMetadataAtHeight(height int64) *TreeMetadata {
	key := heightMetadataKey(height)
	bs := state.appDb.Get(key)
	if bs == nil {
		return nil
	}
	metadata := TreeMetadata{}
	err := json.Unmarshal(bs, &metadata)
	if err != nil {
		log.WithError(err).Panic("Cannot unmarshal tree metadata")
	}
	return &metadata
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
