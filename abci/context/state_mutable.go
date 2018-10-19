package context

import (
	"bytes"
	"encoding/binary"
	"math/big"

	"github.com/tendermint/iavl"
	cmn "github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/libs/db"
)

// IMutableState is an interface for accessing immutable context
type IMutableState interface {
	IImmutableState
	MutableStateTree() *iavl.MutableTree
	MutableWithdrawTree() *iavl.MutableTree
	GetInitialBalance() *big.Int
	GetKeepBlocks() int64
}

// MutableState is a struct contains mutable state
type MutableState struct {
	appDb      db.DB
	stateDb    db.DB
	withdrawDb db.DB

	stateTree    *iavl.MutableTree
	withdrawTree *iavl.MutableTree

	initialBalance *big.Int
	keepBlocks     int64
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

func heightMetadataKey(height int64) []byte {
	buf := new(bytes.Buffer)
	buf.Write(appMetadataAtHeight)
	binary.Write(buf, binary.BigEndian, uint64(height))
	return buf.Bytes()
}

// TreeMetadata is the metadata of the trees, used for querying withdraw proof by height, and also removing outdated
// tree versions
type TreeMetadata struct {
	StateTreeVersion    int64
	WithdrawTreeVersion int64
}

// Bytes encodes the TreeMetadata into byte array so that it could be saved into DB
func (metadata TreeMetadata) Bytes() []byte {
	result := make([]byte, 17)
	// First byte: schema version
	result[0] = 0
	binary.BigEndian.PutUint64(result[1:], uint64(metadata.StateTreeVersion))
	binary.BigEndian.PutUint64(result[9:], uint64(metadata.WithdrawTreeVersion))
	return result
}

// decodeTreeMetadata decode a byte array into TreeMetadata
func decodeTreeMetadata(bs []byte) *TreeMetadata {
	if bs[0] != 0 || len(bs) < 17 {
		return nil
	}
	stateTreeVersion := int64(binary.BigEndian.Uint64(bs[1:]))
	withdrawTreeVersion := int64(binary.BigEndian.Uint64(bs[9:]))
	return &TreeMetadata{
		StateTreeVersion:    stateTreeVersion,
		WithdrawTreeVersion: withdrawTreeVersion,
	}
}

// SetMetadataAtHeight stores the metadata of the trees by height.
func (state *MutableState) SetMetadataAtHeight(height int64, metadata TreeMetadata) {
	key := heightMetadataKey(height)
	bs := metadata.Bytes()
	state.appDb.Set(key, bs)
}

// GetMetadataAtHeight gets the metadata of the trees corresponding to given block height
func (state *MutableState) GetMetadataAtHeight(height int64) *TreeMetadata {
	key := heightMetadataKey(height)
	bs := state.appDb.Get(key)
	if bs == nil {
		return nil
	}
	metadata := decodeTreeMetadata(bs)
	if metadata == nil {
		log.WithField("data", cmn.HexBytes(bs)).Panic("Cannot unmarshal tree metadata")
	}
	return metadata
}

// GC removes outdated versions of the trees
func (state *MutableState) GC(currentHeight int64) {
	removeBeforeHeight := currentHeight - int64(state.GetKeepBlocks())
	if removeBeforeHeight <= 0 {
		return
	}
	keyStart := heightMetadataKey(0)
	keyEnd := heightMetadataKey(removeBeforeHeight + 1)
	it := state.appDb.Iterator(keyStart, keyEnd)
	defer it.Close()
	if !it.Valid() {
		return
	}
	keysToRemove := make([][]byte, 0, 1)
	for ; it.Valid(); it.Next() {
		key := it.Key()
		keysToRemove = append(keysToRemove, key)
		bs := state.appDb.Get(key)
		metadata := decodeTreeMetadata(bs)
		if metadata == nil {
			log.WithField("data", cmn.HexBytes(bs)).Panic("Cannot unmarshal tree metadata")
		}
		state.MutableStateTree().DeleteVersion(metadata.StateTreeVersion)
		state.MutableWithdrawTree().DeleteVersion(metadata.WithdrawTreeVersion)
	}
	for _, key := range keysToRemove {
		state.appDb.Delete(key)
	}
}

// GetInitialBalance returns the initial balance for new account
func (state *MutableState) GetInitialBalance() *big.Int {
	if state.initialBalance == nil {
		initialBalance, success := new(big.Int).SetString(config.InitialBalance, 10)
		if !success {
			initialBalance = big.NewInt(0)
		}
		state.initialBalance = initialBalance
	}
	return state.initialBalance
}

// GetKeepBlocks returns the number of blocks kept in trees
func (state *MutableState) GetKeepBlocks() int64 {
	if state.keepBlocks == 0 {
		keepBlocks := config.KeepBlocks
		if keepBlocks <= 0 {
			keepBlocks = 10000
		}
		state.keepBlocks = keepBlocks
	}
	return state.keepBlocks
}

// Init initializes states
func (state *MutableState) Init() {
	log.Info("Init states")
	state.SetHeight(0)
	state.withdrawTree.Set(initKey, []byte{})
}
