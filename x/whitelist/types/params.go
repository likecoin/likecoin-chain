package types

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

type Params struct {
	Approver sdk.AccAddress `json:"approver" yaml:"approver"`
}

var (
	KeyApprover = []byte("Approver")
)

var _ params.ParamSet = (*Params)(nil)

// Implements params.ParamSet
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		{KeyApprover, &p.Approver},
	}
}

func DefaultParams() Params {
	return Params{}
}

func (p Params) String() string {
	return fmt.Sprintf(`Params:
  Whitelist Approver: %s`, p.Approver)
}

func MustUnmarshalParams(cdc *codec.Codec, value []byte) Params {
	params, err := UnmarshalParams(cdc, value)
	if err != nil {
		panic(err)
	}
	return params
}

func UnmarshalParams(cdc *codec.Codec, value []byte) (params Params, err error) {
	err = cdc.UnmarshalBinaryLengthPrefixed(value, &params)
	if err != nil {
		return
	}
	return
}
