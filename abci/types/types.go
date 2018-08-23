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

func (rawSig *Signature) ToHex() string {
	return common.ToHex(rawSig.Content)
}

// ToBigInt converts BigInteger struct to big Int
func (rawBigInt *BigInteger) ToBigInt() *big.Int {
	bigInt := new(big.Int)
	return bigInt.SetBytes(rawBigInt.Content)
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
