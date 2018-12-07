package withdraw

import (
	"bytes"
	"encoding/json"

	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto/tmhash"
	cmn "github.com/tendermint/tendermint/libs/common"
)

type proofLeafNode struct {
	Key       cmn.HexBytes `json:"key"`
	ValueHash cmn.HexBytes `json:"value"`
	Version   int64        `json:"version"`
}

type proofInnerNode struct {
	Height  int8   `json:"height"`
	Size    int64  `json:"size"`
	Version int64  `json:"version"`
	Left    []byte `json:"left"`
	Right   []byte `json:"right"`
}

type pathToLeaf []proofInnerNode

// RangeProof is the same as the RangeProof in github.com/tendermint/iavl
// The whole struct is copied here to generate contract proof from unexported internal fields.
type RangeProof struct {
	// You don't need the right path because
	// it can be derived from what we have.
	LeftPath   pathToLeaf      `json:"left_path"`
	InnerNodes []pathToLeaf    `json:"inner_nodes"`
	Leaves     []proofLeafNode `json:"leaves"`

	// memoize
	rootVerified bool
	rootHash     []byte // valid iff rootVerified is true
	treeEnd      bool   // valid iff rootVerified is true
}

func (pln proofLeafNode) Hash() []byte {
	hasher := tmhash.New()
	buf := new(bytes.Buffer)

	amino.EncodeInt8(buf, 0)
	amino.EncodeVarint(buf, 1)
	amino.EncodeVarint(buf, pln.Version)
	amino.EncodeByteSlice(buf, pln.Key)
	amino.EncodeByteSlice(buf, pln.ValueHash)
	hasher.Write(buf.Bytes())

	return hasher.Sum(nil)
}

func (pin proofInnerNode) Hash(childHash []byte) []byte {
	hasher := tmhash.New()
	buf := new(bytes.Buffer)

	amino.EncodeInt8(buf, pin.Height)
	amino.EncodeVarint(buf, pin.Size)
	amino.EncodeVarint(buf, pin.Version)

	if len(pin.Left) == 0 {
		amino.EncodeByteSlice(buf, childHash)
		amino.EncodeByteSlice(buf, pin.Right)
	} else {
		amino.EncodeByteSlice(buf, pin.Left)
		amino.EncodeByteSlice(buf, childHash)
	}

	hasher.Write(buf.Bytes())
	return hasher.Sum(nil)
}

// ComputeRootHash compute the root hash of the Merkle tree associated with the proof
func (proof *RangeProof) ComputeRootHash() (rootHash []byte) {
	path := proof.LeftPath
	leaf := proof.Leaves[0]

	rootHash = leaf.Hash()
	for i := len(path) - 1; i >= 0; i-- {
		rootHash = path[i].Hash(rootHash)
	}

	return rootHash
}

// ContractProof generates the proof used in Relay contract
func (proof *RangeProof) ContractProof() []byte {
	output := new(bytes.Buffer)
	buf := new(bytes.Buffer)

	output.WriteByte(0) // reserved proof version field

	leaf := proof.Leaves[0]
	// version bytes & length
	amino.EncodeVarint(buf, leaf.Version)
	output.WriteByte(byte(buf.Len()))
	output.Write(buf.Bytes())

	pathLength := len(proof.LeftPath)
	output.WriteByte(byte(pathLength))

	for i := int(pathLength) - 1; i >= 0; i-- {
		buf.Reset()
		pathNode := proof.LeftPath[i]
		amino.EncodeInt8(buf, pathNode.Height)
		amino.EncodeVarint(buf, pathNode.Size)
		amino.EncodeVarint(buf, pathNode.Version)
		prefixLengthAndOrder := uint8(buf.Len()) & 0x7f
		if len(pathNode.Left) != 0 {
			prefixLengthAndOrder |= 0x80
		}
		output.WriteByte(prefixLengthAndOrder)
		output.Write(buf.Bytes())
		if len(pathNode.Left) != 0 {
			output.Write(pathNode.Left)
		} else {
			output.Write(pathNode.Right)
		}
	}

	return output.Bytes()
}

// ParseRangeProof takes a JSON string in []byte, returns the unmarshaled RangeProof structure
func ParseRangeProof(data []byte) *RangeProof {
	proof := RangeProof{}
	err := json.Unmarshal(data, &proof)
	if err != nil {
		return nil
	}
	return &proof
}
