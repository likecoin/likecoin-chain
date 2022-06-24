package types

import sdk "github.com/cosmos/cosmos-sdk/types"

type RoyaltyAllocation struct {
	Account sdk.AccAddress
	Amount  uint64
}
