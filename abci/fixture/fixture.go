package fixture

import (
	"github.com/likecoin/likechain/abci/types"
)

// User is a person in LikeChain
type User struct {
	ID      *types.LikeChainID
	Address *types.Address
}

func user(idStr string, addrHex string) *User {
	id, err := types.NewLikeChainIDFromString(idStr)
	if err != nil {
		panic(err)
	}
	addr, err := types.NewAddressFromHex(addrHex)
	if err != nil {
		panic(err)
	}
	return &User{
		ID:      id,
		Address: addr,
	}
}

// Alice is a generic first participant, private key: 4B5FA628ABAE47F8D329441DAE5F3B71775523913691C4EDF28AA2D3AFB760AD
var Alice = user(
	"YWxpY2VfX19fX19fX19fX19fX18=",
	"0x064b663abf9d74277a07aa7563a8a64a54de8c0a",
)

// Bob is a generic second participant, private key: 94AFAA67107754F93178942B4262A6092449A073FFE06BFAFF49B19A2E6ECB76
var Bob = user(
	"Ym9iX19fX19fX19fX19fX19fX18=",
	"0xbef509a0ab4a60111a8957322fee016cdf713ad2",
)

// Carol is a generic third participant, private key: 027F85F30CAA3F24AC1A1FD3315F5A8AE027862139BF66D626DF3A05FC26AC1C
var Carol = user(
	"Y2Fyb2xfX19fX19fX19fX19fX18=",
	"0xba0ad74ab6cfea30e0cfa4998392873ad1a11388",
)

// Dave is a generic fourth participant, private key: F7EF0EC3C7FA06EBED1FC93A633CEAE67CE1A8845C35503C0B5F700A33257FF0
var Dave = user(
	"ZGF2ZV9fX19fX19fX19fX19fX18=",
	"0xf6c45b1c4b73ccaeb1d9a37024a6b9fa711d7e7e",
)

// Erin is a generic fifth participant, private key: 66C86947436511D9A931BEEF5B3783066805AB637372BB08F880EF685691DE94
var Erin = user(
	"ZXJpbl9fX19fX19fX19fX19fX18=",
	"0xcc320404b90901fc22401a3c85ac28dfa7295f1e",
)

// Frank is a generic fifth participant, private key: 56D5B1E477E49420CADE0FF0A2755D725C1E4F74B43F61D0202F9C2743BB819E
var Frank = user(
	"ZnJhbmtfX19fX19fX19fX19fX18=",
	"0x342f67f183b4c00097af97cf04b00bfe30c6a4d7",
)

// Mallory is a malicious participant, private key: 4EDA7E263968014D527BA7EE639A9056CF509173C26D2B15F4720A0A2C02993F
var Mallory = user(
	"bWFsbG9yeV9fX19fX19fX19fX18=",
	"0x65f86d54c5e768efe89dd5d07143fd783a3303df",
)
