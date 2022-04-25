package keeper

import (
	"fmt"
	"sort"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/likecoin/likechain/x/likenft/types"
)

func (k Keeper) resolveValidClaimPeriod(ctx sdk.Context, classId string, classData types.ClassData, userAddress sdk.AccAddress) (*types.ClaimPeriod, error) {

	claimPeriods := classData.Config.ClaimPeriods
	if claimPeriods == nil {
		return nil, types.ErrFailedToMintNFT.Wrapf(fmt.Sprintf("Pay to mint is not configurated for the class %s ", classId))
	}

	for _, claimPeriod := range claimPeriods {
		// Check the first applicable claim period
		if claimPeriod.StartTime.Before(ctx.BlockHeader().Time) {
			// Check if the user is allowed to mint the token
			// If allowed address list is nil it means the the class is publically available
			if claimPeriod.AllowedAddressList == nil {
				return claimPeriod, nil
			}
			for _, allowedAddress := range claimPeriod.AllowedAddressList {
				// Ensure the configured allowed address is valid
				wrappedAddress, err := sdk.AccAddressFromBech32(allowedAddress)
				if err != nil {
					return nil, types.ErrFailedToMintNFT.Wrapf(fmt.Sprintf("Failed to parse allowed address %s", allowedAddress))
				}

				if userAddress.Equals(wrappedAddress) {
					return claimPeriod, nil
				}
			}
		}
	}

	return nil, nil
}

func SortClaimPeriod(claimPeriods []*types.ClaimPeriod, descending bool) []*types.ClaimPeriod {
	// Sort the claim periods by start time
	sort.Slice(claimPeriods, func(i, j int) bool {
		if descending {
			i, j = j, i
		}
		return claimPeriods[j].StartTime.After(*claimPeriods[i].StartTime)
	})

	return claimPeriods
}
