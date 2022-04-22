package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgUpdateClaimableNFT = "update_claimable_nft"

var _ sdk.Msg = &MsgUpdateClaimableNFT{}

func NewMsgUpdateClaimableNFT(creator string, classId string, id string, input string) *MsgUpdateClaimableNFT {
	return &MsgUpdateClaimableNFT{
		Creator: creator,
		ClassId: classId,
		Id:      id,
		Input:   input,
	}
}

func (msg *MsgUpdateClaimableNFT) Route() string {
	return RouterKey
}

func (msg *MsgUpdateClaimableNFT) Type() string {
	return TypeMsgUpdateClaimableNFT
}

func (msg *MsgUpdateClaimableNFT) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdateClaimableNFT) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateClaimableNFT) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
