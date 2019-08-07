package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ sdk.Msg = &MsgSetWhitelist{}

type MsgSetWhitelist struct {
	Approver  sdk.AccAddress `json:"approver" yaml:"approver"`
	Whitelist Whitelist      `json:"whitelist" yaml:"whitelist"`
}

func NewMsgSetWhitelist(approver sdk.AccAddress, whitelist Whitelist) MsgSetWhitelist {
	return MsgSetWhitelist{
		Approver:  approver,
		Whitelist: whitelist,
	}
}

func (msg MsgSetWhitelist) Route() string { return RouterKey }
func (msg MsgSetWhitelist) Type() string  { return "set_whitelist" }

func (msg MsgSetWhitelist) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Approver}
}

func (msg MsgSetWhitelist) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgSetWhitelist) ValidateBasic() sdk.Error {
	if msg.Approver.Empty() {
		return ErrInvalidApprover(DefaultCodespace)
	}
	return nil
}
