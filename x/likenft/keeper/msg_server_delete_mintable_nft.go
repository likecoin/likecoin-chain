package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/likecoin/likechain/x/likenft/types"
)

func (k msgServer) DeleteBlindBoxContent(goCtx context.Context, msg *types.MsgDeleteBlindBoxContent) (*types.MsgDeleteBlindBoxContentResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	class, classData, err := k.GetClass(ctx, msg.ClassId)
	if err != nil {
		return nil, err
	}
	parent, err := k.ValidateAndRefreshClassParent(ctx, msg.ClassId, classData.Parent)
	if err != nil {
		return nil, err
	}
	if err := k.validateReqToMutateBlindBoxContent(ctx, msg.Creator, class, classData, parent, false); err != nil {
		return nil, err
	}

	// check id already exists
	if _, exists := k.GetBlindBoxContent(ctx, msg.ClassId, msg.Id); !exists {
		return nil, types.ErrBlindBoxContentNotFound
	}

	// remove record
	k.RemoveBlindBoxContent(ctx, msg.ClassId, msg.Id)

	// Emit event
	ctx.EventManager().EmitTypedEvent(&types.EventDeleteBlindBoxContent{
		ClassId:                 msg.ClassId,
		ContentId:               msg.Id,
		ClassParentIscnIdPrefix: parent.IscnIdPrefix,
		ClassParentAccount:      parent.Account,
	})

	return &types.MsgDeleteBlindBoxContentResponse{}, nil
}
