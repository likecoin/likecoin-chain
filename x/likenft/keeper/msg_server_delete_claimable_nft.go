package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/likecoin/likechain/x/likenft/types"
)

func (k msgServer) DeleteClaimableNFT(goCtx context.Context, msg *types.MsgDeleteClaimableNFT) (*types.MsgDeleteClaimableNFTResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	parentAndOwner, err := k.getParentOwnerAndValidateReqToMutateClaimableNFT(ctx, msg.Creator, msg.ClassId)
	if err != nil {
		return nil, err
	}

	// check id already exists
	if _, exists := k.GetClaimableNFT(ctx, msg.ClassId, msg.Id); !exists {
		return nil, types.ErrClaimableNftNotFound
	}

	// remove record
	k.RemoveClaimableNFT(ctx, msg.ClassId, msg.Id)

	// Emit event
	ctx.EventManager().EmitTypedEvent(&types.EventDeleteClaimableNFT{
		ClassId:                 msg.ClassId,
		ClaimableNFTId:          msg.Id,
		ClassParentIscnIdPrefix: parentAndOwner.ClassParent.IscnIdPrefix,
		ClassParentAccount:      parentAndOwner.ClassParent.Account,
	})

	return &types.MsgDeleteClaimableNFTResponse{}, nil
}
