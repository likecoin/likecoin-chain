package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/likecoin/likechain/x/likenft/types"
)

func (k msgServer) DeleteClaimableNFT(goCtx context.Context, msg *types.MsgDeleteClaimableNFT) (*types.MsgDeleteClaimableNFTResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Handling the message
	_ = ctx

	return &types.MsgDeleteClaimableNFTResponse{}, nil
}
