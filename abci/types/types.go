package types

import (
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

func (rawAddr Address) EthereumAddress() (*common.Address, error) {
	addrBytes := rawAddr.Content
	if len(addrBytes) != 20 {
		return nil, errors.New("Invalid length for Ethereum address")
	}
	ethAddr := common.BytesToAddress(addrBytes)
	return &ethAddr, nil
}

func (rawBigInt BigInteger) BigInt() *big.Int {
	bigInt := new(big.Int)
	return bigInt.SetBytes(rawBigInt.Content)
}
