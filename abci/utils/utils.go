package utils

import (
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/likecoin/likechain/abci/types"
)

func RecoverSignature(hash []byte, sig types.Signature) (common.Address, error) {
	pubKeyBytes, err := crypto.Ecrecover(hash, sig.Content)
	if err != nil {
		return common.Address{}, errors.New("Invalid signature")
	}
	pubKey := crypto.ToECDSAPub(pubKeyBytes)

	return crypto.PubkeyToAddress(*pubKey), nil
}
