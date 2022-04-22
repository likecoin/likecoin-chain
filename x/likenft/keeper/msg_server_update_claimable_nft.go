package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/likecoin/likechain/x/likenft/types"
)

func (k msgServer) UpdateClaimableNFT(goCtx context.Context, msg *types.MsgUpdateClaimableNFT) (*types.MsgUpdateClaimableNFTResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := k.validateRequestToMutateClaimableNFT(ctx, msg.Creator, msg.ClassId); err != nil {
		return nil, err
	}

	// check id already exists
	if _, exists := k.GetClaimableNFT(ctx, msg.ClassId, msg.Id); !exists {
		return nil, types.ErrClaimableNftNotFound
	}

	// set record
	k.SetClaimableNFT(ctx, types.ClaimableNFT{
		ClassId: msg.ClassId,
		Id:      msg.Id,
		Input:   msg.Input,
	})

	// TODO emit event
	return &types.MsgUpdateClaimableNFTResponse{}, nil
}
