package keeper_test

import (
	"testing"
	"time"

	"github.com/likecoin/likechain/x/likenft/keeper"
	"github.com/likecoin/likechain/x/likenft/testutil"
	"github.com/likecoin/likechain/x/likenft/types"
	"github.com/stretchr/testify/require"
)

func TestClaimPeriodSorter(t *testing.T) {
	// Setup

	claimPeriod1 := types.ClaimPeriod{
		StartTime:        testutil.MustParseTime(time.RFC3339, "2019-01-01T00:00:00Z"),
		MintPrice:        uint64(0),
		AllowedAddresses: []string{},
	}
	claimPeriod2 := types.ClaimPeriod{
		StartTime:        testutil.MustParseTime(time.RFC3339, "2020-01-01T00:00:00Z"),
		MintPrice:        uint64(0),
		AllowedAddresses: []string{},
	}
	claimPeriod3 := types.ClaimPeriod{
		StartTime:        testutil.MustParseTime(time.RFC3339, "2022-01-01T00:00:00Z"),
		MintPrice:        uint64(0),
		AllowedAddresses: []string{},
	}

	inputClaimPeriods := []types.ClaimPeriod{
		claimPeriod2,
		claimPeriod1,
		claimPeriod3,
	}

	ascendingRes := keeper.SortClaimPeriod(inputClaimPeriods, false)
	expectedAscendingClaimPeriods := []types.ClaimPeriod{
		claimPeriod1,
		claimPeriod2,
		claimPeriod3,
	}

	require.Equal(t, ascendingRes, expectedAscendingClaimPeriods)

	descendingRes := keeper.SortClaimPeriod(inputClaimPeriods, true)
	expectedDescendingClaimPeriods := []types.ClaimPeriod{
		claimPeriod3,
		claimPeriod2,
		claimPeriod1,
	}

	require.Equal(t, descendingRes, expectedDescendingClaimPeriods)
}
