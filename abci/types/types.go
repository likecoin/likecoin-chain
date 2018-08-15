package types

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
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
