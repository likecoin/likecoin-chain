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

func (k Keeper) ClassesByAccountAll(c context.Context, req *types.QueryAllClassesByAccountRequest) (*types.QueryAllClassesByAccountResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var classesByAccounts []types.ClassesByAccount
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	classesByAccountStore := prefix.NewStore(store, types.KeyPrefix(types.ClassesByAccountKeyPrefix))

	pageRes, err := query.Paginate(classesByAccountStore, req.Pagination, func(key []byte, value []byte) error {
		var classesByAccount types.ClassesByAccount
		if err := k.cdc.Unmarshal(value, &classesByAccount); err != nil {
			return err
		}

		classesByAccounts = append(classesByAccounts, classesByAccount)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllClassesByAccountResponse{ClassesByAccount: classesByAccounts, Pagination: pageRes}, nil
}

func (k Keeper) ClassesByAccount(c context.Context, req *types.QueryGetClassesByAccountRequest) (*types.QueryGetClassesByAccountResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	val, found := k.GetClassesByAccount(
		ctx,
		req.Account,
	)
	if !found {
		return nil, status.Error(codes.InvalidArgument, "not found")
	}

	return &types.QueryGetClassesByAccountResponse{ClassesByAccount: val}, nil
}
