package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/likecoin/likecoin-chain/v3/backport/cosmos-sdk/v0.46.0-alpha2/x/nft"
	"github.com/likecoin/likecoin-chain/v3/x/likenft/types"
)

func (k Keeper) GetClass(ctx sdk.Context, classId string) (nft.Class, types.ClassData, error) {
	// Verify class exists
	class, found := k.nftKeeper.GetClass(ctx, classId)
	if !found {
		return class, types.ClassData{}, types.ErrNftClassNotFound.Wrapf("Class id %s not found", classId)
	}

	// Unmarshal class data
	var classData types.ClassData
	if err := classData.Unmarshal(class.Data.Value); err != nil {
		return class, classData, types.ErrFailedToUnmarshalData.Wrapf(err.Error())
	}

	return class, classData, nil
}
