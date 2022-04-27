package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/likecoin/likechain/x/likenft/types"
)

func (k msgServer) UpdateMintableNFT(goCtx context.Context, msg *types.MsgUpdateMintableNFT) (*types.MsgUpdateMintableNFTResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	parentAndOwner, err := k.getParentOwnerAndValidateReqToMutateMintableNFT(ctx, msg.Creator, msg.ClassId)
	if err != nil {
		return nil, err
	}

	// check id already exists
	if _, exists := k.GetMintableNFT(ctx, msg.ClassId, msg.Id); !exists {
		return nil, types.ErrMintableNftNotFound
	}

	// set record
	k.SetMintableNFT(ctx, types.MintableNFT{
		ClassId: msg.ClassId,
		Id:      msg.Id,
		Input:   msg.Input,
	})

	// Emit event
	ctx.EventManager().EmitTypedEvent(&types.EventUpdateMintableNFT{
		ClassId:                 msg.ClassId,
		MintableNFTId:           msg.Id,
		ClassParentIscnIdPrefix: parentAndOwner.ClassParent.IscnIdPrefix,
		ClassParentAccount:      parentAndOwner.ClassParent.Account,
	})

	return &types.MsgUpdateMintableNFTResponse{}, nil
}
