package types

import (
	"github.com/tendermint/go-amino"
)

var cdc = amino.NewCodec()

func init() {
	cdc.RegisterInterface((*Identifier)(nil), nil)
	cdc.RegisterConcrete(&Address{}, "github.com/likecoin/likechain/Address", nil)
	cdc.RegisterConcrete(&LikeChainID{}, "github.com/likecoin/likechain/LikeChainID", nil)
}

// AminoCodec returns the amino Codec
func AminoCodec() *amino.Codec {
	return cdc
}
