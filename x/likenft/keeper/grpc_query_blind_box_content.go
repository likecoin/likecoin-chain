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

func (k Keeper) BlindBoxContentIndex(c context.Context, req *types.QueryBlindBoxContentIndexRequest) (*types.QueryBlindBoxContentIndexResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var blindBoxContents []types.BlindBoxContent
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	blindBoxContentStore := prefix.NewStore(store, types.KeyPrefix(types.BlindBoxContentKeyPrefix))

	pageRes, err := query.Paginate(blindBoxContentStore, req.Pagination, func(key []byte, value []byte) error {
		var blindBoxContent types.BlindBoxContent
		if err := k.cdc.Unmarshal(value, &blindBoxContent); err != nil {
			return err
		}

		blindBoxContents = append(blindBoxContents, blindBoxContent)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryBlindBoxContentIndexResponse{BlindBoxContents: blindBoxContents, Pagination: pageRes}, nil
}

func (k Keeper) BlindBoxContents(c context.Context, req *types.QueryBlindBoxContentsRequest) (*types.QueryBlindBoxContentsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var blindBoxContents []types.BlindBoxContent
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	subStore := prefix.NewStore(store, append(types.KeyPrefix(types.BlindBoxContentKeyPrefix), types.BlindBoxContentsKey(req.ClassId)...))

	pageRes, err := query.Paginate(subStore, req.Pagination, func(key []byte, value []byte) error {
		var blindBoxContent types.BlindBoxContent
		if err := k.cdc.Unmarshal(value, &blindBoxContent); err != nil {
			return err
		}

		blindBoxContents = append(blindBoxContents, blindBoxContent)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryBlindBoxContentsResponse{BlindBoxContents: blindBoxContents, Pagination: pageRes}, nil
}

func (k Keeper) BlindBoxContent(c context.Context, req *types.QueryBlindBoxContentRequest) (*types.QueryBlindBoxContentResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	val, found := k.GetBlindBoxContent(
		ctx,
		req.ClassId,
		req.Id,
	)
	if !found {
		return nil, status.Error(codes.InvalidArgument, "not found")
	}

	return &types.QueryBlindBoxContentResponse{BlindBoxContent: val}, nil
}
