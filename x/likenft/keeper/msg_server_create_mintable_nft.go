package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/likecoin/likechain/x/likenft/types"
)

func (k msgServer) CreateMintableNFT(goCtx context.Context, msg *types.MsgCreateMintableNFT) (*types.MsgCreateMintableNFTResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	class, classData, err := k.GetClass(ctx, msg.ClassId)
	if err != nil {
		return nil, err
	}
	parent, err := k.ValidateAndRefreshClassParent(ctx, msg.ClassId, classData.Parent)
	if err != nil {
		return nil, err
	}
	if err := k.validateReqToMutateMintableNFT(ctx, msg.Creator, class, classData, parent, true); err != nil {
		return nil, err
	}

	// check id not already exist
	if _, exists := k.GetMintableNFT(ctx, msg.ClassId, msg.Id); exists {
		return nil, types.ErrMintableNftAlreadyExists
	}

	// set record
	k.SetMintableNFT(ctx, types.MintableNFT{
		ClassId: msg.ClassId,
		Id:      msg.Id,
		Input:   msg.Input,
	})

	// Emit event
	ctx.EventManager().EmitTypedEvent(&types.EventCreateMintableNFT{
		ClassId:                 msg.ClassId,
		MintableNftId:           msg.Id,
		ClassParentIscnIdPrefix: parent.IscnIdPrefix,
		ClassParentAccount:      parent.Account,
	})

	return &types.MsgCreateMintableNFTResponse{}, nil
}
