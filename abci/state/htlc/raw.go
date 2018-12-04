package htlc

import (
	"github.com/likecoin/likechain/abci/context"
	logger "github.com/likecoin/likechain/abci/log"
	"github.com/likecoin/likechain/abci/types"
	"github.com/likecoin/likechain/abci/utils"

	cmn "github.com/tendermint/tendermint/libs/common"
)

// HashedTransfer represents a Hashed TimeLock Transfer
type HashedTransfer struct {
	From       types.Identifier
	To         types.Identifier
	Value      types.BigInt
	HashCommit [32]byte
	Expiry     int64
}

// Validate checks whether a HashedTransfer is valid
func (ht *HashedTransfer) Validate() bool {
	if ht.From == nil || ht.To == nil || ht.Value.Int == nil {
		return false
	}
	if !ht.Value.IsWithinRange() || ht.Value.Int.Int64() == 0 {
		return false
	}
	if ht.Expiry <= 0 {
		return false
	}
	return true
}

var (
	htKey = []byte("hashedTx")

	log = logger.L
)

func hashedTransferKey(txHash []byte) []byte {
	return utils.JoinKeys([][]byte{
		htKey,
		txHash,
	})
}

// CreateHashedTransfer stores a HashedTransfer into state tree, associated with a transaction hash
func CreateHashedTransfer(state context.IMutableState, ht *HashedTransfer, txHash []byte) {
	bs, err := types.AminoCodec().MarshalBinaryLengthPrefixed(ht)
	if err != nil {
		log.
			WithField("hashed_transfer", ht).
			WithError(err).
			Panic("Cannot marshal HashedTransfer")
	}
	key := hashedTransferKey(txHash)
	state.MutableStateTree().Set(key, bs)
}

// GetHashedTransfer loads a HashedTransfer associated with the transaction hash from state tree
func GetHashedTransfer(state context.IImmutableState, txHash []byte) *HashedTransfer {
	key := hashedTransferKey(txHash)
	_, bs := state.ImmutableStateTree().Get(key)
	if bs == nil {
		return nil
	}
	ht := HashedTransfer{}
	err := types.AminoCodec().UnmarshalBinaryLengthPrefixed(bs, &ht)
	if err != nil {
		log.
			WithField("data", cmn.HexBytes(bs)).
			WithError(err).
			Panic("Cannot unmarshal HashedTransfer")
	}
	return &ht
}

// RemoveHashedTransfer removes a HashedTransfer associated with the transaction hash from state tree
func RemoveHashedTransfer(state context.IMutableState, txHash []byte) {
	key := hashedTransferKey(txHash)
	state.MutableStateTree().Remove(key)
}
