package account

import (
	"bytes"
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

func getIDAddrPairPrefixKey(id *types.LikeChainID) []byte {
	var buf bytes.Buffer
	buf.WriteString("acc_")
	buf.Write(id.Content)
	buf.WriteString("_addr_")
	return buf.Bytes()
}

func getIDAddrPairKey(id *types.LikeChainID, ethAddr common.Address) []byte {
	var buf bytes.Buffer
	buf.Write(getIDAddrPairPrefixKey(id))
	buf.Write(ethAddr.Bytes())
	return buf.Bytes()
}

func getAddrIDPairKey(ethAddr common.Address) []byte {
	return utils.DbRawKey(ethAddr.Bytes(), "acc", "id")
}

// NewAccount creates a new account
func NewAccount(state context.IMutableState, ethAddr common.Address) (*types.LikeChainID, error) {
	id := generateLikeChainID(state)
	err := NewAccountFromID(state, id, ethAddr)
	return id, err
}

// NewAccountFromID creates a new account from a given LikeChain ID
func NewAccountFromID(state context.IMutableState, id *types.LikeChainID, ethAddr common.Address) error {
	// Save address mapping
	state.MutableStateTree().Set(getAddrIDPairKey(ethAddr), id.Content)
	state.MutableStateTree().Set(getIDAddrPairKey(id, ethAddr), []byte{})

	// Check if address already has balance
	addrIdentifier := types.NewAddressFromHex(ethAddr.Hex()).ToIdentifier()
	addrBalance := FetchRawBalance(state, addrIdentifier)

	var balance *big.Int
	if addrBalance.Cmp(big.NewInt(0)) > 0 {
		// Transfer balance to LikeChain ID
		balance = addrBalance

		// Remove key from db
		key := utils.DbIdentifierKey(addrIdentifier, "acc", "balance")
		state.MutableStateTree().Remove(key)
	} else {
		balance = big.NewInt(0)
	}

	// Initialize account info
	SaveBalance(state, id.ToIdentifier(), balance)
	IncrementNextNonce(state, id)

	return nil
}

func iterateLikeChainIDAddrPair(state context.IImmutableState, id *types.LikeChainID, fn func(id, addr []byte) bool) (isExist bool) {
	startingKey := getIDAddrPairPrefixKey(id)

	// Iterate the tree to check all addresses the given LikeChain ID has been bound
	state.ImmutableStateTree().IterateRange(startingKey, nil, true, func(key, _ []byte) bool {
		splitKey := bytes.Split(key, []byte("acc_"))
		splitKey = bytes.Split(splitKey[1], []byte("_addr_"))
		if len(splitKey) == 2 {
			isExist = fn(splitKey[0], splitKey[1])
		}
		// If isExist becomes true, iteration will be stopped
		return isExist
	})

	return isExist
}

// IsLikeChainIDRegistered checks whether the given LikeChain ID has registered or not
func IsLikeChainIDRegistered(state context.IImmutableState, id *types.LikeChainID) bool {
	return iterateLikeChainIDAddrPair(state, id, func(idBytes, _ []byte) bool {
		return bytes.Compare(idBytes, id.Content) == 0
	})
}

// IsAddressRegistered checks whether the given Address has registered or not
func IsAddressRegistered(state context.IImmutableState, ethAddr common.Address) bool {
	_, value := state.ImmutableStateTree().Get(getAddrIDPairKey(ethAddr))
	return value != nil
}

// IsLikeChainIDHasAddress checks whether the given address has been bound to the given LikeChain ID
func IsLikeChainIDHasAddress(state context.IImmutableState, id *types.LikeChainID, ethAddr common.Address) bool {
	_, value := state.ImmutableStateTree().Get(getIDAddrPairKey(id, ethAddr))
	return value != nil
}

var likeChainIDSeedKey = []byte("$account.likeChainIDSeed")

func generateLikeChainID(state context.IMutableState) *types.LikeChainID {
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

	return &types.LikeChainID{Content: content}
}

// AddressToLikeChainID gets LikeChain ID by Address
func AddressToLikeChainID(state context.IImmutableState, ethAddr common.Address) *types.LikeChainID {
	_, value := state.ImmutableStateTree().Get(getAddrIDPairKey(ethAddr))
	if value != nil {
		return types.NewLikeChainID(value)
	}
	return nil
}

// IdentifierToLikeChainID converts a Identifier to LikeChain ID using
// address - LikeChain ID mapping
func IdentifierToLikeChainID(state context.IImmutableState, identifier *types.Identifier) *types.LikeChainID {
	id := identifier.GetLikeChainID()
	if id != nil && IsLikeChainIDRegistered(state, id) {
		return id
	} else if addr := identifier.GetAddr(); addr != nil {
		return AddressToLikeChainID(state, addr.ToEthereum())
	}

	return nil
}

// NormalizeIdentifier converts an identifier with an address to an identifier
// with LikeChain ID if the address has registered
func NormalizeIdentifier(
	state context.IImmutableState,
	identifier *types.Identifier,
) *types.Identifier {
	if addr := identifier.GetAddr(); addr != nil {
		id := AddressToLikeChainID(state, addr.ToEthereum())
		if id != nil {
			return id.ToIdentifier()
		}
	}
	return identifier
}

// SaveBalance saves account balance by LikeChain ID
func SaveBalance(state context.IMutableState, identifier *types.Identifier, balance *big.Int) error {
	key := utils.DbIdentifierKey(
		NormalizeIdentifier(state, identifier), "acc", "balance")

	state.MutableStateTree().Set(key, balance.Bytes())
	return nil
}

// FetchBalance fetches account balance by normalized Identifier
func FetchBalance(state context.IImmutableState, identifier *types.Identifier) *big.Int {
	return FetchRawBalance(state, NormalizeIdentifier(state, identifier))
}

// FetchRawBalance fetches account balance by Identifier
func FetchRawBalance(state context.IImmutableState, identifier *types.Identifier) *big.Int {
	key := utils.DbIdentifierKey(identifier, "acc", "balance")

	_, value := state.ImmutableStateTree().Get(key)

	balance := big.NewInt(0)
	balance = balance.SetBytes(value)

	return balance
}

// AddBalance adds account balance by Identifier
func AddBalance(state context.IMutableState, identifier *types.Identifier, amount *big.Int) error {
	balance := FetchBalance(state, identifier)
	balance.Add(balance, amount)
	return SaveBalance(state, identifier, balance)
}

// MinusBalance minus account balance by Identifier
func MinusBalance(state context.IMutableState, identifier *types.Identifier, amount *big.Int) error {
	balance := FetchBalance(state, identifier)
	balance.Sub(balance, amount)
	return SaveBalance(state, identifier, balance)
}

// IncrementNextNonce increments next nonce of an account by LikeChain ID
// This also initialize next nonce of an account
func IncrementNextNonce(state context.IMutableState, id *types.LikeChainID) {
	nextNonceInt := FetchNextNonce(state, id) + 1
	nextNonce := make([]byte, 8)
	binary.BigEndian.PutUint64(nextNonce, nextNonceInt)
	state.MutableStateTree().Set(utils.DbIDKey(id, "acc", "nextNonce"), nextNonce)
}

// FetchNextNonce fetches next nonce of an account by LikeChain ID
func FetchNextNonce(state context.IImmutableState, id *types.LikeChainID) uint64 {
	_, bytes := state.ImmutableStateTree().Get(utils.DbIDKey(id, "acc", "nextNonce"))

	if bytes == nil {
		return uint64(0)
	}
	return uint64(binary.BigEndian.Uint64(bytes))
}
