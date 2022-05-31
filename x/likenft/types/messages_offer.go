package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgCreateOffer = "create_offer"
	TypeMsgUpdateOffer = "update_offer"
	TypeMsgDeleteOffer = "delete_offer"
)

var _ sdk.Msg = &MsgCreateOffer{}

func NewMsgCreateOffer(
	creator string,
	classId string,
	nftId string,
	price string,
	expiration string,

) *MsgCreateOffer {
	return &MsgCreateOffer{
		Creator:    creator,
		ClassId:    classId,
		NftId:      nftId,
		Price:      price,
		Expiration: expiration,
	}
}

func (msg *MsgCreateOffer) Route() string {
	return RouterKey
}

func (msg *MsgCreateOffer) Type() string {
	return TypeMsgCreateOffer
}

func (msg *MsgCreateOffer) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgCreateOffer) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCreateOffer) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

var _ sdk.Msg = &MsgUpdateOffer{}

func NewMsgUpdateOffer(
	creator string,
	classId string,
	nftId string,
	price string,
	expiration string,

) *MsgUpdateOffer {
	return &MsgUpdateOffer{
		Creator:    creator,
		ClassId:    classId,
		NftId:      nftId,
		Price:      price,
		Expiration: expiration,
	}
}

func (msg *MsgUpdateOffer) Route() string {
	return RouterKey
}

func (msg *MsgUpdateOffer) Type() string {
	return TypeMsgUpdateOffer
}

func (msg *MsgUpdateOffer) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdateOffer) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateOffer) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

var _ sdk.Msg = &MsgDeleteOffer{}

func NewMsgDeleteOffer(
	creator string,
	classId string,
	nftId string,

) *MsgDeleteOffer {
	return &MsgDeleteOffer{
		Creator: creator,
		ClassId: classId,
		NftId:   nftId,
	}
}
func (msg *MsgDeleteOffer) Route() string {
	return RouterKey
}

func (msg *MsgDeleteOffer) Type() string {
	return TypeMsgDeleteOffer
}

func (msg *MsgDeleteOffer) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgDeleteOffer) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgDeleteOffer) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
