package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (record ContentIdRecord) OwnerAddress() sdk.AccAddress {
	return sdk.AccAddress(record.OwnerAddressBytes)
}
