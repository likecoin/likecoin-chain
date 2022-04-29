package keeper

import (
	"context"
	"fmt"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/likecoin/likechain/backport/cosmos-sdk/v0.46.0-alpha2/x/nft"
	iscntypes "github.com/likecoin/likechain/x/iscn/types"
	"github.com/likecoin/likechain/x/likenft/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) ClassesByISCNIndex(c context.Context, req *types.QueryClassesByISCNIndexRequest) (*types.QueryClassesByISCNIndexResponse, error) {
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

	return &types.QueryClassesByISCNIndexResponse{ClassesByIscn: classesByISCNs, Pagination: pageRes}, nil
}

func (k Keeper) ClassesByISCN(c context.Context, req *types.QueryClassesByISCNRequest) (*types.QueryClassesByISCNResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	iscnId, err := iscntypes.ParseIscnId(req.IscnIdPrefix)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("invalid iscn id: %s", err.Error()))
	}

	val, found := k.GetClassesByISCN(
		ctx,
		iscnId.Prefix.String(),
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

	return &types.QueryClassesByISCNResponse{
		IscnIdPrefix: iscnId.Prefix.String(),
		Classes:      classes,
		Pagination:   pageRes,
	}, nil
}
