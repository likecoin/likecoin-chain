package account

import (
	"encoding/binary"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/likecoin/likechain/abci/context"
	logger "github.com/likecoin/likechain/abci/log"
	"github.com/likecoin/likechain/abci/types"
	"github.com/likecoin/likechain/abci/utils"
)

var log = logger.L

// NewAccount creates a new account
func NewAccount(state context.IMutableState, ethAddr common.Address) (types.LikeChainID, error) {
	id := generateLikeChainID(state)

	// Save address mapping
	state.MutableStateTree().Set(utils.DbAddrKey(ethAddr), id.Content)
	state.MutableStateTree().Set(utils.DbIDKey(id, "acc", "addr"), ethAddr.Bytes())

	// Initialize account info
	SaveBalance(state, id, big.NewInt(0))
	IncrementNextNonce(state, id)

	// Check if address already has balance
	balanceInEthAddr := fetchEthereumAddressBalance(state, ethAddr)
	if balanceInEthAddr != nil {
		// TODO: Transfer balance from ETH address to LikeChain ID
	}

	return id, nil // TODO
}

var likeChainIDSeedKey = []byte("$account.likeChainIDSeed")

func generateLikeChainID(state context.IMutableState) types.LikeChainID {
	var seedInt uint64
	_, seed := state.ImmutableStateTree().Get(likeChainIDSeedKey)
	if seed == nil {
		seedInt = 1
		seed = make([]byte, 8)
		binary.BigEndian.PutUint64(seed, seedInt)
	} else {
		seedInt = uint64(binary.BigEndian.Uint64(seed))
	}

	blockHash := state.GetBlockHash()

	// Concat the seed and the block's hash
	content := make([]byte, len(seed)+len(blockHash))
	copy(content, seed)
	copy(content[len(seed):], blockHash)
	// Take first 20 bytes of Keccak256 hash to be LikeChainID
	content = crypto.Keccak256(content)[:20]

	// Increment and save seed
	seedInt++
	binary.BigEndian.PutUint64(seed, seedInt)
	state.MutableStateTree().Set(likeChainIDSeedKey, seed)

	return types.LikeChainID{Content: content}
}

func AddressToLikeChainID(state context.IImmutableState, addr types.Address) (types.LikeChainID, bool) {
	return types.LikeChainID{}, false // TODO
}

func GetLikeChainID(state context.IImmutableState, identifier types.Identifier) (types.LikeChainID, bool) {
	id := identifier.GetLikeChainID()
	if id != nil {
		return *id, false // TODO: check the existence of this LikeChainID
	}
	addr := identifier.GetAddr()
	// TODO: assert addr != nil
	return AddressToLikeChainID(state, *addr)
}

// SaveBalance saves account balance by LikeChain ID
func SaveBalance(state context.IMutableState, id types.LikeChainID, balance *big.Int) error {
	state.MutableStateTree().Set(utils.DbIDKey(id, "acc", "balance"), balance.Bytes())
	return nil
}

// FetchBalance fetches account balance by LikeChain ID
func FetchBalance(state context.IImmutableState, id types.LikeChainID) *big.Int {
	_, bytes := state.ImmutableStateTree().Get(utils.DbIDKey(id, "acc", "balance"))

	balance := big.NewInt(0)
	balance = balance.SetBytes(bytes)

	return balance
}

func fetchEthereumAddressBalance(state context.IMutableState, addr common.Address) *big.Int {
	return nil // TODO
}

// IncrementNextNonce increments next nonce of an account by LikeChain ID
// This also initialize next nonce of an account
func IncrementNextNonce(state context.IMutableState, id types.LikeChainID) {
	nextNonceInt := FetchNextNonce(state, id) + 1
	nextNonce := make([]byte, 8)
	binary.BigEndian.PutUint64(nextNonce, nextNonceInt)
	state.MutableStateTree().Set(utils.DbIDKey(id, "acc", "nextNonce"), nextNonce)
}

// FetchNextNonce fetches next nonce of an account by LikeChain ID
func FetchNextNonce(state context.IImmutableState, id types.LikeChainID) uint64 {
	_, bytes := state.ImmutableStateTree().Get(utils.DbIDKey(id, "acc", "nextNonce"))

	if bytes == nil {
		return uint64(0)
	}
	return uint64(binary.BigEndian.Uint64(bytes))
}
