package account

import (
	"encoding/binary"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/likecoin/likechain/abci/context"
	"github.com/likecoin/likechain/abci/types"
	"github.com/likecoin/likechain/abci/utils"
)

// NewAccount creates a new account
func NewAccount(ctx context.MutableContext, ethAddr common.Address) (types.LikeChainID, error) {
	id := generateLikeChainID(ctx)

	// Save address mapping
	ctx.MutableStateTree().Set(utils.DbAddrKey(ethAddr), id.Content)
	ctx.MutableStateTree().Set(utils.DbIDKey(id, "acc", "addr"), ethAddr.Bytes())

	// Initialize account info
	SaveBalance(ctx, id, big.NewInt(0))
	IncrementNextNonce(ctx, id)

	// Check if address already has balance
	balanceInEthAddr := fetchEthereumAddressBalance(ctx, ethAddr)
	if balanceInEthAddr != nil {
		// TODO: Transfer balance from ETH address to LikeChain ID
	}

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

// SaveBalance saves account balance by LikeChain ID
func SaveBalance(ctx context.MutableContext, id types.LikeChainID, balance *big.Int) error {
	ctx.MutableStateTree().Set(utils.DbIDKey(id, "acc", "balance"), balance.Bytes())
	return nil
}

// FetchBalance fetches account balance by LikeChain ID
func FetchBalance(ctx context.ImmutableContext, id types.LikeChainID) *big.Int {
	_, bytes := ctx.StateTree().Get(utils.DbIDKey(id, "acc", "balance"))

	balance := big.NewInt(0)
	balance = balance.SetBytes(bytes)

	return balance
}

func fetchEthereumAddressBalance(ctx context.MutableContext, addr common.Address) *big.Int {
	return nil // TODO
}

// IncrementNextNonce increments next nonce of an account by LikeChain ID
// This also initialize next nonce of an account
func IncrementNextNonce(ctx context.MutableContext, id types.LikeChainID) {
	nextNonceInt := FetchNextNonce(ctx, id) + 1
	nextNonce := make([]byte, 8)
	binary.BigEndian.PutUint64(nextNonce, nextNonceInt)
	ctx.MutableStateTree().Set(utils.DbIDKey(id, "acc", "nextNonce"), nextNonce)
}

// FetchNextNonce fetches next nonce of an account by LikeChain ID
func FetchNextNonce(ctx context.ImmutableContext, id types.LikeChainID) uint64 {
	_, bytes := ctx.StateTree().Get(utils.DbIDKey(id, "acc", "nextNonce"))

	if bytes == nil {
		return uint64(0)
	}
	return uint64(binary.BigEndian.Uint64(bytes))
}
