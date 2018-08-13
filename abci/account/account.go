package account

import (
	"github.com/likecoin/likechain/abci/context"
	"github.com/likecoin/likechain/abci/types"
)

// NewAccount creates a new account
func NewAccount(address []byte) (types.Identifier, bool) {
	err := true
	return types.Identifier{}, err // TODO
}

func generateLikeChainID(ctx context.Context) []byte {
	return nil // TODO
}

func SaveBalance(ctx context.Context, id *types.Identifier, balance types.BigInteger) bool {
	err := true
	return err // TODO
}

func FetchBalance(ctx context.Context, id *types.Identifier) types.BigInteger {
	return types.BigInteger{} // TODO
}

func SaveNextNonce(ctx context.Context, id *types.Identifier, nonce uint64) bool {
	err := true
	return err // TODO
}

func FetchNextNonce(ctx context.Context, id *types.Identifier) uint64 {
	return 0 // TODO
}
