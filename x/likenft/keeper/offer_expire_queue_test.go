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

func getAllOfferExpireFromQueue(ctx sdk.Context, keeper *keeper.Keeper) []string {
	var results []string
	entries := keeper.GetOfferExpireQueue(ctx) // wrapper that uses iterator
	for _, entry := range entries {
		results = append(results, fmt.Sprintf("%s:%s", entry.ExpireTime, entry.OfferKey))
	}
	return results
}

func getPendingOfferExpireFromQueue(ctx sdk.Context, keeper *keeper.Keeper, endTime time.Time) []string {
	var results []string
	entries := keeper.GetOfferExpireQueueByTime(ctx, endTime) // wrapper that uses iterator
	for _, entry := range entries {
		results = append(results, fmt.Sprintf("%s:%s", entry.ExpireTime, entry.OfferKey))
	}
	return results
}

func TestOfferExpireQueueInsert(t *testing.T) {
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

	keeper.SetOfferExpireQueueEntry(ctx, types.OfferExpireQueueEntry{
		ExpireTime: time1,
		OfferKey:   key1,
	})
	keeper.SetOfferExpireQueueEntry(ctx, types.OfferExpireQueueEntry{
		ExpireTime: time2,
		OfferKey:   key2,
	})
	keeper.SetOfferExpireQueueEntry(ctx, types.OfferExpireQueueEntry{
		ExpireTime: time3,
		OfferKey:   key3,
	})
	keeper.SetOfferExpireQueueEntry(ctx, types.OfferExpireQueueEntry{
		ExpireTime: time4,
		OfferKey:   key4,
	})

	results := getAllOfferExpireFromQueue(ctx, keeper)
	require.Equal(t, 4, len(results))
}

func TestOfferExpireQueueRemove(t *testing.T) {
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

	keeper.SetOfferExpireQueueEntry(ctx, types.OfferExpireQueueEntry{
		ExpireTime: time1,
		OfferKey:   key1,
	})
	keeper.SetOfferExpireQueueEntry(ctx, types.OfferExpireQueueEntry{
		ExpireTime: time2,
		OfferKey:   key2,
	})
	keeper.SetOfferExpireQueueEntry(ctx, types.OfferExpireQueueEntry{
		ExpireTime: time3,
		OfferKey:   key3,
	})
	keeper.SetOfferExpireQueueEntry(ctx, types.OfferExpireQueueEntry{
		ExpireTime: time4,
		OfferKey:   key4,
	})

	// Remove valid entry
	keeper.RemoveOfferExpireQueueEntry(ctx, time2, key2)
	require.NoError(t, err)
	validRes := getAllOfferExpireFromQueue(ctx, keeper)
	require.Equal(t, 3, len(validRes))
}

func TestOfferExpireQueueUpdate(t *testing.T) {
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

	keeper.SetOfferExpireQueueEntry(ctx, types.OfferExpireQueueEntry{
		ExpireTime: time1,
		OfferKey:   key1,
	})
	keeper.SetOfferExpireQueueEntry(ctx, types.OfferExpireQueueEntry{
		ExpireTime: time2,
		OfferKey:   key2,
	})
	keeper.SetOfferExpireQueueEntry(ctx, types.OfferExpireQueueEntry{
		ExpireTime: time3,
		OfferKey:   key3,
	})
	keeper.SetOfferExpireQueueEntry(ctx, types.OfferExpireQueueEntry{
		ExpireTime: time4,
		OfferKey:   key4,
	})

	// Update entry
	updatedTime, err := time.Parse(time.RFC3339, "2099-01-01T00:00:00Z")
	require.NoError(t, err)

	keeper.UpdateOfferExpireQueueEntry(ctx, time4, key4, updatedTime)

	res := getAllOfferExpireFromQueue(ctx, keeper)
	require.Equal(t, []string{
		fmt.Sprintf("%s:%s", time1, key1),
		fmt.Sprintf("%s:%s", time2, key2),
		fmt.Sprintf("%s:%s", time3, key3),
		fmt.Sprintf("%s:%s", updatedTime, key4),
	}, res)
}

func TestOfferExpireQueueSorting(t *testing.T) {
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

	keeper.SetOfferExpireQueueEntry(ctx, types.OfferExpireQueueEntry{
		ExpireTime: time1,
		OfferKey:   key1,
	})
	keeper.SetOfferExpireQueueEntry(ctx, types.OfferExpireQueueEntry{
		ExpireTime: time2,
		OfferKey:   key2,
	})
	keeper.SetOfferExpireQueueEntry(ctx, types.OfferExpireQueueEntry{
		ExpireTime: time3,
		OfferKey:   key3,
	})
	keeper.SetOfferExpireQueueEntry(ctx, types.OfferExpireQueueEntry{
		ExpireTime: time4,
		OfferKey:   key4,
	})

	res := getAllOfferExpireFromQueue(ctx, keeper)
	require.Equal(t, 4, len(res))

	require.Equal(t, []string{
		fmt.Sprintf("%s:%s", time4, key4),
		fmt.Sprintf("%s:%s", time1, key1),
		fmt.Sprintf("%s:%s", time2, key2),
		fmt.Sprintf("%s:%s", time3, key3),
	}, res)
}

func TestOfferExpireQueueByTimeIterator(t *testing.T) {
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

	keeper.SetOfferExpireQueueEntry(ctx, types.OfferExpireQueueEntry{
		ExpireTime: time1,
		OfferKey:   key1,
	})
	keeper.SetOfferExpireQueueEntry(ctx, types.OfferExpireQueueEntry{
		ExpireTime: time2,
		OfferKey:   key2,
	})
	keeper.SetOfferExpireQueueEntry(ctx, types.OfferExpireQueueEntry{
		ExpireTime: time3,
		OfferKey:   key3,
	})
	keeper.SetOfferExpireQueueEntry(ctx, types.OfferExpireQueueEntry{
		ExpireTime: time4,
		OfferKey:   key4,
	})

	res := getPendingOfferExpireFromQueue(ctx, keeper, time2) // up to 2023-01-01T00:00:00Z exclusive
	require.Equal(t, 2, len(res))

	require.Equal(t, []string{
		fmt.Sprintf("%s:%s", time4, key4),
		fmt.Sprintf("%s:%s", time1, key1),
	}, res)
}
