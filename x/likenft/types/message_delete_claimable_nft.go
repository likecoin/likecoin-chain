package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgDeleteClaimableNFT = "delete_claimable_nft"

var _ sdk.Msg = &MsgDeleteClaimableNFT{}

func NewMsgDeleteClaimableNFT(creator string, classId string, id string) *MsgDeleteClaimableNFT {
	return &MsgDeleteClaimableNFT{
		Creator: creator,
		ClassId: classId,
		Id:      id,
	}
}

func (msg *MsgDeleteClaimableNFT) Route() string {
	return RouterKey
}

func (msg *MsgDeleteClaimableNFT) Type() string {
	return TypeMsgDeleteClaimableNFT
}

func (msg *MsgDeleteClaimableNFT) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgDeleteClaimableNFT) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgDeleteClaimableNFT) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
