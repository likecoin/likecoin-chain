package keeper

import (
	"fmt"
	"sort"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/likecoin/likechain/x/likenft/types"
)

func (k Keeper) resolveValidMintPeriod(ctx sdk.Context, classId string, classData types.ClassData, ownerAddress sdk.AccAddress, userAddress sdk.AccAddress) (*types.MintPeriod, error) {

	mintPeriods := classData.Config.MintPeriods
	if len(mintPeriods) == 0 {
		return nil, sdkerrors.ErrUnauthorized.Wrapf(fmt.Sprintf("No mint period is configured for class %s", classId))
	}

	for _, mintPeriod := range mintPeriods {
		// Check the first applicable mint period
		if mintPeriod.StartTime.Before(ctx.BlockHeader().Time) {
			// Check if the user is allowed to mint the token
			// If the minter is the owner, any mint period that is after the block time is valid
			// If allowed address list is nil it means the the class is publically available
			if ownerAddress.Equals(userAddress) || len(mintPeriod.AllowedAddresses) == 0 {
				return &mintPeriod, nil
			}

			for _, allowedAddress := range mintPeriod.AllowedAddresses {
				// Ensure the configured allowed address is valid
				wrappedAddress, err := sdk.AccAddressFromBech32(allowedAddress)
				if err != nil {
					return nil, types.ErrFailedToMintNFT.Wrapf(fmt.Sprintf("Failed to parse allowed address %s", allowedAddress))
				}

				if userAddress.Equals(wrappedAddress) {
					return &mintPeriod, nil
				}
			}
		}
	}

	return nil, nil
}

func SortMintPeriod(mintPeriods []types.MintPeriod, descending bool) []types.MintPeriod {
	// Sort the mint periods by start time
	sort.Slice(mintPeriods, func(i, j int) bool {
		if descending {
			i, j = j, i
		}

		if mintPeriods[j].StartTime.Equal(*mintPeriods[i].StartTime) {
			return mintPeriods[j].MintPrice > mintPeriods[i].MintPrice
		}

		return mintPeriods[j].StartTime.After(*mintPeriods[i].StartTime)
	})

	return mintPeriods
}
