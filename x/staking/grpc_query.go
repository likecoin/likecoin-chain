package staking

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	prefixstore "github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"

	"github.com/cosmos/cosmos-sdk/x/staking/keeper"
	"github.com/cosmos/cosmos-sdk/x/staking/types"
)

// keeper.DelegationsToDelegationResponses is so inefficient that it does one query for bonded token + one query for
// validator for each delegation, which are both unchanged
func DelegationsToDelegationResponses(
	ctx sdk.Context, k keeper.Keeper, delegations types.Delegations, valAddr sdk.ValAddress,
) (types.DelegationResponses, error) {
	val, found := k.GetValidator(ctx, valAddr)
	if !found {
		return nil, types.ErrNoValidatorFound
	}
	bondDenom := k.BondDenom(ctx)
	resp := make(types.DelegationResponses, len(delegations))
	for i, del := range delegations {
		resp[i] = types.NewDelegationResp(
			del.GetDelegatorAddr(),
			valAddr,
			del.Shares,
			sdk.NewCoin(bondDenom, val.TokensFromShares(del.Shares).TruncateInt()),
		)
	}

	return resp, nil
}

func (q *Querier) ValidatorDelegations(c context.Context, req *types.QueryValidatorDelegationsRequest) (*types.QueryValidatorDelegationsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	blockHeight := ctx.BlockHeight()
	height := q.GetHeight()
	ctx.Logger().Debug(
		"Querying ValidatorDelegations from indexer",
		"req", req,
		"block_height", blockHeight,
		"indexer_height", height,
	)
	if blockHeight != int64(height) {
		// Indexing store does not store old heights, so fallback to the slow path
		// Maybe also because the indexing is not done yet
		ctx.Logger().Info(
			"Querying height not current height in ValidatorDelegations query (could be due to index not yet built), fallback to slow path",
			"block_height", blockHeight,
			"indexer_height", height,
		)
		return q.QueryServer.ValidatorDelegations(c, req)
	}
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	if req.ValidatorAddr == "" {
		return nil, status.Error(codes.InvalidArgument, "validator address cannot be empty")
	}
	valAddr, err := sdk.ValAddressFromBech32(req.ValidatorAddr)
	if err != nil {
		return nil, err
	}
	store := prefixstore.NewStore(q.readStore, getIndexValidatorPrefix(valAddr))
	delegations := []types.Delegation{}
	pageRes, err := query.FilteredPaginate(store, req.Pagination, func(key, value []byte, accumulate bool) (bool, error) {
		if accumulate {
			delegation := types.MustUnmarshalDelegation(q.cdc, value)
			delegations = append(delegations, delegation)
		}
		return true, nil
	})
	heightAfterQuery := q.GetHeight()
	if heightAfterQuery != height {
		// Write occurred during query, queried data could be invalid.
		// Since height changed, the query must be querying an old height (relative to the new height) which we don't
		// support, so fallback to the slow path.

		// Note:
		// In most of the cases, the user actually wants to query the newest height (i.e. height = 0 in query),
		// but the SDK will automatically set the height to current height before passing to the QueryServer,
		// so maybe in this case we actually should not fallback to the slow path?
		// But the SDK will set the result's height to the old height... So it seems not possible to get consistent result.
		ctx.Logger().Debug(
			"Write occurred during ValidatorDelegations query, fallback to slow path",
			"height", height,
			"height_after_query", heightAfterQuery,
		)
		return q.QueryServer.ValidatorDelegations(c, req)
	}
	delResponses := types.DelegationResponses{}
	if len(delegations) > 0 {
		delResponses, err = DelegationsToDelegationResponses(
			ctx, *q.keeper, delegations, delegations[0].GetValidatorAddr(),
		)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}
	return &types.QueryValidatorDelegationsResponse{
		DelegationResponses: delResponses, Pagination: pageRes,
	}, nil
}
