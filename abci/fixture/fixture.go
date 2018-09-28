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
	id, _ := types.NewLikeChainIDFromString(idStr)
	return &User{
		ID:         id,
		RawAddress: types.NewAddressFromHex(addrHex),
		Address:    common.HexToAddress(addrHex),
	}
}

// Alice is a generic first participant
var Alice = NewUser(
	"YWxpY2VfX19fX19fX19fX19fX18=",
	"0x064b663abf9d74277a07aa7563a8a64a54de8c0a",
)

// Bob is a generic second participant
var Bob = NewUser(
	"Ym9iX19fX19fX19fX19fX19fX18=",
	"0xbef509a0ab4a60111a8957322fee016cdf713ad2",
)

// Carol is a generic third participant
var Carol = NewUser(
	"Y2Fyb2xfX19fX19fX19fX19fX18=",
	"0xba0ad74ab6cfea30e0cfa4998392873ad1a11388",
)

// Mallory is a malicious participant
var Mallory = NewUser(
	"Y2Fyb2xfX19fX19fX19fX19fX18=",
	"0x65f86d54c5e768efe89dd5d07143fd783a3303df",
)
