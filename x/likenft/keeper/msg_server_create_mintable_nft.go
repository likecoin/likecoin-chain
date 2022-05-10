package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/likecoin/likechain/x/likenft/types"
)

func (k msgServer) CreateMintableNFT(goCtx context.Context, msg *types.MsgCreateMintableNFT) (*types.MsgCreateMintableNFTResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	parentAndOwner, err := k.getParentOwnerAndValidateReqToMutateMintableNFT(ctx, msg.Creator, msg.ClassId, true)
	if err != nil {
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
		ClassParentIscnIdPrefix: parentAndOwner.ClassParent.IscnIdPrefix,
		ClassParentAccount:      parentAndOwner.ClassParent.Account,
	})

	return &types.MsgCreateMintableNFTResponse{}, nil
}
