package utils

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/likecoin/likechain/abci/types"
)

var SIG_PREFIX = []byte("\x19Ethereum Signed Message:\n")

// Format: "\x19Ethereum Signed Message:\n" + len(message) + message
func RecoverJsonToEthereumAddress(jsonBytes []byte, sig types.Signature) (common.Address, error) {
	l := len(jsonBytes)
	lenString := fmt.Sprintf("%d", l)
	buf := make([]byte, len(SIG_PREFIX)+len(lenString)+l)
	copy(buf, SIG_PREFIX)
	copy(buf[len(SIG_PREFIX):], []byte(lenString))
	copy(buf[len(SIG_PREFIX)+len(lenString):], jsonBytes)
	hash := crypto.Keccak256(buf)
	pubKeyBytes, err := crypto.Ecrecover(hash, sig.Content)
	if err != nil {
		return common.Address{}, errors.New("Invalid signature")
	}
	pubKey := crypto.ToECDSAPub(pubKeyBytes)
	return crypto.PubkeyToAddress(*pubKey), nil
}
