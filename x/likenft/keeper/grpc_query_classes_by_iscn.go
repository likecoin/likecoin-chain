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

func (k Keeper) ClassesByISCNAll(c context.Context, req *types.QueryAllClassesByISCNRequest) (*types.QueryAllClassesByISCNResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var classesByISCNs []types.ClassesByISCN
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	classesByISCNStore := prefix.NewStore(store, types.KeyPrefix(types.ClassesByISCNKeyPrefix))

	pageRes, err := query.Paginate(classesByISCNStore, req.Pagination, func(key []byte, value []byte) error {
		var classesByISCN types.ClassesByISCN
		if err := k.cdc.Unmarshal(value, &classesByISCN); err != nil {
			return err
		}

		classesByISCNs = append(classesByISCNs, classesByISCN)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllClassesByISCNResponse{ClassesByISCN: classesByISCNs, Pagination: pageRes}, nil
}

func (k Keeper) ClassesByISCN(c context.Context, req *types.QueryGetClassesByISCNRequest) (*types.QueryGetClassesByISCNResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	val, found := k.GetClassesByISCN(
		ctx,
		req.IscnIdPrefix,
	)
	if !found {
		return nil, status.Error(codes.InvalidArgument, "not found")
	}

	return &types.QueryGetClassesByISCNResponse{ClassesByISCN: val}, nil
}
