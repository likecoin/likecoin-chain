package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/likecoin/likechain/x/likenft/types"
)

func (k msgServer) DeleteMintableNFT(goCtx context.Context, msg *types.MsgDeleteMintableNFT) (*types.MsgDeleteMintableNFTResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	parentAndOwner, err := k.getParentOwnerAndValidateReqToMutateMintableNFT(ctx, msg.Creator, msg.ClassId)
	if err != nil {
		return nil, err
	}

	// check id already exists
	if _, exists := k.GetMintableNFT(ctx, msg.ClassId, msg.Id); !exists {
		return nil, types.ErrMintableNftNotFound
	}

	// remove record
	k.RemoveMintableNFT(ctx, msg.ClassId, msg.Id)

	// Emit event
	ctx.EventManager().EmitTypedEvent(&types.EventDeleteMintableNFT{
		ClassId:                 msg.ClassId,
		MintableNftId:           msg.Id,
		ClassParentIscnIdPrefix: parentAndOwner.ClassParent.IscnIdPrefix,
		ClassParentAccount:      parentAndOwner.ClassParent.Account,
	})

	return &types.MsgDeleteMintableNFTResponse{}, nil
}
