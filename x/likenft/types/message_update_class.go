package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgUpdateClass = "update_class"

var _ sdk.Msg = &MsgUpdateClass{}

func NewMsgUpdateClass(creator string, classId string, name string, symbol string, description string, uri string, uriHash string, metadata JsonInput, burnable bool) *MsgUpdateClass {
	return &MsgUpdateClass{
		Creator:     creator,
		ClassId:     classId,
		Name:        name,
		Symbol:      symbol,
		Description: description,
		Uri:         uri,
		UriHash:     uriHash,
		Metadata:    metadata,
		Burnable:    burnable,
	}
}

func (msg *MsgUpdateClass) Route() string {
	return RouterKey
}

func (msg *MsgUpdateClass) Type() string {
	return TypeMsgUpdateClass
}

func (msg *MsgUpdateClass) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdateClass) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateClass) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
