package types

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// IsValidFormat checks the length of the address
func (rawAddr *Address) IsValidFormat() bool {
	return len(rawAddr.Content) == 20
}

// ToEthereum converts Address struct to Ethereum address
func (rawAddr *Address) ToEthereum() common.Address {
	addrBytes := rawAddr.Content
	return common.BytesToAddress(addrBytes)
}

func (rawAddr *Address) ToHex() string {
	return rawAddr.ToEthereum().Hex()
}

// NewSignatureFromHex creates Signature from hex string
func NewSignatureFromHex(hex string) *Signature {
	return &Signature{
		Content: common.FromHex(hex),
		Version: 1,
	}
}

// ToHex converts Signature to hex string
func (sig *Signature) ToHex() string {
	return common.ToHex(sig.Content)
}

// IsValidFormat checks the signature version and its length
func (sig *Signature) IsValidFormat() bool {
	switch sig.Version {
	case 1:
		content := sig.Content
		if len(content) != 65 {
			return false
		}
		return true
	default:
		return false
	}
}

// NewBigInteger creates BigInteger from int64
func NewBigInteger(s string) *BigInteger {
	i := new(big.Int)
	i.SetString(s, 10)
	return &BigInteger{Content: i.Bytes()}
}

// ToBigInt converts BigInteger struct to big Int
func (i *BigInteger) ToBigInt() *big.Int {
	bigInt := new(big.Int)
	return bigInt.SetBytes(i.Content)
}

// ToString converts BigInteger struct to string
func (i *BigInteger) ToString() string {
	bigInt := new(big.Int)
	bigInt.SetBytes(i.Content)
	return bigInt.String()
}

// IsValidFormat checks Identifier format
func (id *Identifier) IsValidFormat() bool {
	return (id.GetLikeChainID() != nil && len(id.GetLikeChainID().Content) > 0) ||
		(id.GetAddr() != nil && len(id.GetAddr().Content) > 0)
}

// NewLikeChainID creates a LikeChainID from bytes
func NewLikeChainID(content []byte) LikeChainID {
	return LikeChainID{Content: content}
}

// ToIdentifier converts LikeChainID to Identifier
func (id *LikeChainID) ToIdentifier() *Identifier {
	return &Identifier{
		Id: &Identifier_LikeChainID{
			LikeChainID: id,
		},
	}
}

// NewTransferTarget creates a new TransferTarget
func NewTransferTarget(id *Identifier, value string, remark string) *TransferTransaction_TransferTarget {
	return &TransferTransaction_TransferTarget{
		To:     id,
		Value:  NewBigInteger(value),
		Remark: []byte(remark),
	}
}

// IsValidFormat checks the format of TransferTarget ensuring the `value` and
// `to` are correct
func (t *TransferTransaction_TransferTarget) IsValidFormat() bool {
	value := t.Value
	to := t.To
	if value == nil ||
		to == nil {
		return false
	}
	return value.ToBigInt().Cmp(big.NewInt(0)) >= 0 &&
		to.IsValidFormat()
}

var sigPrefix = []byte("\x19Ethereum Signed Message:\n")

// generateSigningMessageHash wraps a message in follwing format
// `\x19Ethereum Signed Message:\n" + len(message) + message`
// and return Keccak256 hash
func generateSigningMessageHash(msg []byte) []byte {
	hashingMsg := []byte(fmt.Sprintf("%s%d%s", sigPrefix, len(msg), msg))
	return crypto.Keccak256(hashingMsg)
}

// GenerateSigningMessageHash generates a signature from a RegisterTx
func (tx *RegisterTransaction) GenerateSigningMessageHash() ([]byte, error) {
	m := map[string]interface{}{
		"addr": strings.ToLower(tx.Addr.ToEthereum().Hex()),
	}

	msg, err := json.Marshal(m)
	if err != nil {
		return nil, errors.New("Unable to marshal JSON string for RegisterTx")
	}

	return generateSigningMessageHash(msg), nil
}

func (tx *RegisterTransaction) ToString() string {
	return fmt.Sprintf(
		"<Addr: %s, Sig: %s>",
		tx.Addr.ToHex(),
		tx.Sig.ToHex(),
	)
}
