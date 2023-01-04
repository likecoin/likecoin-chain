package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/authz"
)

var _ authz.Authorization = &UpdateAuthorization{}

func NewUpdateAuthorization(iscnIdPrefix string) *UpdateAuthorization {
	return &UpdateAuthorization{
		IscnIdPrefix: iscnIdPrefix,
	}
}

func (a UpdateAuthorization) MsgTypeURL() string {
	return sdk.MsgTypeURL(&MsgUpdateIscnRecord{})
}

func (a UpdateAuthorization) Accept(ctx sdk.Context, msg sdk.Msg) (authz.AcceptResponse, error) {
	msgUpdate, ok := msg.(*MsgUpdateIscnRecord)
	if !ok {
		return authz.AcceptResponse{}, sdkerrors.ErrInvalidType.Wrap("type mismatch")
	}
	iscnId, err := ParseIscnId(msgUpdate.IscnId)
	if err != nil {
		return authz.AcceptResponse{}, ErrInvalidIscnId.Wrapf("%v", err)
	}
	authIscnIdPrefix, err := ParseIscnId(a.IscnIdPrefix)
	if err != nil {
		return authz.AcceptResponse{}, sdkerrors.ErrLogic.Wrapf("authorization has invalid ISCN ID prefix: %v", err)
	}
	if !iscnId.PrefixEqual(&authIscnIdPrefix) {
		return authz.AcceptResponse{}, sdkerrors.ErrUnauthorized.Wrap("ISCN ID prefix mismatch")
	}

	return authz.AcceptResponse{Accept: true}, nil
}

func (a UpdateAuthorization) ValidateBasic() error {
	_, err := ParseIscnId(a.IscnIdPrefix)
	if err != nil {
		return ErrInvalidIscnId.Wrapf("%v", err)
	}
	return nil
}
