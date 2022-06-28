package testutil

import (
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/x/slashing/types"

	"github.com/likecoin/likecoin-chain/v3/bech32-migration/utils"
)

func AssertSlashingAddressBech32(ctx sdk.Context, storeKey sdk.StoreKey, cdc codec.BinaryCodec) bool {
	ok := true
	utils.IterateStoreByPrefix(ctx, storeKey, types.ValidatorSigningInfoKeyPrefix, func(bz []byte) []byte {
		validatorSigningInfo := types.ValidatorSigningInfo{}
		cdc.MustUnmarshal(bz, &validatorSigningInfo)
		if !strings.HasPrefix(validatorSigningInfo.Address, "like") {
			ctx.Logger().Info(fmt.Sprintf("Bad prefix found: %s, interface type: validatorSigningInfo.Address", validatorSigningInfo.Address))
			ok = false
		}
		return bz
	})
	return ok
}
