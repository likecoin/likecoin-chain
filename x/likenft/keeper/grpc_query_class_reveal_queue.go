package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/likecoin/likechain/x/likenft/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) ClassRevealQueueAll(c context.Context, req *types.QueryAllClassRevealQueueRequest) (*types.QueryAllClassRevealQueueResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var classRevealQueues []types.ClassRevealQueue
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	classRevealQueueStore := prefix.NewStore(store, types.KeyPrefix(types.ClassRevealQueueKeyPrefix))

	pageRes, err := query.Paginate(classRevealQueueStore, req.Pagination, func(key []byte, value []byte) error {
		var classRevealQueue types.ClassRevealQueue
		if err := k.cdc.Unmarshal(value, &classRevealQueue); err != nil {
			return err
		}

		classRevealQueues = append(classRevealQueues, classRevealQueue)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllClassRevealQueueResponse{ClassRevealQueue: classRevealQueues, Pagination: pageRes}, nil
}

func (k Keeper) ClassRevealQueue(c context.Context, req *types.QueryGetClassRevealQueueRequest) (*types.QueryGetClassRevealQueueResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	val, found := k.GetClassRevealQueue(
		ctx,
		req.RevealTime,
		req.ClassId,
	)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryGetClassRevealQueueResponse{ClassRevealQueue: val}, nil
}
