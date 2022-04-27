package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgDeleteMintableNFT = "delete_mintable_nft"

var _ sdk.Msg = &MsgDeleteMintableNFT{}

func NewMsgDeleteMintableNFT(creator string, classId string, id string) *MsgDeleteMintableNFT {
	return &MsgDeleteMintableNFT{
		Creator: creator,
		ClassId: classId,
		Id:      id,
	}
}

func (msg *MsgDeleteMintableNFT) Route() string {
	return RouterKey
}

func (msg *MsgDeleteMintableNFT) Type() string {
	return TypeMsgDeleteMintableNFT
}

func (msg *MsgDeleteMintableNFT) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgDeleteMintableNFT) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgDeleteMintableNFT) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
