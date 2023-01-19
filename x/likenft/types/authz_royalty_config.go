package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/authz"
)

var _ authz.Authorization = &CreateRoyaltyConfigAuthorization{}
var _ authz.Authorization = &UpdateRoyaltyConfigAuthorization{}
var _ authz.Authorization = &DeleteRoyaltyConfigAuthorization{}

func NewCreateRoyaltyConfigAuthorization(classId string) *CreateRoyaltyConfigAuthorization {
	return &CreateRoyaltyConfigAuthorization{
		ClassId: classId,
	}
}

func (a CreateRoyaltyConfigAuthorization) MsgTypeURL() string {
	return sdk.MsgTypeURL(&MsgCreateRoyaltyConfig{})
}

func (a CreateRoyaltyConfigAuthorization) Accept(ctx sdk.Context, msg sdk.Msg) (authz.AcceptResponse, error) {
	msgCreate, ok := msg.(*MsgCreateRoyaltyConfig)
	if !ok {
		return authz.AcceptResponse{}, sdkerrors.ErrInvalidType.Wrap("type mismatch")
	}
	if msgCreate.ClassId != a.ClassId {
		return authz.AcceptResponse{}, sdkerrors.ErrUnauthorized.Wrap("class ID mismatch")
	}
	return authz.AcceptResponse{Accept: true}, nil
}

func (a CreateRoyaltyConfigAuthorization) ValidateBasic() error {
	return nil
}

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

func NewDeleteRoyaltyConfigAuthorization(classId string) *DeleteRoyaltyConfigAuthorization {
	return &DeleteRoyaltyConfigAuthorization{
		ClassId: classId,
	}
}

func (a DeleteRoyaltyConfigAuthorization) MsgTypeURL() string {
	return sdk.MsgTypeURL(&MsgDeleteRoyaltyConfig{})
}

func (a DeleteRoyaltyConfigAuthorization) Accept(ctx sdk.Context, msg sdk.Msg) (authz.AcceptResponse, error) {
	msgDelete, ok := msg.(*MsgDeleteRoyaltyConfig)
	if !ok {
		return authz.AcceptResponse{}, sdkerrors.ErrInvalidType.Wrap("type mismatch")
	}
	if msgDelete.ClassId != a.ClassId {
		return authz.AcceptResponse{}, sdkerrors.ErrUnauthorized.Wrap("class ID mismatch")
	}
	return authz.AcceptResponse{Accept: true}, nil
}

func (a DeleteRoyaltyConfigAuthorization) ValidateBasic() error {
	return nil
}
