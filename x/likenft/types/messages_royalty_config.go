package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgCreateRoyaltyConfig = "create_royalty_config"
	TypeMsgUpdateRoyaltyConfig = "update_royalty_config"
	TypeMsgDeleteRoyaltyConfig = "delete_royalty_config"
)

var _ sdk.Msg = &MsgCreateRoyaltyConfig{}

func NewMsgCreateRoyaltyConfig(
	creator string,
	classId string,
	royaltyConfig RoyaltyConfigInput,

) *MsgCreateRoyaltyConfig {
	return &MsgCreateRoyaltyConfig{
		Creator:       creator,
		ClassId:       classId,
		RoyaltyConfig: royaltyConfig,
	}
}

func (msg *MsgCreateRoyaltyConfig) Route() string {
	return RouterKey
}

func (msg *MsgCreateRoyaltyConfig) Type() string {
	return TypeMsgCreateRoyaltyConfig
}

func (msg *MsgCreateRoyaltyConfig) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgCreateRoyaltyConfig) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCreateRoyaltyConfig) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

var _ sdk.Msg = &MsgUpdateRoyaltyConfig{}

func NewMsgUpdateRoyaltyConfig(
	creator string,
	classId string,
	royaltyConfig RoyaltyConfigInput,

) *MsgUpdateRoyaltyConfig {
	return &MsgUpdateRoyaltyConfig{
		Creator:       creator,
		ClassId:       classId,
		RoyaltyConfig: royaltyConfig,
	}
}

func (msg *MsgUpdateRoyaltyConfig) Route() string {
	return RouterKey
}

func (msg *MsgUpdateRoyaltyConfig) Type() string {
	return TypeMsgUpdateRoyaltyConfig
}

func (msg *MsgUpdateRoyaltyConfig) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdateRoyaltyConfig) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateRoyaltyConfig) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

var _ sdk.Msg = &MsgDeleteRoyaltyConfig{}

func NewMsgDeleteRoyaltyConfig(
	creator string,
	classId string,

) *MsgDeleteRoyaltyConfig {
	return &MsgDeleteRoyaltyConfig{
		Creator: creator,
		ClassId: classId,
	}
}
func (msg *MsgDeleteRoyaltyConfig) Route() string {
	return RouterKey
}

func (msg *MsgDeleteRoyaltyConfig) Type() string {
	return TypeMsgDeleteRoyaltyConfig
}

func (msg *MsgDeleteRoyaltyConfig) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgDeleteRoyaltyConfig) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgDeleteRoyaltyConfig) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
