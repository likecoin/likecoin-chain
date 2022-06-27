package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgUpdateBlindBoxContent = "update_blind_box_content"

var _ sdk.Msg = &MsgUpdateBlindBoxContent{}

func NewMsgUpdateBlindBoxContent(creator string, classId string, id string, input NFTInput) *MsgUpdateBlindBoxContent {
	return &MsgUpdateBlindBoxContent{
		Creator: creator,
		ClassId: classId,
		Id:      id,
		Input:   input,
	}
}

func (msg *MsgUpdateBlindBoxContent) Route() string {
	return RouterKey
}

func (msg *MsgUpdateBlindBoxContent) Type() string {
	return TypeMsgUpdateBlindBoxContent
}

func (msg *MsgUpdateBlindBoxContent) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdateBlindBoxContent) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateBlindBoxContent) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
