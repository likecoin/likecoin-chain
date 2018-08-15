package account

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/likecoin/likechain/abci/context"
	"github.com/likecoin/likechain/abci/types"
)

// NewAccount creates a new account
func NewAccount(address common.Address) (types.Identifier, error) {
	return types.Identifier{}, nil // TODO
}

func generateLikeChainID(ctx context.Context) types.LikeChainID {
	return types.LikeChainID{} // TODO
}

func AddressToLikeChainID(ctx context.Context, addr types.Address) (types.LikeChainID, bool) {
	return types.LikeChainID{}, false // TODO
}

func GetLikeChainID(ctx context.Context, identifier types.Identifier) (types.LikeChainID, bool) {
	id := identifier.GetLikeChainID()
	if id != nil {
		return *id, false // TODO: check the existence of this LikeChainID
	}
	addr := identifier.GetAddr()
	// TODO: assert addr != nil
	return AddressToLikeChainID(ctx, *addr)
}

func SaveBalance(ctx context.Context, id types.LikeChainID, balance types.BigInteger) error {
	return nil // TODO
}

func FetchBalance(ctx context.Context, id types.LikeChainID) *big.Int {
	return nil // TODO
}

func SaveNextNonce(ctx context.Context, id types.LikeChainID, nonce uint64) error {
	return nil
}

func FetchNextNonce(ctx context.Context, id types.LikeChainID) uint64 {
	return 0 // TODO
}
