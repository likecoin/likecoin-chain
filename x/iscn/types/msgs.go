package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgCreateIscnRecord          = "create_iscn_record"
	TypeMsgUpdateIscnRecord          = "update_iscn_record"
	TypeMsgChangeIscnRecordOwnership = "msg_change_iscn_record_ownership"
)

var _ sdk.Msg = &MsgCreateIscnRecord{}
var _ sdk.Msg = &MsgUpdateIscnRecord{}
var _ sdk.Msg = &MsgChangeIscnRecordOwnership{}

func NewMsgCreateIscnRecord(from sdk.AccAddress, record *IscnRecord) *MsgCreateIscnRecord {
	return &MsgCreateIscnRecord{
		From:   from.String(),
		Record: *record,
	}
}

func (m MsgCreateIscnRecord) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&m)
	return sdk.MustSortJSON(bz)
}

func (m MsgCreateIscnRecord) GetSigners() []sdk.AccAddress {
	from, _ := sdk.AccAddressFromBech32(m.From)
	return []sdk.AccAddress{from}
}

func (msg MsgCreateIscnRecord) Route() string { return RouterKey }

func (msg MsgCreateIscnRecord) Type() string { return TypeMsgCreateIscnRecord }

func (msg MsgCreateIscnRecord) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err.Error())
	}
	err = msg.Record.Validate()
	if err != nil {
		return sdkerrors.Wrapf(ErrInvalidIscnRecord, "%s", err.Error())
	}
	return nil
}

func NewMsgUpdateIscnRecord(from sdk.AccAddress, iscnId IscnId, record *IscnRecord) *MsgUpdateIscnRecord {
	return &MsgUpdateIscnRecord{
		From:   from.String(),
		IscnId: iscnId.String(),
		Record: *record,
	}
}

func (m MsgUpdateIscnRecord) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&m)
	return sdk.MustSortJSON(bz)
}

func (m MsgUpdateIscnRecord) GetSigners() []sdk.AccAddress {
	from, _ := sdk.AccAddressFromBech32(m.From)
	return []sdk.AccAddress{from}
}

func (msg MsgUpdateIscnRecord) Route() string { return RouterKey }

func (msg MsgUpdateIscnRecord) Type() string { return TypeMsgUpdateIscnRecord }

func (msg MsgUpdateIscnRecord) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address: %s", err)
	}
	id, err := ParseIscnID(msg.IscnId)
	if err != nil {
		return sdkerrors.Wrapf(ErrInvalidIscnId, "%s", err.Error())
	}
	if id.Version == 0 {
		return sdkerrors.Wrapf(ErrInvalidIscnId, "invalid ISCN ID version")
	}
	err = msg.Record.Validate()
	if err != nil {
		return sdkerrors.Wrapf(ErrInvalidIscnRecord, "%s", err.Error())
	}
	return nil
}

func NewMsgChangeIscnRecordOwnership(from sdk.AccAddress, iscnId IscnId, newOwner sdk.AccAddress) *MsgChangeIscnRecordOwnership {
	return &MsgChangeIscnRecordOwnership{
		From:     from.String(),
		IscnId:   iscnId.String(),
		NewOwner: newOwner.String(),
	}
}

func (m MsgChangeIscnRecordOwnership) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&m)
	return sdk.MustSortJSON(bz)
}

func (m MsgChangeIscnRecordOwnership) GetSigners() []sdk.AccAddress {
	from, _ := sdk.AccAddressFromBech32(m.From)
	return []sdk.AccAddress{from}
}

func (msg MsgChangeIscnRecordOwnership) Route() string { return RouterKey }

func (msg MsgChangeIscnRecordOwnership) Type() string { return TypeMsgChangeIscnRecordOwnership }

func (msg MsgChangeIscnRecordOwnership) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address: %s", err.Error())
	}
	_, err = sdk.AccAddressFromBech32(msg.NewOwner)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid new owner address: %s", err.Error())
	}
	_, err = ParseIscnID(msg.IscnId)
	if err != nil {
		return sdkerrors.Wrapf(ErrInvalidIscnId, "%s", err.Error())
	}
	return nil
}
