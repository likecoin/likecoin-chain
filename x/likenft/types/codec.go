package types

import (
	"github.com/gogo/protobuf/proto"

	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	"github.com/cosmos/cosmos-sdk/x/authz"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgNewClass{}, "likenft/NewClass", nil)
	cdc.RegisterConcrete(&MsgUpdateClass{}, "likenft/UpdateClass", nil)
	cdc.RegisterConcrete(&MsgMintNFT{}, "likenft/MintNFT", nil)
	cdc.RegisterConcrete(&MsgBurnNFT{}, "likenft/BurnNFT", nil)
	cdc.RegisterConcrete(&MsgCreateBlindBoxContent{}, "likenft/CreateBlindBoxContent", nil)
	cdc.RegisterConcrete(&MsgUpdateBlindBoxContent{}, "likenft/UpdateBlindBoxContent", nil)
	cdc.RegisterConcrete(&MsgDeleteBlindBoxContent{}, "likenft/DeleteBlindBoxContent", nil)
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

	cdc.RegisterConcrete(&CreateRoyaltyConfigAuthorization{}, "likenft/CreateRoyaltyConfigAuthorization", nil)
	cdc.RegisterConcrete(&UpdateRoyaltyConfigAuthorization{}, "likenft/UpdateRoyaltyConfigAuthorization", nil)
	cdc.RegisterConcrete(&DeleteRoyaltyConfigAuthorization{}, "likenft/DeleteRoyaltyConfigAuthorization", nil)
	cdc.RegisterConcrete(&CreateListingAuthorization{}, "likenft/CreateListingAuthorization", nil)
	cdc.RegisterConcrete(&UpdateListingAuthorization{}, "likenft/UpdateListingAuthorization", nil)
	cdc.RegisterConcrete(&DeleteListingAuthorization{}, "likenft/DeleteListingAuthorization", nil)
	cdc.RegisterConcrete(&CreateOfferAuthorization{}, "likenft/CreateOfferAuthorization", nil)
	cdc.RegisterConcrete(&UpdateOfferAuthorization{}, "likenft/UpdateOfferAuthorization", nil)
	cdc.RegisterConcrete(&DeleteOfferAuthorization{}, "likenft/DeleteOfferAuthorization", nil)
	cdc.RegisterConcrete(&NewClassAuthorization{}, "likenft/NewClassAuthorization", nil)
	cdc.RegisterConcrete(&UpdateClassAuthorization{}, "likenft/UpdateClassAuthorization", nil)
	cdc.RegisterConcrete(&SendNFTAuthorization{}, "likenft/SendNFTAuthorization", nil)
	cdc.RegisterConcrete(&MintNFTAuthorization{}, "likenft/MintNFTAuthorization", nil)
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
		&MsgCreateBlindBoxContent{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUpdateBlindBoxContent{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgDeleteBlindBoxContent{},
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

	registry.RegisterImplementations(
		(*authz.Authorization)(nil),
		&CreateRoyaltyConfigAuthorization{},
		&UpdateRoyaltyConfigAuthorization{},
		&DeleteRoyaltyConfigAuthorization{},
		&CreateListingAuthorization{},
		&UpdateListingAuthorization{},
		&DeleteListingAuthorization{},
		&CreateOfferAuthorization{},
		&UpdateOfferAuthorization{},
		&DeleteOfferAuthorization{},
		&NewClassAuthorization{},
		&UpdateClassAuthorization{},
		&SendNFTAuthorization{},
		&MintNFTAuthorization{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
