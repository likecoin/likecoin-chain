package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

var (
	DefaultMintPriceDenom = "nanolike"
)

var (
	ParamKeyMintPriceDenom = []byte("MintPriceDenom")
)

var _ paramtypes.ParamSet = (*Params)(nil)

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return Params{
		MintPriceDenom: DefaultMintPriceDenom,
	}
}

// Validate Mint price denom type
func validateMintPriceDenom(i interface{}) error {
	s, ok := i.(string)
	if !ok {
		return fmt.Errorf("Mint price denom must be string, got type %T", i)
	}
	if s == "" {
		return fmt.Errorf("Mint price denom is empty")
	}
	if err := sdk.ValidateDenom(s); err != nil {
		return err
	}

	return nil
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(ParamKeyMintPriceDenom, &p.MintPriceDenom, validateMintPriceDenom),
	}
}

// Validate validates the set of params
func (p Params) Validate() error {
	var err error
	err = validateMintPriceDenom(p.MintPriceDenom)
	if err != nil {
		return err
	}
	return nil
}

// String implements the Stringer interface.
func (p Params) String() string {
	return fmt.Sprintf(`Params:
	Mint Price denom: %s`, p.MintPriceDenom)
}
