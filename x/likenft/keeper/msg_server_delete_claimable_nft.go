package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/likecoin/likechain/x/likenft/types"
)

func (k msgServer) DeleteClaimableNFT(goCtx context.Context, msg *types.MsgDeleteClaimableNFT) (*types.MsgDeleteClaimableNFTResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := k.validateRequestToMutateClaimableNFT(ctx, msg.Creator, msg.ClassId); err != nil {
		return nil, err
	}

	// check id already exists
	if _, exists := k.GetClaimableNFT(ctx, msg.ClassId, msg.Id); !exists {
		return nil, types.ErrClaimableNftNotFound
	}

	// set record
	k.RemoveClaimableNFT(ctx, msg.ClassId, msg.Id)

	// TODO emit event
	return &types.MsgDeleteClaimableNFTResponse{}, nil
}
