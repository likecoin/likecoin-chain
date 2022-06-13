package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

var (
	DefaultPriceDenom = "nanolike"
)

var (
	ParamKeyPriceDenom = []byte("PriceDenom")
)

var _ paramtypes.ParamSet = (*Params)(nil)

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return Params{
		PriceDenom: DefaultPriceDenom,
	}
}

// Validate price denom type
func validatePriceDenom(i interface{}) error {
	s, ok := i.(string)
	if !ok {
		return fmt.Errorf("Price denom must be string, got type %T", i)
	}
	if s == "" {
		return fmt.Errorf("Price denom is empty")
	}
	if err := sdk.ValidateDenom(s); err != nil {
		return err
	}

	return nil
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(ParamKeyPriceDenom, &p.PriceDenom, validatePriceDenom),
	}
}

// Validate validates the set of params
func (p Params) Validate() error {
	var err error
	err = validatePriceDenom(p.PriceDenom)
	if err != nil {
		return err
	}
	return nil
}

// String implements the Stringer interface.
func (p Params) String() string {
	return fmt.Sprintf(`Params:
	Price denom: %s`, p.PriceDenom)
}
