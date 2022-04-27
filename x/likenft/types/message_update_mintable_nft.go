package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgUpdateMintableNFT = "update_mintable_nft"

var _ sdk.Msg = &MsgUpdateMintableNFT{}

func NewMsgUpdateMintableNFT(creator string, classId string, id string, input NFTInput) *MsgUpdateMintableNFT {
	return &MsgUpdateMintableNFT{
		Creator: creator,
		ClassId: classId,
		Id:      id,
		Input:   input,
	}
}

func (msg *MsgUpdateMintableNFT) Route() string {
	return RouterKey
}

func (msg *MsgUpdateMintableNFT) Type() string {
	return TypeMsgUpdateMintableNFT
}

func (msg *MsgUpdateMintableNFT) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdateMintableNFT) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateMintableNFT) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
