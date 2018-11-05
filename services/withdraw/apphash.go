package withdraw

import (
	"io"

	amino "github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto/merkle"
	"github.com/tendermint/tendermint/crypto/tmhash"
	cmn "github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/types"

	"github.com/likecoin/likechain/services/tendermint"
)

type hasher struct {
	item interface{}
}

func (h hasher) Hash() []byte {
	hasher := tmhash.New()
	if h.item != nil && !cmn.IsTypedNil(h.item) && !cmn.IsEmpty(h.item) {
		bz, err := tendermint.AminoCodec().MarshalBinaryBare(h.item)
		if err != nil {
			panic(err)
		}
		_, err = hasher.Write(bz)
		if err != nil {
			panic(err)
		}
	}
	return hasher.Sum(nil)

}

func aminoHash(item interface{}) []byte {
	h := hasher{item}
	return h.Hash()
}

func aminoHasher(item interface{}) Hasher {
	return hasher{item}
}

// Hasher represents a hashable piece of data which can be hashed in the Tree.
type Hasher interface {
	Hash() []byte
}

// KVPair A local extension to KVPair that can be hashed.
// Key and value are length prefixed and concatenated,
// then hashed.
type KVPair cmn.KVPair

// Uvarint length prefixed byteslice
func encodeByteSlice(w io.Writer, bz []byte) (err error) {
	return amino.EncodeByteSlice(w, bz)
}

// Hash the KVPair
func (kv KVPair) Hash() []byte {
	hasher := tmhash.New()
	err := encodeByteSlice(hasher, kv.Key)
	if err != nil {
		panic(err)
	}
	err = encodeByteSlice(hasher, kv.Value)
	if err != nil {
		panic(err)
	}
	result := hasher.Sum(nil)
	return result
}

func simpleHashFromHashes(hashes [][]byte) []byte {
	// Recursive impl.
	switch len(hashes) {
	case 0:
		return nil
	case 1:
		return hashes[0]
	default:
		left := simpleHashFromHashes(hashes[:(len(hashes)+1)/2])
		right := simpleHashFromHashes(hashes[(len(hashes)+1)/2:])
		return merkle.SimpleHashFromTwoHashes(left, right)
	}
}

// Proof returns the merkle proof for the AppHash in the Tendermint header
func Proof(h *types.Header) [][]byte {
	pairs := [][]byte{
		KVPair{Key: []byte("App"), Value: aminoHasher(h.AppHash).Hash()}.Hash(),
		KVPair{Key: []byte("ChainID"), Value: aminoHasher(h.ChainID).Hash()}.Hash(),
		KVPair{Key: []byte("Consensus"), Value: aminoHasher(h.ConsensusHash).Hash()}.Hash(),
		KVPair{Key: []byte("Data"), Value: aminoHasher(h.DataHash).Hash()}.Hash(),
		KVPair{Key: []byte("Evidence"), Value: aminoHasher(h.EvidenceHash).Hash()}.Hash(),
		KVPair{Key: []byte("Height"), Value: aminoHasher(h.Height).Hash()}.Hash(),
		KVPair{Key: []byte("LastBlockID"), Value: aminoHasher(h.LastBlockID).Hash()}.Hash(),
		KVPair{Key: []byte("LastCommit"), Value: aminoHasher(h.LastCommitHash).Hash()}.Hash(),
		KVPair{Key: []byte("NumTxs"), Value: aminoHasher(h.NumTxs).Hash()}.Hash(),
		KVPair{Key: []byte("Results"), Value: aminoHasher(h.LastResultsHash).Hash()}.Hash(),
		KVPair{Key: []byte("Time"), Value: aminoHasher(h.Time).Hash()}.Hash(),
		KVPair{Key: []byte("TotalTxs"), Value: aminoHasher(h.TotalTxs).Hash()}.Hash(),
		KVPair{Key: []byte("Validators"), Value: aminoHasher(h.ValidatorsHash).Hash()}.Hash(),
	}
	// 13
	// 7 6
	// 4 3 3 3
	// 2 2 2 1 2 1 2 1
	proof := make([][]byte, 4)
	rootHash := pairs[0]
	proof[0] = pairs[1]
	rootHash = merkle.SimpleHashFromTwoHashes(rootHash, proof[0])
	proof[1] = simpleHashFromHashes(pairs[2:4])
	rootHash = merkle.SimpleHashFromTwoHashes(rootHash, proof[1])
	proof[2] = simpleHashFromHashes(pairs[4:7])
	rootHash = merkle.SimpleHashFromTwoHashes(rootHash, proof[2])
	proof[3] = simpleHashFromHashes(pairs[7:13])
	return proof
}
