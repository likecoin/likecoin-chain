package likenft

import (
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/likecoin/likechain/x/likenft/keeper"
	"github.com/likecoin/likechain/x/likenft/types"
)

func tryRevealClassCatchPanic(ctx sdk.Context, keeper keeper.Keeper, classId string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%s", r)
		}
	}()
	err = keeper.RevealMintableNFTs(ctx, classId)
	return
}

// EndBlocker called every block, process class reveal queue.
func EndBlocker(ctx sdk.Context, keeper keeper.Keeper) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyEndBlocker)

	// Reveal classes with reveal time < current block header time
	keeper.IterateClassRevealQueueByTime(ctx, ctx.BlockHeader().Time, func(entry types.ClassRevealQueueEntry) (stop bool) {
		err := tryRevealClassCatchPanic(ctx, keeper, entry.ClassId)

		if err != nil {
			ctx.EventManager().EmitTypedEvent(&types.EventRevealClass{
				ClassId: entry.ClassId,
				Success: false,
				Error:   err.Error(),
			})
		} else {
			ctx.EventManager().EmitTypedEvent(&types.EventRevealClass{
				ClassId: entry.ClassId,
				Success: true,
			})
		}

		keeper.RemoveClassRevealQueueEntry(ctx, entry.RevealTime, entry.ClassId)
		return false
	})
}
