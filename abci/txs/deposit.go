package txs

import (
	"math/big"

	"github.com/likecoin/likechain/abci/context"
	"github.com/likecoin/likechain/abci/response"
	"github.com/likecoin/likechain/abci/types"
)

// DepositInput represents one of the inputs in a Deposit transaction
type DepositInput struct {
	FromAddr types.Address
	Value    types.BigInt
}

// DepositTransaction represents a Deposit transaction
type DepositTransaction struct {
	BlockNumber uint64
	Inputs      []DepositInput
}

// ValidateFormat checks if a transaction has invalid format, e.g. nil fields, negative transfer amounts
func (tx *DepositTransaction) ValidateFormat() bool {
	if len(tx.Inputs) == 0 {
		return false
	}
	zero := big.NewInt(0)
	limit := new(big.Int).Exp(big.NewInt(2), big.NewInt(256), nil)
	for _, input := range tx.Inputs {
		if input.Value.Int.Cmp(zero) <= 0 || input.Value.Int.Cmp(limit) >= 0 {
			return false
		}
	}
	return true
}

// CheckTx checks the transaction to see if it should be executed
func (tx *DepositTransaction) CheckTx(state context.IImmutableState) response.R {
	return response.R{} // TODO
}

// DeliverTx checks the transaction to see if it should be executed
func (tx *DepositTransaction) DeliverTx(state context.IMutableState, txHash []byte) response.R {
	return response.R{} // TODO
}
