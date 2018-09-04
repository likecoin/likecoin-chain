package fixture

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/likecoin/likechain/abci/types"
)

// User is a person in LikeChain
type User struct {
	ID         *types.LikeChainID
	RawAddress *types.Address
	Address    common.Address
}

// NewUser creates an User with LikeChain ID and address in string
func NewUser(idStr string, addrHex string) *User {
	return &User{
		ID:         types.NewLikeChainID([]byte(idStr)),
		RawAddress: types.NewAddressFromHex(addrHex),
		Address:    common.HexToAddress(addrHex),
	}
}

// Alice is a generic first participant
var Alice = NewUser(
	"alice",
	"0x064b663abf9d74277a07aa7563a8a64a54de8c0a",
)

// Bob is a generic second participant
var Bob = NewUser(
	"bob",
	"0xbef509a0ab4a60111a8957322fee016cdf713ad2",
)

// Carol is a generic third participant
var Carol = NewUser(
	"carol",
	"0xba0ad74ab6cfea30e0cfa4998392873ad1a11388",
)
