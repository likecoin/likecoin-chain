package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/authz"
)

var _ authz.Authorization = &CreateListingAuthorization{}
var _ authz.Authorization = &UpdateListingAuthorization{}
var _ authz.Authorization = &DeleteListingAuthorization{}

func NewCreateListingAuthorization(classId string, nftId string) *CreateListingAuthorization {
	return &CreateListingAuthorization{
		ClassId: classId,
		NftId:   nftId,
	}
}

func (a CreateListingAuthorization) MsgTypeURL() string {
	return sdk.MsgTypeURL(&MsgCreateListing{})
}

func (a CreateListingAuthorization) Accept(ctx sdk.Context, msg sdk.Msg) (authz.AcceptResponse, error) {
	msgCreate, ok := msg.(*MsgCreateListing)
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

func (a CreateListingAuthorization) ValidateBasic() error {
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
	if a.NftId != "" && msgUpdate.NftId != a.NftId {
		return authz.AcceptResponse{}, sdkerrors.ErrUnauthorized.Wrap("NFT ID mismatch")
	}
	return authz.AcceptResponse{Accept: true}, nil
}

func (a UpdateListingAuthorization) ValidateBasic() error {
	return nil
}

func NewDeleteListingAuthorization(classId string, nftId string) *DeleteListingAuthorization {
	return &DeleteListingAuthorization{
		ClassId: classId,
		NftId:   nftId,
	}
}

func (a DeleteListingAuthorization) MsgTypeURL() string {
	return sdk.MsgTypeURL(&MsgDeleteListing{})
}

func (a DeleteListingAuthorization) Accept(ctx sdk.Context, msg sdk.Msg) (authz.AcceptResponse, error) {
	msgDelete, ok := msg.(*MsgDeleteListing)
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

func (a DeleteListingAuthorization) ValidateBasic() error {
	return nil
}
