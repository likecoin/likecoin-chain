package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgCreateMintableNFT = "create_mintable_nft"

var _ sdk.Msg = &MsgCreateMintableNFT{}

func NewMsgCreateMintableNFT(creator string, classId string, id string, input NFTInput) *MsgCreateMintableNFT {
	return &MsgCreateMintableNFT{
		Creator: creator,
		ClassId: classId,
		Id:      id,
		Input:   input,
	}
}

func (msg *MsgCreateMintableNFT) Route() string {
	return RouterKey
}

func (msg *MsgCreateMintableNFT) Type() string {
	return TypeMsgCreateMintableNFT
}

func (msg *MsgCreateMintableNFT) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgCreateMintableNFT) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCreateMintableNFT) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
