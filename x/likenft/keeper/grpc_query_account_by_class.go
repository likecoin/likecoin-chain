package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/likecoin/likecoin-chain/v4/x/likenft/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) AccountByClass(goCtx context.Context, req *types.QueryAccountByClassRequest) (*types.QueryAccountByClassResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// Get account from class data
	class, classData, err := k.GetClass(ctx, req.ClassId)
	if err != nil {
		return nil, err
	}
	// Check parent is class
	if classData.Parent.Type != types.ClassParentType_ACCOUNT {
		return nil, types.ErrNftClassNotRelatedToAnyAccount.Wrapf("NFT Class is related to a %s, not account", classData.Parent.Type.String())
	}
	if err := k.validateClassParentRelation(ctx, class.Id, classData.Parent); err != nil {
		return nil, err
	}

	// Return account address
	return &types.QueryAccountByClassResponse{
		Address: classData.Parent.Account,
	}, nil
}
