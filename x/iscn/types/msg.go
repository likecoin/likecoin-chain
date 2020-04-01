package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ sdk.Msg = &MsgCreateIscn{}

type MsgCreateIscn struct {
	From       sdk.AccAddress `json:"from" yaml:"from"`
	IscnRecord IscnRecord     `json:"iscnRecord" yaml:"iscnRecord"`
}

func NewMsgCreateIscn(from sdk.AccAddress, iscnRecord IscnRecord) MsgCreateIscn {
	return MsgCreateIscn{
		From:       from,
		IscnRecord: iscnRecord,
	}
}

func (msg MsgCreateIscn) Route() string { return RouterKey }
func (msg MsgCreateIscn) Type() string  { return "create_iscn" }

func (msg MsgCreateIscn) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.From}
}

func (msg MsgCreateIscn) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgCreateIscn) ValidateBasic() sdk.Error {
	if msg.From.Empty() {
		return ErrInvalidApprover(DefaultCodespace)
	}
	// TODO: validate IscnRecord
	return nil
}
