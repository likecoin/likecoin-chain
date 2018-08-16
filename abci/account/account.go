package account

import (
	"encoding/binary"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/likecoin/likechain/abci/context"
	"github.com/likecoin/likechain/abci/types"
)

// NewAccount creates a new account
func NewAccount(ctx context.MutableContext, address common.Address) (types.LikeChainID, error) {
	id := generateLikeChainID(ctx)
	// TODO: save to db
	// TODO: initialize account info
	// TODO: check if address already has balance
	return id, nil // TODO
}

var likeChainIDSeedKey = []byte("$account.likeChainIDSeed")

func generateLikeChainID(ctx context.MutableContext) types.LikeChainID {
	var seedInt uint64
	_, seed := ctx.StateTree().Get(likeChainIDSeedKey)
	if seed == nil {
		seedInt = 1
		seed = make([]byte, 8)
		binary.BigEndian.PutUint64(seed, seedInt)
	} else {
		seedInt = uint64(binary.BigEndian.Uint64(seed))
	}

	blockHash := ctx.GetBlockHash()

	// Concat the seed and the block's hash
	content := make([]byte, len(seed)+len(blockHash))
	copy(content, seed)
	copy(content[len(seed):], blockHash)
	// Take first 20 bytes of Keccak256 hash to be LikeChainID
	content = crypto.Keccak256(content)[:20]

	// Increment and save seed
	seedInt++
	binary.BigEndian.PutUint64(seed, seedInt)
	ctx.MutableStateTree().Set(likeChainIDSeedKey, seed)

	return types.LikeChainID{Content: content}
}

func AddressToLikeChainID(ctx context.ImmutableContext, addr types.Address) (types.LikeChainID, bool) {
	return types.LikeChainID{}, false // TODO
}

func GetLikeChainID(ctx context.ImmutableContext, identifier types.Identifier) (types.LikeChainID, bool) {
	id := identifier.GetLikeChainID()
	if id != nil {
		return *id, false // TODO: check the existence of this LikeChainID
	}
	addr := identifier.GetAddr()
	// TODO: assert addr != nil
	return AddressToLikeChainID(ctx, *addr)
}

func SaveBalance(ctx context.MutableContext, id types.LikeChainID, balance *big.Int) error {
	return nil // TODO
}

func FetchBalance(ctx context.ImmutableContext, id types.LikeChainID) *big.Int {
	return nil // TODO
}

func FetchEthereumAddressBalance(ctx context.MutableContext, addr common.Address) *big.Int {
	return nil // TODO
}

func IncrementNextNonce(ctx context.MutableContext, id types.LikeChainID) {
	// TODO
}

func FetchNextNonce(ctx context.ImmutableContext, id types.LikeChainID) uint64 {
	return 0 // TODO
}
