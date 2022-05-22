package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/likecoin/likechain/x/likenft/types"
)

func (k Keeper) NewClassParentFromInput(ctx sdk.Context, input types.ClassParentInput, parentAccBech32 string) (types.ClassParentWithOwner, error) {
	if input.Type == types.ClassParentType_ISCN {
		iscnId, iscnRecord, err := k.resolveIscnIdAndRecord(ctx, input.IscnIdPrefix)
		if err != nil {
			return types.ClassParentWithOwner{}, err
		}
		return types.ClassParentWithOwner{
			ClassParent: types.ClassParent{
				Type:              types.ClassParentType_ISCN,
				IscnIdPrefix:      iscnId.Prefix.String(),
				IscnVersionAtMint: iscnRecord.LatestVersion,
			},
			Owner: iscnRecord.OwnerAddress(),
		}, nil
	} else if input.Type == types.ClassParentType_ACCOUNT {
		parentAcc, err := sdk.AccAddressFromBech32(parentAccBech32)
		if err != nil {
			return types.ClassParentWithOwner{}, sdkerrors.ErrInvalidAddress.Wrapf(err.Error())
		}
		return types.ClassParentWithOwner{
			ClassParent: types.ClassParent{
				Type:    types.ClassParentType_ACCOUNT,
				Account: parentAcc.String(),
			},
			Owner: parentAcc,
		}, nil
	} else {
		return types.ClassParentWithOwner{}, sdkerrors.ErrInvalidRequest.Wrapf("Unsupported parent type %s in nft class", input.Type.String())
	}
}

func (k Keeper) ValidateAndRefreshClassParent(ctx sdk.Context, classId string, parent types.ClassParent) (types.ClassParentWithOwner, error) {
	if err := k.validateClassParentRelation(ctx, classId, parent); err != nil {
		return types.ClassParentWithOwner{}, err
	}
	return k.NewClassParentFromInput(ctx, parent.ToInput(), parent.Account)
}
