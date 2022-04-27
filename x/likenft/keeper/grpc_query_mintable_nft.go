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

func (k Keeper) MintableNFTIndex(c context.Context, req *types.QueryMintableNFTIndexRequest) (*types.QueryMintableNFTIndexResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var mintableNFTs []types.MintableNFT
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	mintableNFTStore := prefix.NewStore(store, types.KeyPrefix(types.MintableNFTKeyPrefix))

	pageRes, err := query.Paginate(mintableNFTStore, req.Pagination, func(key []byte, value []byte) error {
		var mintableNFT types.MintableNFT
		if err := k.cdc.Unmarshal(value, &mintableNFT); err != nil {
			return err
		}

		mintableNFTs = append(mintableNFTs, mintableNFT)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryMintableNFTIndexResponse{MintableNFT: mintableNFTs, Pagination: pageRes}, nil
}

func (k Keeper) MintableNFTs(c context.Context, req *types.QueryMintableNFTsRequest) (*types.QueryMintableNFTsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var mintableNFTs []types.MintableNFT
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	subStore := prefix.NewStore(store, append(types.KeyPrefix(types.MintableNFTKeyPrefix), types.MintableNFTsKey(req.ClassId)...))

	pageRes, err := query.Paginate(subStore, req.Pagination, func(key []byte, value []byte) error {
		var mintableNFT types.MintableNFT
		if err := k.cdc.Unmarshal(value, &mintableNFT); err != nil {
			return err
		}

		mintableNFTs = append(mintableNFTs, mintableNFT)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryMintableNFTsResponse{MintableNFTs: mintableNFTs, Pagination: pageRes}, nil
}

func (k Keeper) MintableNFT(c context.Context, req *types.QueryMintableNFTRequest) (*types.QueryMintableNFTResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	val, found := k.GetMintableNFT(
		ctx,
		req.ClassId,
		req.Id,
	)
	if !found {
		return nil, status.Error(codes.InvalidArgument, "not found")
	}

	return &types.QueryMintableNFTResponse{MintableNFT: val}, nil
}
