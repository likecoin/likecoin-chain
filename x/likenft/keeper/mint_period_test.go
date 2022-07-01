package keeper_test

import (
	"testing"
	"time"

	"github.com/likecoin/likecoin-chain/v3/x/likenft/keeper"
	"github.com/likecoin/likecoin-chain/v3/x/likenft/testutil"
	"github.com/likecoin/likecoin-chain/v3/x/likenft/types"
	"github.com/stretchr/testify/require"
)

func TestMintPeriodSorter(t *testing.T) {
	// Setup

	mintPeriod1 := types.MintPeriod{
		StartTime:        *testutil.MustParseTime(time.RFC3339, "2019-01-01T00:00:00Z"),
		MintPrice:        uint64(0),
		AllowedAddresses: []string{},
	}
	mintPeriod2 := types.MintPeriod{
		StartTime:        *testutil.MustParseTime(time.RFC3339, "2020-01-01T00:00:00Z"),
		MintPrice:        uint64(0),
		AllowedAddresses: []string{},
	}
	mintPeriod3 := types.MintPeriod{
		StartTime:        *testutil.MustParseTime(time.RFC3339, "2020-01-01T00:00:00Z"),
		MintPrice:        uint64(400),
		AllowedAddresses: []string{},
	}
	mintPeriod4 := types.MintPeriod{
		StartTime:        *testutil.MustParseTime(time.RFC3339, "2022-01-01T00:00:00Z"),
		MintPrice:        uint64(0),
		AllowedAddresses: []string{},
	}
	mintPeriod5 := types.MintPeriod{
		StartTime:        *testutil.MustParseTime(time.RFC3339, "2022-01-01T00:00:00Z"),
		MintPrice:        uint64(100),
		AllowedAddresses: []string{},
	}

	inputMintPeriods := []types.MintPeriod{
		mintPeriod2,
		mintPeriod1,
		mintPeriod3,
		mintPeriod5,
		mintPeriod4,
	}

	ascendingRes := keeper.SortMintPeriod(inputMintPeriods, false)
	expectedAscendingMintPeriods := []types.MintPeriod{
		mintPeriod1,
		mintPeriod2,
		mintPeriod3,
		mintPeriod4,
		mintPeriod5,
	}

	require.Equal(t, ascendingRes, expectedAscendingMintPeriods)

	descendingRes := keeper.SortMintPeriod(inputMintPeriods, true)
	expectedDescendingMintPeriods := []types.MintPeriod{
		mintPeriod5,
		mintPeriod4,
		mintPeriod3,
		mintPeriod2,
		mintPeriod1,
	}

	require.Equal(t, descendingRes, expectedDescendingMintPeriods)
}
