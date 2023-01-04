package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/authz"
)

var _ authz.Authorization = &CreateOfferAuthorization{}
var _ authz.Authorization = &UpdateOfferAuthorization{}
var _ authz.Authorization = &DeleteOfferAuthorization{}

func NewCreateOfferAuthorization(classId string, nftId string) *CreateOfferAuthorization {
	return &CreateOfferAuthorization{
		ClassId: classId,
		NftId:   nftId,
	}
}

func (a CreateOfferAuthorization) MsgTypeURL() string {
	return sdk.MsgTypeURL(&MsgCreateOffer{})
}

func (a CreateOfferAuthorization) Accept(ctx sdk.Context, msg sdk.Msg) (authz.AcceptResponse, error) {
	msgCreate, ok := msg.(*MsgCreateOffer)
	if !ok {
		return authz.AcceptResponse{}, sdkerrors.ErrInvalidType.Wrap("type mismatch")
	}
	if msgCreate.ClassId != a.ClassId {
		return authz.AcceptResponse{}, sdkerrors.ErrUnauthorized.Wrap("class ID mismatch")
	}
	if a.NftId != "" && msgCreate.NftId != a.NftId {
		return authz.AcceptResponse{}, sdkerrors.ErrUnauthorized.Wrap("NFT ID mismatch")
	}
	return authz.AcceptResponse{Accept: true}, nil
}

func (a CreateOfferAuthorization) ValidateBasic() error {
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
	if a.NftId != "" && msgUpdate.NftId != a.NftId {
		return authz.AcceptResponse{}, sdkerrors.ErrUnauthorized.Wrap("NFT ID mismatch")
	}
	return authz.AcceptResponse{Accept: true}, nil
}

func (a UpdateOfferAuthorization) ValidateBasic() error {
	return nil
}

// implement Authorization interface for DeleteOfferAuthorization
func NewDeleteOfferAuthorization(classId string, nftId string) *DeleteOfferAuthorization {
	return &DeleteOfferAuthorization{
		ClassId: classId,
		NftId:   nftId,
	}
}

func (a DeleteOfferAuthorization) MsgTypeURL() string {
	return sdk.MsgTypeURL(&MsgDeleteOffer{})
}

func (a DeleteOfferAuthorization) Accept(ctx sdk.Context, msg sdk.Msg) (authz.AcceptResponse, error) {
	msgDelete, ok := msg.(*MsgDeleteOffer)
	if !ok {
		return authz.AcceptResponse{}, sdkerrors.ErrInvalidType.Wrap("type mismatch")
	}
	if msgDelete.ClassId != a.ClassId {
		return authz.AcceptResponse{}, sdkerrors.ErrUnauthorized.Wrap("class ID mismatch")
	}
	if a.NftId != "" && msgDelete.NftId != a.NftId {
		return authz.AcceptResponse{}, sdkerrors.ErrUnauthorized.Wrap("NFT ID mismatch")
	}
	return authz.AcceptResponse{Accept: true}, nil
}

func (a DeleteOfferAuthorization) ValidateBasic() error {
	return nil
}
