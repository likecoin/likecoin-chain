package withdraw

import (
	amino "github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto/merkle"
	"github.com/tendermint/tendermint/crypto/tmhash"
	cmn "github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/types"
)

func simpleHashFromTwoHashes(left, right []byte) []byte {
	var hasher = tmhash.New()
	err := amino.EncodeByteSlice(hasher, left)
	if err != nil {
		panic(err)
	}
	err = amino.EncodeByteSlice(hasher, right)
	if err != nil {
		panic(err)
	}
	return hasher.Sum(nil)
}

func cdcEncode(item interface{}) []byte {
	if item != nil && !cmn.IsTypedNil(item) && !cmn.IsEmpty(item) {
		return types.GetCodec().MustMarshalBinaryBare(item)
	}
	return nil
}

// Proof returns the merkle proof for the AppHash in the Tendermint header
func Proof(h *types.Header) (root []byte, proof [][]byte) {
	lv0 := merkle.SimpleHashFromByteSlices([][]byte{cdcEncode(h.AppHash)})
	lv1 := merkle.SimpleHashFromByteSlices([][]byte{cdcEncode(h.LastResultsHash)})
	lv2 := merkle.SimpleHashFromByteSlices([][]byte{
		cdcEncode(h.EvidenceHash),
		cdcEncode(h.ProposerAddress),
	})
	lv3 := merkle.SimpleHashFromByteSlices([][]byte{
		cdcEncode(h.DataHash),
		cdcEncode(h.ValidatorsHash),
		cdcEncode(h.NextValidatorsHash),
		cdcEncode(h.ConsensusHash),
	})
	lv4 := merkle.SimpleHashFromByteSlices([][]byte{
		cdcEncode(h.Version),
		cdcEncode(h.ChainID),
		cdcEncode(h.Height),
		cdcEncode(h.Time),
		cdcEncode(h.NumTxs),
		cdcEncode(h.TotalTxs),
		cdcEncode(h.LastBlockID),
		cdcEncode(h.LastCommitHash),
	})
	root = lv0
	root = simpleHashFromTwoHashes(root, lv1)
	root = simpleHashFromTwoHashes(root, lv2)
	root = simpleHashFromTwoHashes(lv3, root)
	root = simpleHashFromTwoHashes(lv4, root)
	return root, [][]byte{lv1, lv2, lv3, lv4}
}
