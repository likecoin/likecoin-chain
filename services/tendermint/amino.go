package tendermint

import (
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

var cdc = amino.NewCodec()

// AminoCodec returns the Codec struct shared among packages
func AminoCodec() *amino.Codec {
	return cdc
}

func init() {
	cdc.RegisterInterface((*crypto.Signature)(nil), nil)
	cdc.RegisterConcrete(secp256k1.SignatureSecp256k1{}, "tendermint/SignatureSecp256k1", nil)
	cdc.RegisterInterface((*crypto.PubKey)(nil), nil)
	cdc.RegisterConcrete(secp256k1.PubKeySecp256k1{}, "tendermint/PubKeySecp256k1", nil)
}
