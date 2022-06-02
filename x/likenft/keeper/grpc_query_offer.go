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

func (k Keeper) OfferIndex(c context.Context, req *types.QueryOfferIndexRequest) (*types.QueryOfferIndexResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var offers []types.Offer
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	offerStore := prefix.NewStore(store, types.KeyPrefix(types.OfferKeyPrefix))

	pageRes, err := query.Paginate(offerStore, req.Pagination, func(key []byte, value []byte) error {
		var offer types.OfferStoreRecord
		if err := k.cdc.Unmarshal(value, &offer); err != nil {
			return err
		}

		offers = append(offers, offer.ToPublicRecord())
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryOfferIndexResponse{Offers: offers, Pagination: pageRes}, nil
}

func (k Keeper) Offer(c context.Context, req *types.QueryOfferRequest) (*types.QueryOfferResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	buyer, err := sdk.AccAddressFromBech32(req.Buyer)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	val, found := k.GetOffer(
		ctx,
		req.ClassId,
		req.NftId,
		buyer,
	)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryOfferResponse{Offer: val}, nil
}

func (k Keeper) OffersByClass(goCtx context.Context, req *types.QueryOffersByClassRequest) (*types.QueryOffersByClassResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var offers []types.Offer
	ctx := sdk.UnwrapSDKContext(goCtx)

	store := ctx.KVStore(k.storeKey)
	subStore := prefix.NewStore(store, append(types.KeyPrefix(types.OfferKeyPrefix), types.OffersByClassKey(req.ClassId)...))

	pageRes, err := query.Paginate(subStore, req.Pagination, func(key []byte, value []byte) error {
		var offer types.OfferStoreRecord
		if err := k.cdc.Unmarshal(value, &offer); err != nil {
			return err
		}

		offers = append(offers, offer.ToPublicRecord())
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryOffersByClassResponse{
		Offers:     offers,
		Pagination: pageRes,
	}, nil
}

func (k Keeper) OffersByNFT(goCtx context.Context, req *types.QueryOffersByNFTRequest) (*types.QueryOffersByNFTResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var offers []types.Offer
	ctx := sdk.UnwrapSDKContext(goCtx)

	store := ctx.KVStore(k.storeKey)
	subStore := prefix.NewStore(store, append(types.KeyPrefix(types.OfferKeyPrefix), types.OffersByNFTKey(req.ClassId, req.NftId)...))

	pageRes, err := query.Paginate(subStore, req.Pagination, func(key []byte, value []byte) error {
		var offer types.OfferStoreRecord
		if err := k.cdc.Unmarshal(value, &offer); err != nil {
			return err
		}

		offers = append(offers, offer.ToPublicRecord())
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryOffersByNFTResponse{
		Offers:     offers,
		Pagination: pageRes,
	}, nil
}
