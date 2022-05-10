package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/likecoin/likechain/backport/cosmos-sdk/v0.46.0-alpha2/x/nft"
	"github.com/likecoin/likechain/x/likenft/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) ClassesByAccountIndex(c context.Context, req *types.QueryClassesByAccountIndexRequest) (*types.QueryClassesByAccountIndexResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var classesByAccounts []types.ClassesByAccount
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	classesByAccountStore := prefix.NewStore(store, types.KeyPrefix(types.ClassesByAccountKeyPrefix))

	pageRes, err := query.Paginate(classesByAccountStore, req.Pagination, func(key []byte, value []byte) error {
		var storeRecord types.ClassesByAccountStoreRecord
		if err := k.cdc.Unmarshal(value, &storeRecord); err != nil {
			return err
		}

		classesByAccounts = append(classesByAccounts, storeRecord.ToPublicRecord())
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryClassesByAccountIndexResponse{ClassesByAccounts: classesByAccounts, Pagination: pageRes}, nil
}

func (k Keeper) ClassesByAccount(c context.Context, req *types.QueryClassesByAccountRequest) (*types.QueryClassesByAccountResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	acc, err := sdk.AccAddressFromBech32(req.Account)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid address")
	}

	val, found := k.GetClassesByAccount(
		ctx,
		acc,
	)
	if !found {
		return nil, status.Error(codes.InvalidArgument, "not found")
	}

	var classes []nft.Class
	pageRes, err := PaginateStringArray(val.ClassIds, req.Pagination, func(i int, val string) error {
		class, found := k.nftKeeper.GetClass(ctx, val)
		if !found { // not found, fill in id and return rest fields as empty
			class.Id = val
		}
		classes = append(classes, class)
		return nil
	}, 20, 50) // TODO refactor this in oursky/likecoin-chain#98
	if err != nil {
		// we will not throw error in onResult, so error must be bad pagination argument
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &types.QueryClassesByAccountResponse{
		Account:    req.Account,
		Classes:    classes,
		Pagination: pageRes,
	}, nil
}
