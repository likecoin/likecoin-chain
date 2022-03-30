package auth

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/likecoin/likechain/bech32-migration/utils"
)

func MigrateAddressBech32(ctx sdk.Context, storeKey sdk.StoreKey, cdc codec.BinaryCodec) {
	ctx.Logger().Info("Migration of address bech32 for auth module begin")
	migratedAccountCount := uint64(0)
	utils.IterateStoreByPrefix(ctx, storeKey, types.AddressStoreKeyPrefix, func(bz []byte) []byte {
		var accountI types.AccountI
		err := cdc.UnmarshalInterface(bz, &accountI)
		if err != nil {
			panic(err)
		}
		switch accountI.(type) {
		case *types.BaseAccount:
			acc := accountI.(*types.BaseAccount)
			acc.Address = utils.ConvertAccAddr(acc.Address)
		case *types.ModuleAccount:
			acc := accountI.(*types.ModuleAccount)
			acc.Address = utils.ConvertAccAddr(acc.Address)
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
		bz, err = cdc.MarshalInterface(accountI)
		if err != nil {
			panic(err)
		}
		migratedAccountCount++
		return bz
	})
	ctx.Logger().Info(
		"Migration of address bech32 for auth module done",
		"migrated_account_count", migratedAccountCount,
	)
}
