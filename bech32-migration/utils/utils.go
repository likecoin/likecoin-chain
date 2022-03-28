package utils

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func ConvertValAddr(valAddr string) string {
	parsedValAddr, err := sdk.ValAddressFromBech32(valAddr)
	if err != nil {
		return valAddr
	}
	return parsedValAddr.String()
}

func ConvertAccAddr(accAddr string) string {
	parsedAccAddr, err := sdk.AccAddressFromBech32(accAddr)
	if err != nil {
		return accAddr
	}
	return parsedAccAddr.String()
}

func ConvertConsAddr(consAddr string) string {
	parsedConsAddr, err := sdk.ConsAddressFromBech32(consAddr)
	if err != nil {
		return consAddr
	}
	return parsedConsAddr.String()
}

func IterateStoreByPrefix(
	ctx sdk.Context, storeKey sdk.StoreKey, prefix []byte,
	fn func(value []byte) (newValue []byte),
) {
	store := ctx.KVStore(storeKey)
	iterator := sdk.KVStorePrefixIterator(store, prefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		newValue := fn(iterator.Value())
		store.Set(iterator.Key(), newValue)
	}
}
