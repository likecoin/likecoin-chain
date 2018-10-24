package types

import "math/big"

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
