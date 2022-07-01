package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	iscntypes "github.com/likecoin/likecoin-chain/v3/x/iscn/types"
	"github.com/likecoin/likecoin-chain/v3/x/likenft/types"
)

func (k Keeper) resolveIscnIdAndRecord(ctx sdk.Context, iscnIdStr string) (iscntypes.IscnId, iscntypes.ContentIdRecord, error) {
	iscnId, err := iscntypes.ParseIscnId(iscnIdStr)
	if err != nil {
		return iscnId, iscntypes.ContentIdRecord{}, types.ErrInvalidIscnId.Wrapf("%s", err.Error())
	}
	iscnRecord := k.iscnKeeper.GetContentIdRecord(ctx, iscnId.Prefix)
	if iscnRecord == nil {
		return iscnId, iscntypes.ContentIdRecord{}, types.ErrIscnRecordNotFound.Wrapf("ISCN %s not found", iscnId.Prefix.String())
	}
	return iscnId, *iscnRecord, nil
}
