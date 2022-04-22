package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/likecoin/likechain/x/likenft/types"
)

func (k msgServer) UpdateClaimableNFT(goCtx context.Context, msg *types.MsgUpdateClaimableNFT) (*types.MsgUpdateClaimableNFTResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Handling the message
	_ = ctx

	return &types.MsgUpdateClaimableNFTResponse{}, nil
}
