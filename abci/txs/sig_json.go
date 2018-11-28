package txs

import (
	"encoding/json"
	"fmt"

	"github.com/likecoin/likechain/abci/types"
	"github.com/likecoin/likechain/abci/utils"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// JSONSignature is the signature format using deterministic JSON as message representation
type JSONSignature [65]byte

func (sig *JSONSignature) String() string {
	return common.ToHex(sig[:])
}

// JSONMap represents a JSON object for signing message
type JSONMap map[string]interface{}

// Hash takes a JSONMap representing a JSON object, returns the hash for signing the message
func (jsonMap JSONMap) Hash() ([]byte, error) {
	msg, err := json.Marshal(jsonMap)
	if err != nil {
		return nil, err
	}
	sigPrefix := "\x19Ethereum Signed Message:\n"
	hashingMsg := []byte(fmt.Sprintf("%s%d%s", sigPrefix, len(msg), msg))
	return crypto.Keccak256(hashingMsg), nil
}

// RecoverAddress recover the signature to address by the deterministic JSON representation of the message
func (sig *JSONSignature) RecoverAddress(jsonMap JSONMap) (*types.Address, error) {
	hash, err := jsonMap.Hash()
	if err != nil {
		return nil, err
	}
	addr, err := recoverEthSignature(hash, *sig)
	return addr, err
}

// Sig transforms a hex string into [65]byte which could be converted into signatures, panic if the string is not a
// valid signature
func Sig(sigHex string) (sig JSONSignature) {
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
