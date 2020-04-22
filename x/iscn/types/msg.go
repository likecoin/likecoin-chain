package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ sdk.Msg = &MsgCreateIscn{}
var _ sdk.Msg = &MsgAddAuthor{}

type MsgAddAuthor struct {
	From       sdk.AccAddress `json:"from" yaml:"from"`
	AuthorInfo Author         `json:"authorInfo" yaml:"authorInfo"`
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

type MsgCreateIscn struct {
	From       sdk.AccAddress `json:"from" yaml:"from"`
	IscnRecord IscnRecord     `json:"iscnRecord" yaml:"iscnRecord"`
}

func NewMsgAddAuthor(from sdk.AccAddress, authorInfo Author) MsgAddAuthor {
	return MsgAddAuthor{
		From:       from,
		AuthorInfo: authorInfo,
	}
}

func (msg MsgAddAuthor) Route() string { return RouterKey }
func (msg MsgAddAuthor) Type() string  { return "add_author" }

func (msg MsgAddAuthor) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.From}
}

func (msg MsgAddAuthor) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgAddAuthor) ValidateBasic() sdk.Error {
	if msg.From.Empty() {
		return ErrInvalidApprover(DefaultCodespace)
	}
	// TODO: validate IscnRecord
	return nil
}
