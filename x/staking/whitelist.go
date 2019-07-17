package staking

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func checkWhitelistIncludeValidator(whitelist []sdk.ValAddress, valAddr sdk.ValAddress) bool {
	for _, v := range whitelist {
		if v.Equals(valAddr) {
			return true
		}
	}
	return false
}

func checkCanCreateValidator(whitelist []sdk.ValAddress, valAddr sdk.ValAddress) bool {
	return len(whitelist) == 0 || checkWhitelistIncludeValidator(whitelist, valAddr)
}
