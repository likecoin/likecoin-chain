package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	"github.com/gogo/protobuf/proto"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgNewClass{}, "likenft/NewClass", nil)
	cdc.RegisterConcrete(&MsgUpdateClass{}, "likenft/UpdateClass", nil)
	cdc.RegisterConcrete(&MsgMintNFT{}, "likenft/MintNFT", nil)
	cdc.RegisterConcrete(&MsgBurnNFT{}, "likenft/BurnNFT", nil)
	cdc.RegisterConcrete(&MsgCreateMintableNFT{}, "likenft/CreateMintableNFT", nil)
	cdc.RegisterConcrete(&MsgUpdateMintableNFT{}, "likenft/UpdateMintableNFT", nil)
	cdc.RegisterConcrete(&MsgDeleteMintableNFT{}, "likenft/DeleteMintableNFT", nil)
	cdc.RegisterConcrete(&MsgCreateOffer{}, "likenft/CreateOffer", nil)
	cdc.RegisterConcrete(&MsgUpdateOffer{}, "likenft/UpdateOffer", nil)
	cdc.RegisterConcrete(&MsgDeleteOffer{}, "likenft/DeleteOffer", nil)
	cdc.RegisterConcrete(&MsgCreateListing{}, "likenft/CreateListing", nil)
	cdc.RegisterConcrete(&MsgUpdateListing{}, "likenft/UpdateListing", nil)
	cdc.RegisterConcrete(&MsgDeleteListing{}, "likenft/DeleteListing", nil)
	cdc.RegisterConcrete(&MsgSellNFT{}, "likenft/SellNFT", nil)
	cdc.RegisterConcrete(&MsgBuyNFT{}, "likenft/BuyNFT", nil)
	cdc.RegisterConcrete(&MsgCreateRoyaltyConfig{}, "likenft/CreateRoyaltyConfig", nil)
	cdc.RegisterConcrete(&MsgUpdateRoyaltyConfig{}, "likenft/UpdateRoyaltyConfig", nil)
	cdc.RegisterConcrete(&MsgDeleteRoyaltyConfig{}, "likenft/DeleteRoyaltyConfig", nil)
	// this line is used by starport scaffolding # 2
	cdc.RegisterConcrete(&ClassData{}, "likenft/ClassData", nil)
	cdc.RegisterConcrete(&ClassParent{}, "likenft/ClassParent", nil)
	cdc.RegisterConcrete(&ClassConfig{}, "likenft/ClassConfig", nil)
	cdc.RegisterConcrete(&NFTData{}, "likenft/NFTData", nil)
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgNewClass{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUpdateClass{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgMintNFT{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgBurnNFT{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateMintableNFT{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUpdateMintableNFT{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgDeleteMintableNFT{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateOffer{},
		&MsgUpdateOffer{},
		&MsgDeleteOffer{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateListing{},
		&MsgUpdateListing{},
		&MsgDeleteListing{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgSellNFT{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgBuyNFT{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateRoyaltyConfig{},
		&MsgUpdateRoyaltyConfig{},
		&MsgDeleteRoyaltyConfig{},
	)
	// this line is used by starport scaffolding # 3
	registry.RegisterImplementations((*proto.Message)(nil), &ClassData{})
	registry.RegisterImplementations((*proto.Message)(nil), &ClassParent{})
	registry.RegisterImplementations((*proto.Message)(nil), &ClassConfig{})
	registry.RegisterImplementations((*proto.Message)(nil), &NFTData{})

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
