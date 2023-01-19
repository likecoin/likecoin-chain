package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/authz"
)

var _ authz.Authorization = &NewClassAuthorization{}
var _ authz.Authorization = &UpdateClassAuthorization{}

// implement Authorization interface for NewClassAuthorization
func NewNewClassAuthorization(iscnIdPrefix string) *NewClassAuthorization {
	return &NewClassAuthorization{
		IscnIdPrefix: iscnIdPrefix,
	}
}

func (a NewClassAuthorization) MsgTypeURL() string {
	return sdk.MsgTypeURL(&MsgNewClass{})
}

func (a NewClassAuthorization) Accept(ctx sdk.Context, msg sdk.Msg) (authz.AcceptResponse, error) {
	msgNewClass, ok := msg.(*MsgNewClass)
	if !ok {
		return authz.AcceptResponse{}, sdkerrors.ErrInvalidType.Wrap("type mismatch")
	}
	if msgNewClass.Parent.Type != ClassParentType_ISCN {
		return authz.AcceptResponse{}, sdkerrors.ErrUnauthorized.Wrap("unsupported parent type")
	}
	if msgNewClass.Parent.IscnIdPrefix != a.IscnIdPrefix {
		return authz.AcceptResponse{}, sdkerrors.ErrUnauthorized.Wrap("ISCN ID prefix mismatch")
	}
	return authz.AcceptResponse{Accept: true}, nil
}

func (a NewClassAuthorization) ValidateBasic() error {
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
