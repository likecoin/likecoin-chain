package account

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/likecoin/likechain/abci/context"
	"github.com/likecoin/likechain/abci/types"
)

// NewAccount creates a new account
func NewAccount(ctx context.Context, address common.Address) (types.LikeChainID, error) {
	id := generateLikeChainID(ctx)
	// TODO: save to db
	// TODO: initialize account info
	// TODO: check if address already has balance
	return id, nil // TODO
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

func SaveBalance(ctx context.Context, id types.LikeChainID, balance *big.Int) error {
	return nil // TODO
}

func FetchBalance(ctx context.Context, id types.LikeChainID) *big.Int {
	return nil // TODO
}

func FetchEthereumAddressBalance(ctx context.Context, addr common.Address) *big.Int {
	return nil // TODO
}

func IncrementNextNonce(ctx context.Context, id types.LikeChainID) {
	// TODO
}

func FetchNextNonce(ctx context.Context, id types.LikeChainID) uint64 {
	return 0 // TODO
}
