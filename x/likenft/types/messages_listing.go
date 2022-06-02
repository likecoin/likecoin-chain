package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgCreateListing = "create_listing"
	TypeMsgUpdateListing = "update_listing"
	TypeMsgDeleteListing = "delete_listing"
)

var _ sdk.Msg = &MsgCreateListing{}

func NewMsgCreateListing(
	creator string,
	classId string,
	nftId string,
	price uint64,
	expiration time.Time,

) *MsgCreateListing {
	return &MsgCreateListing{
		Creator:    creator,
		ClassId:    classId,
		NftId:      nftId,
		Price:      price,
		Expiration: expiration,
	}
}

func (msg *MsgCreateListing) Route() string {
	return RouterKey
}

func (msg *MsgCreateListing) Type() string {
	return TypeMsgCreateListing
}

func (msg *MsgCreateListing) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgCreateListing) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCreateListing) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

var _ sdk.Msg = &MsgUpdateListing{}

func NewMsgUpdateListing(
	creator string,
	classId string,
	nftId string,
	price uint64,
	expiration time.Time,

) *MsgUpdateListing {
	return &MsgUpdateListing{
		Creator:    creator,
		ClassId:    classId,
		NftId:      nftId,
		Price:      price,
		Expiration: expiration,
	}
}

func (msg *MsgUpdateListing) Route() string {
	return RouterKey
}

func (msg *MsgUpdateListing) Type() string {
	return TypeMsgUpdateListing
}

func (msg *MsgUpdateListing) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdateListing) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateListing) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

var _ sdk.Msg = &MsgDeleteListing{}

func NewMsgDeleteListing(
	creator string,
	classId string,
	nftId string,
) *MsgDeleteListing {
	return &MsgDeleteListing{
		Creator: creator,
		ClassId: classId,
		NftId:   nftId,
	}
}
func (msg *MsgDeleteListing) Route() string {
	return RouterKey
}

func (msg *MsgDeleteListing) Type() string {
	return TypeMsgDeleteListing
}

func (msg *MsgDeleteListing) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgDeleteListing) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgDeleteListing) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
