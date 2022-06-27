package keeper

import (
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/likecoin/likechain/x/likenft/types"
)

func (k Keeper) incBlindBoxContentCount(ctx sdk.Context, classId string) error {
	return k._setBlindBoxContentCount(ctx, classId, func(count uint64) uint64 {
		return count + 1
	})
}

func (k Keeper) decBlindBoxContentCount(ctx sdk.Context, classId string) error {
	return k._setBlindBoxContentCount(ctx, classId, func(count uint64) uint64 {
		return count - 1
	})
}

func (k Keeper) setBlindBoxContentCount(ctx sdk.Context, classId string, count uint64) error {
	return k._setBlindBoxContentCount(ctx, classId, func(_ uint64) uint64 {
		return count
	})
}

func (k Keeper) _setBlindBoxContentCount(ctx sdk.Context, classId string, edit func(uint64) uint64) error {
	class, classData, err := k.GetClass(ctx, classId)
	if err != nil {
		return err
	}
	classData.BlindBoxState.ContentCount = edit(classData.BlindBoxState.ContentCount)
	classDataInAny, err := cdctypes.NewAnyWithValue(&classData)
	if err != nil {
		return types.ErrFailedToMarshalData.Wrapf("%s", err.Error())
	}
	class.Data = classDataInAny
	if err := k.nftKeeper.UpdateClass(ctx, class); err != nil {
		return types.ErrFailedToUpdateClass.Wrapf("%s", err.Error())
	}
	return nil
}
