package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgBuyNFT = "buy_nft"

var _ sdk.Msg = &MsgBuyNFT{}

func NewMsgBuyNFT(creator string, classId string, nftId string, seller string, price uint64) *MsgBuyNFT {
	return &MsgBuyNFT{
		Creator: creator,
		ClassId: classId,
		NftId:   nftId,
		Seller:  seller,
		Price:   price,
	}
}

func (msg *MsgBuyNFT) Route() string {
	return RouterKey
}

func (msg *MsgBuyNFT) Type() string {
	return TypeMsgBuyNFT
}

func (msg *MsgBuyNFT) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgBuyNFT) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgBuyNFT) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
