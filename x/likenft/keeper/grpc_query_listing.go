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

func (k Keeper) ListingIndex(c context.Context, req *types.QueryListingIndexRequest) (*types.QueryListingIndexResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var listings []types.Listing
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	listingStore := prefix.NewStore(store, types.KeyPrefix(types.ListingKeyPrefix))

	pageRes, err := query.Paginate(listingStore, req.Pagination, func(key []byte, value []byte) error {
		var storeRecord types.ListingStoreRecord
		if err := k.cdc.Unmarshal(value, &storeRecord); err != nil {
			return err
		}

		listings = append(listings, storeRecord.ToPublicRecord())
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryListingIndexResponse{Listings: listings, Pagination: pageRes}, nil
}

func (k Keeper) Listing(c context.Context, req *types.QueryListingRequest) (*types.QueryListingResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	seller, err := sdk.AccAddressFromBech32(req.Seller)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	val, found := k.GetListing(
		ctx,
		req.ClassId,
		req.NftId,
		seller,
	)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryListingResponse{Listing: val.ToPublicRecord()}, nil
}

func (k Keeper) ListingsByClass(goCtx context.Context, req *types.QueryListingsByClassRequest) (*types.QueryListingsByClassResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var listings []types.Listing
	ctx := sdk.UnwrapSDKContext(goCtx)

	store := ctx.KVStore(k.storeKey)
	subStore := prefix.NewStore(store, append(types.KeyPrefix(types.ListingKeyPrefix), types.ListingsByClassKey(req.ClassId)...))

	pageRes, err := query.Paginate(subStore, req.Pagination, func(key []byte, value []byte) error {
		var storeRecord types.ListingStoreRecord
		if err := k.cdc.Unmarshal(value, &storeRecord); err != nil {
			return err
		}

		listings = append(listings, storeRecord.ToPublicRecord())
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryListingsByClassResponse{
		Listings:   listings,
		Pagination: pageRes,
	}, nil
}

func (k Keeper) ListingsByNFT(goCtx context.Context, req *types.QueryListingsByNFTRequest) (*types.QueryListingsByNFTResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var listings []types.Listing
	ctx := sdk.UnwrapSDKContext(goCtx)

	store := ctx.KVStore(k.storeKey)
	subStore := prefix.NewStore(store, append(types.KeyPrefix(types.ListingKeyPrefix), types.ListingsByNFTKey(req.ClassId, req.NftId)...))

	pageRes, err := query.Paginate(subStore, req.Pagination, func(key []byte, value []byte) error {
		var storeRecord types.ListingStoreRecord
		if err := k.cdc.Unmarshal(value, &storeRecord); err != nil {
			return err
		}

		listings = append(listings, storeRecord.ToPublicRecord())
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryListingsByNFTResponse{
		Listings:   listings,
		Pagination: pageRes,
	}, nil
}
