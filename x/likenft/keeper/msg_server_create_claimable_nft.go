package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/likecoin/likechain/x/likenft/types"
)

func (k msgServer) CreateClaimableNFT(goCtx context.Context, msg *types.MsgCreateClaimableNFT) (*types.MsgCreateClaimableNFTResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Handling the message
	_ = ctx

	return &types.MsgCreateClaimableNFTResponse{}, nil
}
