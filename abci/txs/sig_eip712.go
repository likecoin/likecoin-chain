package txs

import (
	"bytes"

	"github.com/likecoin/likechain/abci/types"
	"github.com/likecoin/likechain/abci/utils"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

const (
	domainName    = "LikeChain Signature"
	domainVersion = "1"
)

var eip712Prefix []byte

// EIP712Type represents the type and value of a field in EIP-712
type EIP712Type interface {
	TypeName() string
	EncodedValue() []byte
}

// EIP712String represents a string type field in EIP-712 standard
type EIP712String string

// TypeName returns the name of the type, which is used in type hash of EIP-712
func (s EIP712String) TypeName() string {
	return "string"
}

// EncodedValue returns the value encoded in binary following EIP-712 standard
func (s EIP712String) EncodedValue() []byte {
	return crypto.Keccak256([]byte(s))
}

// EIP712Identifier represents an LikeChain Identifier type field in EIP-712 standard
type EIP712Identifier struct {
	identifier types.Identifier
}

// TypeName returns the name of the type, which is used in type hash of EIP-712
func (iden EIP712Identifier) TypeName() string {
	return "string"
}

// EncodedValue returns the value encoded in binary following EIP-712 standard
func (iden EIP712Identifier) EncodedValue() []byte {
	return crypto.Keccak256([]byte(iden.identifier.EIP712String()))
}

// EIP712Uint256 represents a uint256 type field in EIP-712 standard
type EIP712Uint256 types.BigInt

// TypeName returns the name of the type, which is used in type hash of EIP-712
func (n EIP712Uint256) TypeName() string {
	return "uint256"
}

// EncodedValue returns the value encoded in binary following EIP-712 standard
func (n EIP712Uint256) EncodedValue() []byte {
	return types.BigInt(n).ToUint256Bytes()
}

// EIP712Uint64 represents a uint64 type field in EIP-712 standard
type EIP712Uint64 uint64

// TypeName returns the name of the type, which is used in type hash of EIP-712
func (n EIP712Uint64) TypeName() string {
	return "uint64"
}

// EncodedValue returns the value encoded in binary following EIP-712 standard
func (n EIP712Uint64) EncodedValue() []byte {
	return types.NewBigInt(int64(n)).ToUint256Bytes()
}

// EIP712Address represents an address type field in EIP-712 standard
type EIP712Address types.Address

// TypeName returns the name of the type, which is used in type hash of EIP-712
func (addr EIP712Address) TypeName() string {
	return "address"
}

// EncodedValue returns the value encoded in binary following EIP-712 standard
func (addr EIP712Address) EncodedValue() []byte {
	bs := make([]byte, 32)
	copy(bs[12:], addr[:])
	return bs
}

// EIP712Bytes32 represents a bytes32 type field in EIP-712 standard
type EIP712Bytes32 []byte

// TypeName returns the name of the type, which is used in type hash of EIP-712
func (bs EIP712Bytes32) TypeName() string {
	return "bytes32"
}

// EncodedValue returns the value encoded in binary following EIP-712 standard
func (bs EIP712Bytes32) EncodedValue() []byte {
	l := len(bs)
	if l == 32 {
		return bs
	}
	// l should not be greater than 32
	result := make([]byte, 32)
	copy(result, bs)
	return result
}

// EIP712Field represents a field in a struct in EIP-712 standard. It includes the field name, type and value
type EIP712Field struct {
	Name  string
	Value EIP712Type
}

// EIP712SignData represents the whole sign data in EIP-712 standard, including domain, struct name and fields
type EIP712SignData struct {
	Name   string
	Fields []EIP712Field
}

// Hash takes a EIP-712 sign data, returns the hash for signing the message
func (signData EIP712SignData) Hash() ([]byte, error) {
	buf := new(bytes.Buffer)
	buf.Write([]byte(signData.Name))
	buf.WriteByte('(')
	for i := 0; i < len(signData.Fields)-1; i++ {
		field := signData.Fields[i]
		buf.Write([]byte(field.Value.TypeName()))
		buf.WriteByte(' ')
		buf.Write([]byte(field.Name))
		buf.WriteByte(',')
	}
	field := signData.Fields[len(signData.Fields)-1]
	buf.Write([]byte(field.Value.TypeName()))
	buf.WriteByte(' ')
	buf.Write([]byte(field.Name))
	buf.WriteByte(')')
	bs := buf.Bytes()
	typeHash := crypto.Keccak256(bs)

	buf = new(bytes.Buffer)
	buf.Write(typeHash)
	for _, field := range signData.Fields {
		buf.Write((field.Value.EncodedValue()))
	}
	bs = buf.Bytes()
	structHash := crypto.Keccak256(bs)

	buf = new(bytes.Buffer)
	buf.Write(eip712Prefix)
	buf.Write(structHash)

	bs = buf.Bytes()
	return crypto.Keccak256(bs), nil
}

// EIP712Signature is the signature format using EIP-712 standard
type EIP712Signature [65]byte

func (sig *EIP712Signature) String() string {
	return common.ToHex(sig[:])
}

// RecoverAddress recover the signature to address by the schema of the message
func (sig *EIP712Signature) RecoverAddress(signData EIP712SignData) (*types.Address, error) {
	hash, err := signData.Hash()
	if err != nil {
		return nil, err
	}
	addr, err := recoverEthSignature(hash, *sig)
	return addr, nil
}

func init() {
	typeString := "EIP712Domain(string name,string version)"
	typeHash := crypto.Keccak256([]byte(typeString))

	buf := new(bytes.Buffer)
	buf.Write(typeHash)
	buf.Write(EIP712String(domainName).EncodedValue())
	buf.Write(EIP712String(domainVersion).EncodedValue())
	bs := buf.Bytes()
	structHash := crypto.Keccak256(bs)

	buf = new(bytes.Buffer)
	buf.Write([]byte{0x19, 0x01})
	buf.Write(structHash)
	eip712Prefix = buf.Bytes()
}

// SigEIP712 transforms a hex string into [65]byte which could be converted into signatures, panic if the string is not
// a valid signature
func SigEIP712(sigHex string) (sig EIP712Signature) {
	sigBytes, err := utils.Hex2Bytes(sigHex)
	if err != nil {
		panic(err)
	}
	if len(sigBytes) != 65 {
		panic("Invalid signature length")
	}
	copy(sig[:], sigBytes)
	return sig
}
