package testutil

import (
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/x/auth/types"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"

	"github.com/likecoin/likechain/bech32-migration/utils"
)

func AssertAuthAddressBech32(ctx sdk.Context, storeKey sdk.StoreKey, cdc codec.BinaryCodec) bool {
	ok := true
	utils.IterateStoreByPrefix(ctx, storeKey, types.AddressStoreKeyPrefix, func(bz []byte) []byte {
		var accountI types.AccountI
		err := cdc.UnmarshalInterface(bz, &accountI)
		if err != nil {
			panic(err)
		}
		var acc string
		var itype string
		switch accountI.(type) {
		case *types.BaseAccount:
			acc = accountI.(*types.BaseAccount).Address
			itype = "BaseAccount"
		case *types.ModuleAccount:
			acc = accountI.(*types.ModuleAccount).Address
			itype = "ModuleAccount"
		case *vestingtypes.BaseVestingAccount:
			acc = accountI.(*vestingtypes.BaseVestingAccount).Address
			itype = "BaseVestingAccount"
		case *vestingtypes.ContinuousVestingAccount:
			acc = accountI.(*vestingtypes.ContinuousVestingAccount).Address
			itype = "ContinuousVestingAccount"
		case *vestingtypes.DelayedVestingAccount:
			acc = accountI.(*vestingtypes.DelayedVestingAccount).Address
			itype = "DelayedVestingAccount"
		case *vestingtypes.PeriodicVestingAccount:
			acc = accountI.(*vestingtypes.PeriodicVestingAccount).Address
			itype = "PeriodicVestingAccount"
		case *vestingtypes.PermanentLockedAccount:
			acc = accountI.(*vestingtypes.PermanentLockedAccount).Address
			itype = "PermanentLockedAccount"
		default:
			ctx.Logger().Info(
				"Warning: unknown account type, skipping migration",
				"address", accountI.GetAddress().String(),
				"account_number", accountI.GetAccountNumber(),
				"public_key", accountI.GetPubKey(),
				"sequence", accountI.GetSequence(),
			)
			return bz
		}
		if !strings.HasPrefix(acc, "like") {
			ctx.Logger().Info(fmt.Sprintf("Bad prefix found: %s, interface type: %s", acc, itype))
			ok = false
		}
		return bz
	})
	return ok
}
