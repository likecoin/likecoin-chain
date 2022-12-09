package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgSellNFT = "sell_nft"

var _ sdk.Msg = &MsgSellNFT{}

func NewMsgSellNFT(creator string, classId string, nftId string, buyer string, price uint64, fullPayToRoyalty bool) *MsgSellNFT {
	return &MsgSellNFT{
		Creator:          creator,
		ClassId:          classId,
		NftId:            nftId,
		Buyer:            buyer,
		Price:            price,
		FullPayToRoyalty: fullPayToRoyalty,
	}
}

func (msg *MsgSellNFT) Route() string {
	return RouterKey
}

func (msg *MsgSellNFT) Type() string {
	return TypeMsgSellNFT
}

func (msg *MsgSellNFT) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgSellNFT) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgSellNFT) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
