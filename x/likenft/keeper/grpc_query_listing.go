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

func (k Keeper) ListingAll(c context.Context, req *types.QueryAllListingRequest) (*types.QueryAllListingResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var listings []types.Listing
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	listingStore := prefix.NewStore(store, types.KeyPrefix(types.ListingKeyPrefix))

	pageRes, err := query.Paginate(listingStore, req.Pagination, func(key []byte, value []byte) error {
		var listing types.Listing
		if err := k.cdc.Unmarshal(value, &listing); err != nil {
			return err
		}

		listings = append(listings, listing)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllListingResponse{Listings: listings, Pagination: pageRes}, nil
}

func (k Keeper) Listing(c context.Context, req *types.QueryGetListingRequest) (*types.QueryGetListingResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	val, found := k.GetListing(
		ctx,
		req.ClassId,
		req.NftId,
		req.Seller,
	)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryGetListingResponse{Listing: val}, nil
}
