package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

var (
	DefaultPriceDenom             = "nanolike"
	DefaultFeePerByteDenom        = "nanolike"
	DefaultFeePerByteAmount int64 = 10000
	DefaultFeePerByte             = sdk.NewDecCoin(
		DefaultFeePerByteDenom, sdk.NewInt(DefaultFeePerByteAmount),
	)
	DefaultMaxOfferDurationDays   uint64 = 180
	DefaultMaxListingDurationDays uint64 = 180
	DefaultMaxRoyaltyBasisPoints  uint64 = 1000 // 10%
)

var (
	ParamKeyPriceDenom             = []byte("PriceDenom")
	ParamKeyFeePerByte             = []byte("FeePerByte")
	ParamKeyMaxOfferDurationDays   = []byte("MaxOfferDurationDays")
	ParamKeyMaxListingDurationDays = []byte("MaxListingDurationDays")
	ParamKeyMaxRoyaltyBasisPoints  = []byte("MaxRoyaltyBasisPoints")
)

var _ paramtypes.ParamSet = (*Params)(nil)

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return Params{
		PriceDenom:             DefaultPriceDenom,
		FeePerByte:             DefaultFeePerByte,
		MaxOfferDurationDays:   DefaultMaxOfferDurationDays,
		MaxListingDurationDays: DefaultMaxListingDurationDays,
		MaxRoyaltyBasisPoints:  DefaultMaxRoyaltyBasisPoints,
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

// Validate fee per byte
func validateFeePerByte(i interface{}) error {
	v, ok := i.(sdk.DecCoin)
	if !ok {
		return fmt.Errorf("LikeNFT fee per byte has invalid type: %T", i)
	}
	if v.Denom == "" {
		return fmt.Errorf("LikeNFT fee per byte has empty denom")
	}
	if v.IsNegative() {
		return fmt.Errorf("LikeNFT fee per byte must be non-negative, got %s", v.Amount.String())
	}
	return nil
}

// Validate max offer duration
func validateMaxOfferDurationDays(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("Max offer duration has invalid type: %T", i)
	}
	if v == 0 {
		return fmt.Errorf("Max offer duration is zero")
	}
	return nil
}

// Validate max listing duration
func validateMaxListingDurationDays(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("Max listing duration has invalid type: %T", i)
	}
	if v == 0 {
		return fmt.Errorf("Max listing duration is zero")
	}
	return nil
}

// Validate max royalty basis points
func validateMaxRoyaltyBasisPoints(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("Max royalty basis points has invalid type: %T", i)
	}
	if v > 10000 {
		return fmt.Errorf("Max royalty basis points is larger than 10000 (100%%)")
	}
	return nil
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(ParamKeyPriceDenom, &p.PriceDenom, validatePriceDenom),
		paramtypes.NewParamSetPair(ParamKeyFeePerByte, &p.FeePerByte, validateFeePerByte),
		paramtypes.NewParamSetPair(ParamKeyMaxOfferDurationDays, &p.MaxOfferDurationDays, validateMaxOfferDurationDays),
		paramtypes.NewParamSetPair(ParamKeyMaxListingDurationDays, &p.MaxListingDurationDays, validateMaxListingDurationDays),
		paramtypes.NewParamSetPair(ParamKeyMaxRoyaltyBasisPoints, &p.MaxRoyaltyBasisPoints, validateMaxRoyaltyBasisPoints),
	}
}

// Validate validates the set of params
func (p Params) Validate() error {
	var err error
	err = validatePriceDenom(p.PriceDenom)
	if err != nil {
		return err
	}
	err = validateFeePerByte(p.FeePerByte)
	if err != nil {
		return err
	}
	err = validateMaxOfferDurationDays(p.MaxOfferDurationDays)
	if err != nil {
		return err
	}
	err = validateMaxListingDurationDays(p.MaxListingDurationDays)
	if err != nil {
		return err
	}
	err = validateMaxRoyaltyBasisPoints(p.MaxRoyaltyBasisPoints)
	if err != nil {
		return err
	}
	return nil
}

// String implements the Stringer interface.
func (p Params) String() string {
	return fmt.Sprintf(`Params:
	Price denom: %s,
	Fee per byte: %s
	Max offer duration days: %d
	Max listing duration days: %d
	Max royalty basis points: %d`,
		p.PriceDenom,
		p.FeePerByte,
		p.MaxOfferDurationDays,
		p.MaxListingDurationDays,
		p.MaxRoyaltyBasisPoints,
	)
}
