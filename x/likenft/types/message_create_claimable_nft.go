package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgCreateClaimableNFT = "create_claimable_nft"

var _ sdk.Msg = &MsgCreateClaimableNFT{}

func NewMsgCreateClaimableNFT(creator string, classId string, id string, input NFTInput) *MsgCreateClaimableNFT {
	return &MsgCreateClaimableNFT{
		Creator: creator,
		ClassId: classId,
		Id:      id,
		Input:   input,
	}
}

func (msg *MsgCreateClaimableNFT) Route() string {
	return RouterKey
}

func (msg *MsgCreateClaimableNFT) Type() string {
	return TypeMsgCreateClaimableNFT
}

func (msg *MsgCreateClaimableNFT) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgCreateClaimableNFT) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCreateClaimableNFT) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
