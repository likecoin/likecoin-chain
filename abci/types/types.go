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

func (rawAddr *Address) IsValidFormat() bool {
	return len(rawAddr.Content) == 20
}

func (rawAddr *Address) ToEthereum() common.Address {
	addrBytes := rawAddr.Content
	return common.BytesToAddress(addrBytes)
}

func (rawBigInt *BigInteger) ToBigInt() *big.Int {
	bigInt := new(big.Int)
	return bigInt.SetBytes(rawBigInt.Content)
}

var sigPrefix = []byte("\x19Ethereum Signed Message:\n")

// generateSigningMessageHash wraps a message in follwing format
// `\x19Ethereum Signed Message:\n" + len(message) + message`
// and return Keccak256 hash
func generateSigningMessageHash(msg []byte) []byte {
	msgLen := len(msg)
	msgLenStr := fmt.Sprintf("%d", msgLen)

	hashingMsg := make([]byte, len(sigPrefix)+len(msgLenStr)+msgLen)
	copy(hashingMsg, sigPrefix)
	copy(hashingMsg[len(sigPrefix):], []byte(msgLenStr))
	copy(hashingMsg[len(sigPrefix)+len(msgLenStr):], msg)

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
