package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/likecoin/likecoin-chain/v4/x/likenft/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) RoyaltyConfigIndex(c context.Context, req *types.QueryRoyaltyConfigIndexRequest) (*types.QueryRoyaltyConfigIndexResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var royaltyConfigByClasss []types.RoyaltyConfigByClass
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	royaltyConfigByClassStore := prefix.NewStore(store, types.KeyPrefix(types.RoyaltyConfigByClassKeyPrefix))

	pageRes, err := query.Paginate(royaltyConfigByClassStore, req.Pagination, func(key []byte, value []byte) error {
		var royaltyConfigByClass types.RoyaltyConfigByClass
		if err := k.cdc.Unmarshal(value, &royaltyConfigByClass); err != nil {
			return err
		}

		royaltyConfigByClasss = append(royaltyConfigByClasss, royaltyConfigByClass)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryRoyaltyConfigIndexResponse{RoyaltyConfigByClass: royaltyConfigByClasss, Pagination: pageRes}, nil
}

func (k Keeper) RoyaltyConfig(c context.Context, req *types.QueryRoyaltyConfigRequest) (*types.QueryRoyaltyConfigResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	val, found := k.GetRoyaltyConfig(
		ctx,
		req.ClassId,
	)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryRoyaltyConfigResponse{RoyaltyConfig: val}, nil
}
