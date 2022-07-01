package keeper_test

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	keepertest "github.com/likecoin/likecoin-chain/v3/testutil/keeper"
	"github.com/likecoin/likecoin-chain/v3/x/likenft/keeper"
	"github.com/likecoin/likecoin-chain/v3/x/likenft/types"
	"github.com/stretchr/testify/require"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func getAllListingExpireFromQueue(ctx sdk.Context, keeper *keeper.Keeper) []string {
	var results []string
	entries := keeper.GetListingExpireQueue(ctx) // wrapper that uses iterator
	for _, entry := range entries {
		results = append(results, fmt.Sprintf("%s:%s", entry.ExpireTime, entry.ListingKey))
	}
	return results
}

func getPendingListingExpireFromQueue(ctx sdk.Context, keeper *keeper.Keeper, endTime time.Time) []string {
	var results []string
	entries := keeper.GetListingExpireQueueByTime(ctx, endTime) // wrapper that uses iterator
	for _, entry := range entries {
		results = append(results, fmt.Sprintf("%s:%s", entry.ExpireTime, entry.ListingKey))
	}
	return results
}

func TestListingExpireQueueInsert(t *testing.T) {
	keeper, ctx := keepertest.LikenftKeeper(t)

	key1 := []byte("1")
	time1, err := time.Parse(time.RFC3339, "2020-01-01T00:00:00Z")
	require.NoError(t, err)

	key2 := []byte("2")
	time2, err := time.Parse(time.RFC3339, "2023-01-01T00:00:00Z")
	require.NoError(t, err)

	key3 := []byte("3")
	time3, err := time.Parse(time.RFC3339, "2027-01-01T00:00:00Z")
	require.NoError(t, err)

	key4 := []byte("4")
	time4, err := time.Parse(time.RFC3339, "2009-01-01T00:00:00Z")
	require.NoError(t, err)

	keeper.SetListingExpireQueueEntry(ctx, types.ListingExpireQueueEntry{
		ExpireTime: time1,
		ListingKey: key1,
	})
	keeper.SetListingExpireQueueEntry(ctx, types.ListingExpireQueueEntry{
		ExpireTime: time2,
		ListingKey: key2,
	})
	keeper.SetListingExpireQueueEntry(ctx, types.ListingExpireQueueEntry{
		ExpireTime: time3,
		ListingKey: key3,
	})
	keeper.SetListingExpireQueueEntry(ctx, types.ListingExpireQueueEntry{
		ExpireTime: time4,
		ListingKey: key4,
	})

	results := getAllListingExpireFromQueue(ctx, keeper)
	require.Equal(t, 4, len(results))
}

func TestListingExpireQueueRemove(t *testing.T) {
	keeper, ctx := keepertest.LikenftKeeper(t)

	key1 := []byte("1")
	time1, err := time.Parse(time.RFC3339, "2020-01-01T00:00:00Z")
	require.NoError(t, err)

	key2 := []byte("2")
	time2, err := time.Parse(time.RFC3339, "2023-01-01T00:00:00Z")
	require.NoError(t, err)

	key3 := []byte("3")
	time3, err := time.Parse(time.RFC3339, "2027-01-01T00:00:00Z")
	require.NoError(t, err)

	key4 := []byte("4")
	time4, err := time.Parse(time.RFC3339, "2009-01-01T00:00:00Z")
	require.NoError(t, err)

	keeper.SetListingExpireQueueEntry(ctx, types.ListingExpireQueueEntry{
		ExpireTime: time1,
		ListingKey: key1,
	})
	keeper.SetListingExpireQueueEntry(ctx, types.ListingExpireQueueEntry{
		ExpireTime: time2,
		ListingKey: key2,
	})
	keeper.SetListingExpireQueueEntry(ctx, types.ListingExpireQueueEntry{
		ExpireTime: time3,
		ListingKey: key3,
	})
	keeper.SetListingExpireQueueEntry(ctx, types.ListingExpireQueueEntry{
		ExpireTime: time4,
		ListingKey: key4,
	})

	// Remove valid entry
	keeper.RemoveListingExpireQueueEntry(ctx, time2, key2)
	require.NoError(t, err)
	validRes := getAllListingExpireFromQueue(ctx, keeper)
	require.Equal(t, 3, len(validRes))
}

func TestListingExpireQueueUpdate(t *testing.T) {
	keeper, ctx := keepertest.LikenftKeeper(t)

	key1 := []byte("1")
	time1, err := time.Parse(time.RFC3339, "2020-01-01T00:00:00Z")
	require.NoError(t, err)

	key2 := []byte("2")
	time2, err := time.Parse(time.RFC3339, "2023-01-01T00:00:00Z")
	require.NoError(t, err)

	key3 := []byte("3")
	time3, err := time.Parse(time.RFC3339, "2027-01-01T00:00:00Z")
	require.NoError(t, err)

	key4 := []byte("4")
	time4, err := time.Parse(time.RFC3339, "2009-01-01T00:00:00Z")
	require.NoError(t, err)

	keeper.SetListingExpireQueueEntry(ctx, types.ListingExpireQueueEntry{
		ExpireTime: time1,
		ListingKey: key1,
	})
	keeper.SetListingExpireQueueEntry(ctx, types.ListingExpireQueueEntry{
		ExpireTime: time2,
		ListingKey: key2,
	})
	keeper.SetListingExpireQueueEntry(ctx, types.ListingExpireQueueEntry{
		ExpireTime: time3,
		ListingKey: key3,
	})
	keeper.SetListingExpireQueueEntry(ctx, types.ListingExpireQueueEntry{
		ExpireTime: time4,
		ListingKey: key4,
	})

	// Update entry
	updatedTime, err := time.Parse(time.RFC3339, "2099-01-01T00:00:00Z")
	require.NoError(t, err)

	keeper.UpdateListingExpireQueueEntry(ctx, time4, key4, updatedTime)

	res := getAllListingExpireFromQueue(ctx, keeper)
	require.Equal(t, []string{
		fmt.Sprintf("%s:%s", time1, key1),
		fmt.Sprintf("%s:%s", time2, key2),
		fmt.Sprintf("%s:%s", time3, key3),
		fmt.Sprintf("%s:%s", updatedTime, key4),
	}, res)
}

func TestListingExpireQueueSorting(t *testing.T) {
	keeper, ctx := keepertest.LikenftKeeper(t)

	key1 := []byte("1")
	time1, err := time.Parse(time.RFC3339, "2020-01-01T00:00:00Z")
	require.NoError(t, err)

	key2 := []byte("2")
	time2, err := time.Parse(time.RFC3339, "2023-01-01T00:00:00Z")
	require.NoError(t, err)

	key3 := []byte("3")
	time3, err := time.Parse(time.RFC3339, "2027-01-01T00:00:00Z")
	require.NoError(t, err)

	key4 := []byte("4")
	time4, err := time.Parse(time.RFC3339, "2009-01-01T00:00:00Z")
	require.NoError(t, err)

	keeper.SetListingExpireQueueEntry(ctx, types.ListingExpireQueueEntry{
		ExpireTime: time1,
		ListingKey: key1,
	})
	keeper.SetListingExpireQueueEntry(ctx, types.ListingExpireQueueEntry{
		ExpireTime: time2,
		ListingKey: key2,
	})
	keeper.SetListingExpireQueueEntry(ctx, types.ListingExpireQueueEntry{
		ExpireTime: time3,
		ListingKey: key3,
	})
	keeper.SetListingExpireQueueEntry(ctx, types.ListingExpireQueueEntry{
		ExpireTime: time4,
		ListingKey: key4,
	})

	res := getAllListingExpireFromQueue(ctx, keeper)
	require.Equal(t, 4, len(res))

	require.Equal(t, []string{
		fmt.Sprintf("%s:%s", time4, key4),
		fmt.Sprintf("%s:%s", time1, key1),
		fmt.Sprintf("%s:%s", time2, key2),
		fmt.Sprintf("%s:%s", time3, key3),
	}, res)
}

func TestListingExpireQueueByTimeIterator(t *testing.T) {
	keeper, ctx := keepertest.LikenftKeeper(t)

	key1 := []byte("1")
	time1, err := time.Parse(time.RFC3339, "2020-01-01T00:00:00Z")
	require.NoError(t, err)

	key2 := []byte("2")
	time2, err := time.Parse(time.RFC3339, "2023-01-01T00:00:00Z")
	require.NoError(t, err)

	key3 := []byte("3")
	time3, err := time.Parse(time.RFC3339, "2027-01-01T00:00:00Z")
	require.NoError(t, err)

	key4 := []byte("4")
	time4, err := time.Parse(time.RFC3339, "2009-01-01T00:00:00Z")
	require.NoError(t, err)

	keeper.SetListingExpireQueueEntry(ctx, types.ListingExpireQueueEntry{
		ExpireTime: time1,
		ListingKey: key1,
	})
	keeper.SetListingExpireQueueEntry(ctx, types.ListingExpireQueueEntry{
		ExpireTime: time2,
		ListingKey: key2,
	})
	keeper.SetListingExpireQueueEntry(ctx, types.ListingExpireQueueEntry{
		ExpireTime: time3,
		ListingKey: key3,
	})
	keeper.SetListingExpireQueueEntry(ctx, types.ListingExpireQueueEntry{
		ExpireTime: time4,
		ListingKey: key4,
	})

	res := getPendingListingExpireFromQueue(ctx, keeper, time2) // up to 2023-01-01T00:00:00Z exclusive
	require.Equal(t, 2, len(res))

	require.Equal(t, []string{
		fmt.Sprintf("%s:%s", time4, key4),
		fmt.Sprintf("%s:%s", time1, key1),
	}, res)
}
