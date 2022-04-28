package keeper_test

import (
	"strconv"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	keepertest "github.com/likecoin/likechain/testutil/keeper"
	"github.com/likecoin/likechain/x/likenft/keeper"
	"github.com/likecoin/likechain/x/likenft/types"
	"github.com/stretchr/testify/require"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func getAllClassIdsFromQueue(ctx sdk.Context, keeper *keeper.Keeper) []string {
	var resultClassIds []string
	keeper.IterateClassRevealQueue(ctx, func(val types.ClassRevealQueueEntry) (stop bool) {
		resultClassIds = append(resultClassIds, val.ClassId)
		return false
	})
	return resultClassIds
}

func TestClassRevealQueueInsert(t *testing.T) {
	keeper, ctx := keepertest.LikenftKeeper(t)

	classId1 := "likenft10"
	revealTime1, err := time.Parse(time.RFC3339, "2020-01-01T00:00:00Z")
	require.NoError(t, err)

	classId2 := "likenft11"
	revealTime2, err := time.Parse(time.RFC3339, "2023-01-01T00:00:00Z")
	require.NoError(t, err)

	classId3 := "likenft12"
	revealTime3, err := time.Parse(time.RFC3339, "2027-01-01T00:00:00Z")
	require.NoError(t, err)

	classId4 := "likenft13"
	revealTime4, err := time.Parse(time.RFC3339, "2009-01-01T00:00:00Z")
	require.NoError(t, err)

	keeper.SetClassRevealQueueEntry(ctx, types.ClassRevealQueueEntry{
		RevealTime: revealTime1,
		ClassId:    classId1,
	})
	keeper.SetClassRevealQueueEntry(ctx, types.ClassRevealQueueEntry{
		RevealTime: revealTime2,
		ClassId:    classId2,
	})
	keeper.SetClassRevealQueueEntry(ctx, types.ClassRevealQueueEntry{
		RevealTime: revealTime3,
		ClassId:    classId3,
	})
	keeper.SetClassRevealQueueEntry(ctx, types.ClassRevealQueueEntry{
		RevealTime: revealTime4,
		ClassId:    classId4,
	})

	resultClassIds := getAllClassIdsFromQueue(ctx, keeper)

	require.Equal(t, 4, len(resultClassIds))
}

func TestClassRevealQueueRemove(t *testing.T) {
	keeper, ctx := keepertest.LikenftKeeper(t)

	classId1 := "likenft10"
	revealTime1, err := time.Parse(time.RFC3339, "2020-01-01T00:00:00Z")
	require.NoError(t, err)

	classId2 := "likenft11"
	revealTime2, err := time.Parse(time.RFC3339, "2023-01-01T00:00:00Z")
	require.NoError(t, err)

	classId3 := "likenft12"
	revealTime3, err := time.Parse(time.RFC3339, "2027-01-01T00:00:00Z")
	require.NoError(t, err)

	classId4 := "likenft13"
	revealTime4, err := time.Parse(time.RFC3339, "2009-01-01T00:00:00Z")
	require.NoError(t, err)

	keeper.SetClassRevealQueueEntry(ctx, types.ClassRevealQueueEntry{
		RevealTime: revealTime1,
		ClassId:    classId1,
	})
	keeper.SetClassRevealQueueEntry(ctx, types.ClassRevealQueueEntry{
		RevealTime: revealTime2,
		ClassId:    classId2,
	})
	keeper.SetClassRevealQueueEntry(ctx, types.ClassRevealQueueEntry{
		RevealTime: revealTime3,
		ClassId:    classId3,
	})
	keeper.SetClassRevealQueueEntry(ctx, types.ClassRevealQueueEntry{
		RevealTime: revealTime4,
		ClassId:    classId4,
	})

	// Remove valid entry
	keeper.RemoveClassRevealQueueEntry(ctx, revealTime2, classId2)
	require.NoError(t, err)
	validClassIds := getAllClassIdsFromQueue(ctx, keeper)
	require.Equal(t, 3, len(validClassIds))
}

func TestClassRevealQueueUpdate(t *testing.T) {
	keeper, ctx := keepertest.LikenftKeeper(t)

	classId1 := "likenft10"
	revealTime1, err := time.Parse(time.RFC3339, "2020-01-01T00:00:00Z")
	require.NoError(t, err)

	classId2 := "likenft11"
	revealTime2, err := time.Parse(time.RFC3339, "2023-01-01T00:00:00Z")
	require.NoError(t, err)

	classId3 := "likenft12"
	revealTime3, err := time.Parse(time.RFC3339, "2027-01-01T00:00:00Z")
	require.NoError(t, err)

	classId4 := "likenft13"
	revealTime4, err := time.Parse(time.RFC3339, "2009-01-01T00:00:00Z")
	require.NoError(t, err)

	keeper.SetClassRevealQueueEntry(ctx, types.ClassRevealQueueEntry{
		RevealTime: revealTime1,
		ClassId:    classId1,
	})
	keeper.SetClassRevealQueueEntry(ctx, types.ClassRevealQueueEntry{
		RevealTime: revealTime2,
		ClassId:    classId2,
	})
	keeper.SetClassRevealQueueEntry(ctx, types.ClassRevealQueueEntry{
		RevealTime: revealTime3,
		ClassId:    classId3,
	})
	keeper.SetClassRevealQueueEntry(ctx, types.ClassRevealQueueEntry{
		RevealTime: revealTime4,
		ClassId:    classId4,
	})

	// Update entry
	updatedRevealTime, err := time.Parse(time.RFC3339, "2099-01-01T00:00:00Z")
	require.NoError(t, err)

	keeper.UpdateClassRevealQueueEntry(ctx, revealTime4, classId4, updatedRevealTime)

	classIds := getAllClassIdsFromQueue(ctx, keeper)
	require.Equal(t, []string{
		classId1,
		classId2,
		classId3,
		classId4,
	}, classIds)
}

func TestClassRevealQueueSorting(t *testing.T) {
	keeper, ctx := keepertest.LikenftKeeper(t)

	classId1 := "likenft10"
	revealTime1, err := time.Parse(time.RFC3339, "2020-01-01T00:00:00Z")
	require.NoError(t, err)

	classId2 := "likenft11"
	revealTime2, err := time.Parse(time.RFC3339, "2023-01-01T00:00:00Z")
	require.NoError(t, err)

	classId3 := "likenft12"
	revealTime3, err := time.Parse(time.RFC3339, "2027-01-01T00:00:00Z")
	require.NoError(t, err)

	classId4 := "likenft13"
	revealTime4, err := time.Parse(time.RFC3339, "2009-01-01T00:00:00Z")
	require.NoError(t, err)

	keeper.SetClassRevealQueueEntry(ctx, types.ClassRevealQueueEntry{
		RevealTime: revealTime1,
		ClassId:    classId1,
	})
	keeper.SetClassRevealQueueEntry(ctx, types.ClassRevealQueueEntry{
		RevealTime: revealTime2,
		ClassId:    classId2,
	})
	keeper.SetClassRevealQueueEntry(ctx, types.ClassRevealQueueEntry{
		RevealTime: revealTime3,
		ClassId:    classId3,
	})
	keeper.SetClassRevealQueueEntry(ctx, types.ClassRevealQueueEntry{
		RevealTime: revealTime4,
		ClassId:    classId4,
	})

	classIds := getAllClassIdsFromQueue(ctx, keeper)
	require.Equal(t, 4, len(classIds))

	require.Equal(t, []string{
		classId4,
		classId1,
		classId2,
		classId3,
	}, classIds)
}
