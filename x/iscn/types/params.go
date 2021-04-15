package types

import (
	"fmt"
	"regexp"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

var (
	DefaultFeePerByteDenom        = "nanolike"
	DefaultFeePerByteAmount int64 = 10000
	DefaultRegistryId             = "likecoin-chain"
	DefaultFeePerByte             = sdk.NewDecCoin(
		DefaultFeePerByteDenom, sdk.NewInt(DefaultFeePerByteAmount),
	)
)

var (
	ParamKeyRegistryId = []byte("RegistryId")
	ParamKeyFeePerByte = []byte("FeePerByte")
)

func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

var _ paramtypes.ParamSet = (*Params)(nil)

func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(ParamKeyRegistryId, &p.RegistryId, validateRegistryId),
		paramtypes.NewParamSetPair(ParamKeyFeePerByte, &p.FeePerByte, validateFeePerByte),
	}
}

func validateRegistryId(i interface{}) error {
	s, ok := i.(string)
	if !ok {
		return fmt.Errorf("ISCN registry ID must be string, got type %T", i)
	}
	if s == "" {
		return fmt.Errorf("ISCN registry ID is empty")
	}
	// TODO: see if need to add any other characters to valid character list
	if regexp.MustCompile("^[-_.:a-zA-Z0-9]+$").FindString(s) != s {
		return fmt.Errorf("ISCN registry ID contains invalid character")
	}
	return nil
}

func validateFeePerByte(i interface{}) error {
	v, ok := i.(sdk.DecCoin)
	if !ok {
		return fmt.Errorf("ISCN fee per byte has invalid type: %T", i)
	}
	if v.Denom == "" {
		return fmt.Errorf("ISCN fee per byte has empty denom")
	}
	if v.IsNegative() {
		return fmt.Errorf("ISCN fee per byte must be non-negative, got %s", v.Amount.String())
	}
	return nil
}

func DefaultParams() Params {
	return Params{
		RegistryId: DefaultRegistryId,
		FeePerByte: DefaultFeePerByte,
	}
}

func (p Params) Validate() error {
	var err error
	err = validateRegistryId(p.RegistryId)
	if err != nil {
		return err
	}
	err = validateFeePerByte(p.FeePerByte)
	if err != nil {
		return err
	}
	return nil
}

func (p Params) String() string {
	return fmt.Sprintf(`Params:
  Registry ID: %s,
  Fee per byte: %s`,
		p.RegistryId,
		p.FeePerByte.String(),
	)
}
