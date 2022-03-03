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

func (k Keeper) concretizeClassesByISCN(ctx sdk.Context, val types.ClassesByISCN) types.ConcreteClassesByISCN {

	classes := make([]*nft.Class, len(val.ClassIds))
	for i, classId := range val.ClassIds {
		class, found := k.nftKeeper.GetClass(ctx, classId)
		if found {
			classes[i] = &class
		} else {
			classes[i] = nil
		}
	}

	return types.ConcreteClassesByISCN{
		IscnIdPrefix: val.IscnIdPrefix,
		Classes:      classes,
	}
}

func (k Keeper) ClassesByISCNIndex(c context.Context, req *types.QueryClassesByISCNIndexRequest) (*types.QueryClassesByISCNIndexResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var concreteClassesByISCNs []types.ConcreteClassesByISCN
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	classesByISCNStore := prefix.NewStore(store, types.KeyPrefix(types.ClassesByISCNKeyPrefix))

	pageRes, err := query.Paginate(classesByISCNStore, req.Pagination, func(key []byte, value []byte) error {
		var classesByISCN types.ClassesByISCN
		if err := k.cdc.Unmarshal(value, &classesByISCN); err != nil {
			return err
		}

		concretized := k.concretizeClassesByISCN(ctx, classesByISCN)

		concreteClassesByISCNs = append(concreteClassesByISCNs, concretized)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryClassesByISCNIndexResponse{ClassesByISCN: concreteClassesByISCNs, Pagination: pageRes}, nil
}

func (k Keeper) ClassesByISCN(c context.Context, req *types.QueryClassesByISCNRequest) (*types.QueryClassesByISCNResponse, error) {
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

	concretized := k.concretizeClassesByISCN(ctx, val)

	return &types.QueryClassesByISCNResponse{ClassesByISCN: concretized}, nil
}
