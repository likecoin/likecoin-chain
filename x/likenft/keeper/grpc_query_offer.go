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
		var offer types.Offer
		if err := k.cdc.Unmarshal(value, &offer); err != nil {
			return err
		}

		offers = append(offers, offer)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryOfferIndexResponse{Offer: offers, Pagination: pageRes}, nil
}

func (k Keeper) Offer(c context.Context, req *types.QueryOfferRequest) (*types.QueryOfferResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	val, found := k.GetOffer(
		ctx,
		req.ClassId,
		req.NftId,
		req.Buyer,
	)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryOfferResponse{Offer: val}, nil
}
