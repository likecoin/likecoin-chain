package types

// Identifier is either a LikeChain ID or an Ethereum address
type Identifier interface {
	Equals(iden Identifier) bool
	Bytes() []byte
	String() string
	DBKey(prefix string, suffix string) []byte
}

// NewIdentifier constructs and returns either a LikeChainID or an Address
func NewIdentifier(s string) Identifier {
	id, err := NewLikeChainIDFromString(s)
	if err == nil {
		return id
	}
	addr, err := NewAddressFromHex(s)
	if err == nil {
		return addr
	}
	return nil
}
