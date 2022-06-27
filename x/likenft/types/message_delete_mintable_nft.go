package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgDeleteBlindBoxContent = "delete_blind_box_content"

var _ sdk.Msg = &MsgDeleteBlindBoxContent{}

func NewMsgDeleteBlindBoxContent(creator string, classId string, id string) *MsgDeleteBlindBoxContent {
	return &MsgDeleteBlindBoxContent{
		Creator: creator,
		ClassId: classId,
		Id:      id,
	}
}

func (msg *MsgDeleteBlindBoxContent) Route() string {
	return RouterKey
}

func (msg *MsgDeleteBlindBoxContent) Type() string {
	return TypeMsgDeleteBlindBoxContent
}

func (msg *MsgDeleteBlindBoxContent) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgDeleteBlindBoxContent) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgDeleteBlindBoxContent) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
