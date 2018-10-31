package types

import (
	"errors"
	"math/big"
)

// BigInt is an adaptor of big.Int, implementing AminoMarshaler and AminoUnmarshaler
type BigInt struct {
	*big.Int
}

// NewBigInt returns a BigInt initialized with an int64
func NewBigInt(n int64) BigInt {
	return BigInt{big.NewInt(n)}
}

// MarshalAmino implements AminoMarshaler
func (n BigInt) MarshalAmino() ([]byte, error) {
	return n.Bytes(), nil
}

// UnmarshalAmino implements AminoUnmarshaler
func (n *BigInt) UnmarshalAmino(bs []byte) error {
	n.Int = new(big.Int).SetBytes(bs)
	return nil
}

// MarshalJSON implements json.Marshaler
func (n *BigInt) MarshalJSON() ([]byte, error) {
	return []byte(`"` + n.Int.String() + `"`), nil
}

// UnmarshalJSON implements json.Unmarshaler
func (n *BigInt) UnmarshalJSON(bs []byte) error {
	if len(bs) > 2 && bs[0] == '"' || bs[len(bs)-1] == '"' {
		bs = bs[1 : len(bs)-1]
	}
	v, ok := new(big.Int).SetString(string(bs), 10)
	if !ok {
		return errors.New("Cannot parse BigInt string")
	}
	n.Int = v
	return nil
}

// NewBigIntFromString returns a BigInt from the input string with base 10, and a boolean indicates success.
// It fails when the string is an invalid number or when the bytes of the number cannot reproduce the number
// (e.g. negative numbers)
func NewBigIntFromString(s string) (BigInt, bool) {
	n, ok := new(big.Int).SetString(s, 10)
	if !ok || n.Cmp(new(big.Int).SetBytes(n.Bytes())) != 0 {
		return BigInt{}, false
	}
	return BigInt{n}, true
}
