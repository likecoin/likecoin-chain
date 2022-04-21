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

func (k Keeper) ClaimableNFTAll(c context.Context, req *types.QueryAllClaimableNFTRequest) (*types.QueryAllClaimableNFTResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var claimableNFTs []types.ClaimableNFT
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	claimableNFTStore := prefix.NewStore(store, types.KeyPrefix(types.ClaimableNFTKeyPrefix))

	pageRes, err := query.Paginate(claimableNFTStore, req.Pagination, func(key []byte, value []byte) error {
		var claimableNFT types.ClaimableNFT
		if err := k.cdc.Unmarshal(value, &claimableNFT); err != nil {
			return err
		}

		claimableNFTs = append(claimableNFTs, claimableNFT)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllClaimableNFTResponse{ClaimableNFT: claimableNFTs, Pagination: pageRes}, nil
}

func (k Keeper) ClaimableNFT(c context.Context, req *types.QueryGetClaimableNFTRequest) (*types.QueryGetClaimableNFTResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	val, found := k.GetClaimableNFT(
		ctx,
		req.ClassId,
		req.Id,
	)
	if !found {
		return nil, status.Error(codes.InvalidArgument, "not found")
	}

	return &types.QueryGetClaimableNFTResponse{ClaimableNFT: val}, nil
}
