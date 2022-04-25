package likenft

import (
	"time"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/likecoin/likechain/x/likenft/keeper"
	"github.com/likecoin/likechain/x/likenft/types"
)

func tryRevealClassCatchPanic(ctx sdk.Context, keeper keeper.Keeper, classId string) error {
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()
	err = keeper.RevealMintableNFTs(ctx, classId)
	return err
}

// EndBlocker called every block, process class reveal queue.
func EndBlocker(ctx sdk.Context, keeper keeper.Keeper) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyEndBlocker)

	// Reveal class when its reveal time is reached.
	keeper.IterateClassRevealQueue(ctx, func(entry types.ClassRevealQueueEntry) bool {
		if entry.RevealTime.After(ctx.BlockHeader().Time) {
			// Processed all pending entries already, stop
			return true
		}

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
