package types

import (
	"bytes"
	"errors"
	"strings"

	"github.com/likecoin/likechain/abci/utils"

	"github.com/ethereum/go-ethereum/common"
)

// Address represents an Ethereum address
type Address common.Address

// Equals returns true if the other identifier is exactly the same address as the receiver
func (addr *Address) Equals(iden Identifier) bool {
	switch iden.(type) {
	case *Address:
		addr2 := iden.(*Address)
		return bytes.Compare(addr[:], addr2[:]) == 0
	default:
		return false
	}
}

// Bytes returns the bytes of the Address
func (addr *Address) Bytes() []byte {
	return addr[:]
}

// DBKey returns a key with Ethereum address in `{prefix}:addr:_{addr}_{suffix}` format
func (addr *Address) DBKey(prefix string, suffix string) []byte {
	var buf bytes.Buffer
	buf.WriteString(prefix)
	buf.WriteString(":addr:")
	return utils.DbRawKey(addr.Bytes(), buf.String(), suffix)
}

// NewAddress creates an Address from bytes
func NewAddress(bs []byte) (*Address, error) {
	if len(bs) != 20 {
		return nil, errors.New("Invalid Address length")
	}
	result := Address{}
	copy(result[:], bs)
	return &result, nil
}

// NewAddressFromHex creates an Address from hex string
func NewAddressFromHex(s string) (*Address, error) {
	bs, err := utils.Hex2Bytes(s)
	if err != nil {
		return nil, err
	}
	return NewAddress(bs)
}

func (addr *Address) String() string {
	return strings.ToLower(common.Address(*addr).Hex())
}

// MarshalJSON implements json.Marshaler
func (addr *Address) MarshalJSON() ([]byte, error) {
	return []byte(`"` + addr.String() + `"`), nil
}

// UnmarshalJSON implements json.Unmarshaler
func (addr *Address) UnmarshalJSON(bs []byte) error {
	if len(bs) < 2 || bs[0] != '"' || bs[len(bs)-1] != '"' {
		return errors.New("Invalid input for Address JSON serialization data")
	}
	bs = bs[1 : len(bs)-1]
	tmpAddr, err := NewAddressFromHex(string(bs))
	if err != nil {
		return err
	}
	*addr = *tmpAddr
	return nil
}

// Addr transforms a hex string into address, panic if the string is not a valid address
func Addr(s string) *Address {
	addr, err := NewAddressFromHex(s)
	if err != nil {
		panic(err)
	}
	return addr
}
