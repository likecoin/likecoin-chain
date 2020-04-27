package whitelist

import (
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case QueryApprover:
			return queryApprover(ctx, req, k)
		case QueryWhitelist:
			return queryWhitelist(ctx, req, k)
		default:
			return nil, sdk.ErrUnknownRequest("unknown whitelist query endpoint")
		}
	}
}

func queryApprover(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	approver := k.Approver(ctx)

	res, err := codec.MarshalJSONIndent(ModuleCdc, approver)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("failed to JSON marshal result", err.Error()))
	}

	return res, nil
}

func queryWhitelist(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	whitelist := k.GetWhitelist(ctx)

	res, err := codec.MarshalJSONIndent(ModuleCdc, whitelist)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("failed to JSON marshal result", err.Error()))
	}

	return res, nil
}
