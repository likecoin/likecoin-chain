package types

import (
	"bytes"
	"encoding/base64"
	"errors"

	"github.com/likecoin/likechain/abci/utils"
)

// LikeChainID is the account primary key (except for address-only accounts) in LikeChain
type LikeChainID [20]byte

// Equals returns true if the other identifier is exactly the same LikeChainID as the receiver
func (id *LikeChainID) Equals(iden Identifier) bool {
	switch iden.(type) {
	case *LikeChainID:
		id2 := iden.(*LikeChainID)
		return bytes.Compare(id[:], id2[:]) == 0
	default:
		return false
	}
}

// Bytes returns the bytes of the LikeChainID
func (id *LikeChainID) Bytes() []byte {
	return id[:]
}

// DBKey composes a key with LikeChain ID in `{prefix}:id:_{id}_{suffix}` format
func (id *LikeChainID) DBKey(prefix string, suffix string) []byte {
	var buf bytes.Buffer
	buf.WriteString(prefix)
	buf.WriteString(":id:")
	return utils.DbRawKey(id.Bytes(), buf.String(), suffix)
}

// NewLikeChainID creates a LikeChainID from bytes
func NewLikeChainID(bs []byte) (*LikeChainID, error) {
	if len(bs) != 20 {
		return nil, errors.New("Invalid LikeChainID length")
	}
	result := LikeChainID{}
	copy(result[:], bs)
	return &result, nil
}

// NewLikeChainIDFromString creates a LikeChainID from string
func NewLikeChainIDFromString(s string) (*LikeChainID, error) {
	bs, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return nil, err
	}
	return NewLikeChainID(bs)
}

func (id *LikeChainID) String() string {
	return base64.StdEncoding.EncodeToString(id[:])
}

// MarshalJSON implements json.Marshaler
func (id *LikeChainID) MarshalJSON() ([]byte, error) {
	return []byte(`"` + id.String() + `"`), nil
}

// UnmarshalJSON implements json.Unmarshaler
func (id *LikeChainID) UnmarshalJSON(bs []byte) error {
	if len(bs) < 2 || bs[0] != '"' || bs[len(bs)-1] != '"' {
		return errors.New("Invalid input for LikeChainID JSON serialization data")
	}
	bs = bs[1 : len(bs)-1]
	tmpID, err := NewLikeChainIDFromString(string(bs))
	if err != nil {
		return err
	}
	*id = *tmpID
	return nil
}

// ID is the shortcut of types.NewLikeChainID
func ID(bs []byte) *LikeChainID {
	id, err := NewLikeChainID(bs)
	if err != nil {
		panic(err)
	}
	return id
}

// IDStr transforms a base-64 string into LikeChainID, panic if the string is not a valid LikeChainID
func IDStr(s string) *LikeChainID {
	id, err := NewLikeChainIDFromString(s)
	if err != nil {
		panic(err)
	}
	return id
}
