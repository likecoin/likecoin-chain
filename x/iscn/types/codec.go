package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgCreateIscnRecord{}, "likecoin-chain/MsgCreateIscnRecord", nil)
	cdc.RegisterConcrete(&MsgUpdateIscnRecord{}, "likecoin-chain/MsgUpdateIscnRecord", nil)
	cdc.RegisterConcrete(&MsgChangeIscnRecordOwnership{}, "likecoin-chain/MsgChangeIscnRecordOwnership", nil)
}

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateIscnRecord{},
		&MsgUpdateIscnRecord{},
		&MsgChangeIscnRecordOwnership{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewAminoCodec(amino)
)

func init() {
	RegisterLegacyAminoCodec(amino)
	cryptocodec.RegisterCrypto(amino)
	amino.Seal()
}
