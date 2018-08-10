package account

import (
	"github.com/likecoin/likechain/abci/context"
	"github.com/likecoin/likechain/abci/types"
)

// NewAccount creates a new account
func NewAccount(address []byte) types.Identifier {
	return types.Identifier{} // TODO
}

func generateLikeChainID(ctx context.Context) []byte {
	return nil // TODO
}

func SaveBalance(ctx context.Context, balance types.BigInteger) {
	// TODO
}

func FetchBalance(ctx context.Context, id *types.Identifier) types.BigInteger {
	return types.BigInteger{} // TODO
}

func SaveNextNonce(ctx context.Context, id *types.Identifier) {
	// TODO
}

func FetchNextNonce(ctx context.Context, id *types.Identifier) uint64 {
	return 0 // TODO
}
