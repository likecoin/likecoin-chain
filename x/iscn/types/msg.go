package types

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
	cbornode "github.com/ipfs/go-ipld-cbor"
)

var _ sdk.Msg = &MsgCreateIscn{}
var _ sdk.Msg = &MsgAddEntity{}

type IscnInput []byte // CBOR encoded

func (input IscnInput) MarshalJSON() ([]byte, error) {
	rawMap := RawIscnMap{}
	err := cbornode.DecodeInto(input, &rawMap)
	if err != nil {
		return nil, err
	}
	return json.Marshal(rawMap)
}

func (input *IscnInput) UnmarshalJSON(bz []byte) error {
	rawMap := RawIscnMap{}
	err := json.Unmarshal(bz, &rawMap)
	if err != nil {
		return err
	}
	bz, err = cbornode.DumpObject(rawMap)
	if err != nil {
		return err
	}
	*input = bz
	return nil
}

type MsgCreateIscn struct {
	From       sdk.AccAddress `json:"from" yaml:"from"`
	IscnKernel IscnInput      `json:"iscnKernel" yaml:"iscnKernel"`
}

func NewMsgCreateIscn(from sdk.AccAddress, iscnKernel IscnInput) MsgCreateIscn {
	return MsgCreateIscn{
		From:       from,
		IscnKernel: iscnKernel,
	}
}

func (msg MsgCreateIscn) Route() string { return RouterKey }
func (msg MsgCreateIscn) Type() string  { return "create_iscn" }

func (msg MsgCreateIscn) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.From}
}

func (msg MsgCreateIscn) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	// TODO: unmarshal from CBOR, then marshal as JSON
	bz = sdk.MustSortJSON(bz)
	return bz
}

func (msg MsgCreateIscn) ValidateBasic() sdk.Error {
	if msg.From.Empty() {
		return ErrInvalidSender(DefaultCodespace)
	}
	// TODO: validate IscnRecord
	// 1. timestamps
	// 2. if parent is empty, version should be 1
	return nil
}

type MsgAddEntity struct {
	From   sdk.AccAddress `json:"from" yaml:"from"`
	Entity IscnInput      `json:"entity" yaml:"entity"`
}

func NewMsgAddEntity(from sdk.AccAddress, entity IscnInput) MsgAddEntity {
	return MsgAddEntity{
		From:   from,
		Entity: entity,
	}
}

func (msg MsgAddEntity) Route() string { return RouterKey }
func (msg MsgAddEntity) Type() string  { return "add_entity" }

func (msg MsgAddEntity) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.From}
}

func (msg MsgAddEntity) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	// TODO: unmarshal from CBOR, then marshal as JSON
	return sdk.MustSortJSON(bz)
}

func (msg MsgAddEntity) ValidateBasic() sdk.Error {
	if msg.From.Empty() {
		return ErrInvalidSender(DefaultCodespace)
	}
	return nil
}

type MsgAddRightTerms struct {
	From       sdk.AccAddress `json:"from" yaml:"from"`
	RightTerms string         `json:"rightTerms" yaml:"rightTerms"`
}

func NewMsgAddRightTerms(from sdk.AccAddress, rightTerms string) MsgAddRightTerms {
	return MsgAddRightTerms{
		From:       from,
		RightTerms: rightTerms,
	}
}

func (msg MsgAddRightTerms) Route() string { return RouterKey }
func (msg MsgAddRightTerms) Type() string  { return "add_right_terms" }

func (msg MsgAddRightTerms) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.From}
}

func (msg MsgAddRightTerms) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgAddRightTerms) ValidateBasic() sdk.Error {
	if msg.From.Empty() {
		return ErrInvalidSender(DefaultCodespace)
	}
	return nil
}
