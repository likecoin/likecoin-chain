package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgMintNFT = "mint_nft"

var _ sdk.Msg = &MsgMintNFT{}

func NewMsgMintNFT(creator string, classId string, id string, input NFTInput) *MsgMintNFT {
	return &MsgMintNFT{
		Creator: creator,
		ClassId: classId,
		Id:      id,
		Input:   input,
	}
}

func (msg *MsgMintNFT) Route() string {
	return RouterKey
}

func (msg *MsgMintNFT) Type() string {
	return TypeMsgMintNFT
}

func (msg *MsgMintNFT) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgMintNFT) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgMintNFT) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
