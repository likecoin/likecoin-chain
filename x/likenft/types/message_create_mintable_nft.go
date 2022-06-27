package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgCreateBlindBoxContent = "create_mintable_nft"

var _ sdk.Msg = &MsgCreateBlindBoxContent{}

func NewMsgCreateBlindBoxContent(creator string, classId string, id string, input NFTInput) *MsgCreateBlindBoxContent {
	return &MsgCreateBlindBoxContent{
		Creator: creator,
		ClassId: classId,
		Id:      id,
		Input:   input,
	}
}

func (msg *MsgCreateBlindBoxContent) Route() string {
	return RouterKey
}

func (msg *MsgCreateBlindBoxContent) Type() string {
	return TypeMsgCreateBlindBoxContent
}

func (msg *MsgCreateBlindBoxContent) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgCreateBlindBoxContent) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCreateBlindBoxContent) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
