package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/authz"
)

var _ authz.Authorization = &UpdateRoyaltyConfigAuthorization{}
var _ authz.Authorization = &UpdateListingAuthorization{}
var _ authz.Authorization = &UpdateOfferAuthorization{}
var _ authz.Authorization = &UpdateClassAuthorization{}
var _ authz.Authorization = &MintNFTAuthorization{}

func NewUpdateRoyaltyConfigAuthorization(classId string) *UpdateRoyaltyConfigAuthorization {
	return &UpdateRoyaltyConfigAuthorization{
		ClassId: classId,
	}
}

func (a UpdateRoyaltyConfigAuthorization) MsgTypeURL() string {
	return sdk.MsgTypeURL(&MsgUpdateRoyaltyConfig{})
}

func (a UpdateRoyaltyConfigAuthorization) Accept(ctx sdk.Context, msg sdk.Msg) (authz.AcceptResponse, error) {
	msgUpdate, ok := msg.(*MsgUpdateRoyaltyConfig)
	if !ok {
		return authz.AcceptResponse{}, sdkerrors.ErrInvalidType.Wrap("type mismatch")
	}
	if msgUpdate.ClassId != a.ClassId {
		return authz.AcceptResponse{}, sdkerrors.ErrUnauthorized.Wrap("class ID mismatch")
	}
	return authz.AcceptResponse{Accept: true}, nil
}

func (a UpdateRoyaltyConfigAuthorization) ValidateBasic() error {
	return nil
}

func NewUpdateListingAuthorization(classId string, nftId string) *UpdateListingAuthorization {
	return &UpdateListingAuthorization{
		ClassId: classId,
		NftId:   nftId,
	}
}

func (a UpdateListingAuthorization) MsgTypeURL() string {
	return sdk.MsgTypeURL(&MsgUpdateListing{})
}

func (a UpdateListingAuthorization) Accept(ctx sdk.Context, msg sdk.Msg) (authz.AcceptResponse, error) {
	msgUpdate, ok := msg.(*MsgUpdateListing)
	if !ok {
		return authz.AcceptResponse{}, sdkerrors.ErrInvalidType.Wrap("type mismatch")
	}
	if msgUpdate.ClassId != a.ClassId {
		return authz.AcceptResponse{}, sdkerrors.ErrUnauthorized.Wrap("class ID mismatch")
	}
	if msgUpdate.NftId != a.NftId {
		return authz.AcceptResponse{}, sdkerrors.ErrUnauthorized.Wrap("NFT ID mismatch")
	}
	return authz.AcceptResponse{Accept: true}, nil
}

func (a UpdateListingAuthorization) ValidateBasic() error {
	return nil
}

func NewUpdateOfferAuthorization(classId string, nftId string) *UpdateOfferAuthorization {
	return &UpdateOfferAuthorization{
		ClassId: classId,
		NftId:   nftId,
	}
}

func (a UpdateOfferAuthorization) MsgTypeURL() string {
	return sdk.MsgTypeURL(&MsgUpdateOffer{})
}

func (a UpdateOfferAuthorization) Accept(ctx sdk.Context, msg sdk.Msg) (authz.AcceptResponse, error) {
	msgUpdate, ok := msg.(*MsgUpdateOffer)
	if !ok {
		return authz.AcceptResponse{}, sdkerrors.ErrInvalidType.Wrap("type mismatch")
	}
	if msgUpdate.ClassId != a.ClassId {
		return authz.AcceptResponse{}, sdkerrors.ErrUnauthorized.Wrap("class ID mismatch")
	}
	if msgUpdate.NftId != a.NftId {
		return authz.AcceptResponse{}, sdkerrors.ErrUnauthorized.Wrap("NFT ID mismatch")
	}
	return authz.AcceptResponse{Accept: true}, nil
}

func (a UpdateOfferAuthorization) ValidateBasic() error {
	return nil
}

func NewUpdateClassAuthorization(classId string) *UpdateClassAuthorization {
	return &UpdateClassAuthorization{
		ClassId: classId,
	}
}

func (a UpdateClassAuthorization) MsgTypeURL() string {
	return sdk.MsgTypeURL(&MsgUpdateClass{})
}

func (a UpdateClassAuthorization) Accept(ctx sdk.Context, msg sdk.Msg) (authz.AcceptResponse, error) {
	msgUpdate, ok := msg.(*MsgUpdateClass)
	if !ok {
		return authz.AcceptResponse{}, sdkerrors.ErrInvalidType.Wrap("type mismatch")
	}
	if msgUpdate.ClassId != a.ClassId {
		return authz.AcceptResponse{}, sdkerrors.ErrUnauthorized.Wrap("class ID mismatch")
	}
	return authz.AcceptResponse{Accept: true}, nil
}

func (a UpdateClassAuthorization) ValidateBasic() error {
	return nil
}

func NewMintNFTAuthorization(classId string) *MintNFTAuthorization {
	return &MintNFTAuthorization{
		ClassId: classId,
	}
}

func (a MintNFTAuthorization) MsgTypeURL() string {
	return sdk.MsgTypeURL(&MsgMintNFT{})
}

func (a MintNFTAuthorization) Accept(ctx sdk.Context, msg sdk.Msg) (authz.AcceptResponse, error) {
	msgMint, ok := msg.(*MsgMintNFT)
	if !ok {
		return authz.AcceptResponse{}, sdkerrors.ErrInvalidType.Wrap("type mismatch")
	}
	if msgMint.ClassId != a.ClassId {
		return authz.AcceptResponse{}, sdkerrors.ErrUnauthorized.Wrap("class ID mismatch")
	}
	return authz.AcceptResponse{Accept: true}, nil
}

func (a MintNFTAuthorization) ValidateBasic() error {
	return nil
}
