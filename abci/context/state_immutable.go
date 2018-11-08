package context

import (
	"bytes"
	"encoding/binary"
	"encoding/json"

	"github.com/likecoin/likechain/abci/utils"

	"github.com/tendermint/iavl"
	"github.com/tendermint/tendermint/crypto/tmhash"
	"github.com/tendermint/tendermint/libs/db"
)

// IImmutableState is an interface for accessing mutable context
type IImmutableState interface {
	ImmutableStateTree() *iavl.ImmutableTree
	ImmutableWithdrawTree() *iavl.ImmutableTree
	GetBlockHash() []byte
	GetBlockTime() int64
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

// GetBlockTime returns the block time of the executing block
func (state *ImmutableState) GetBlockTime() int64 {
	bs := state.appDb.Get(appBlockTimeKey)
	if len(bs) == 0 {
		return 0
	}
	return int64(utils.DecodeUint64(bs))
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

var zeros = make([]byte, tmhash.Size)

func generateAppHash(stateHash, withdrawHash []byte) (hash []byte) {
	hashBuf := new(bytes.Buffer)
	// Indended to put withdraw tree hashBuf first,
	// easier for Relay contract to parse
	if len(withdrawHash) > 0 {
		hashBuf.Write(withdrawHash)
	} else {
		hashBuf.Write(zeros)
	}
	if len(stateHash) > 0 {
		hashBuf.Write(stateHash)
	} else {
		hashBuf.Write(zeros)
	}
	return hashBuf.Bytes()
}

// GetAppHash returns the app hash of the most recently saved state
func (state *ImmutableState) GetAppHash() []byte {
	if state.GetHeight() == 0 {
		return nil
	}
	return generateAppHash(state.stateTree.Hash(), state.withdrawTree.Hash())
}
