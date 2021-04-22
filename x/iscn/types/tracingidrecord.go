package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (record TracingIdRecord) OwnerAddress() sdk.AccAddress {
	return sdk.AccAddress(record.OwnerAddressBytes)
}
