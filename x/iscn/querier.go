package iscn

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/likecoin/likechain/x/iscn/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case QueryRecord:
			return queryRecord(ctx, req, k)
		default:
			return nil, sdk.ErrUnknownRequest("unknown iscn query endpoint")
		}
	}
}

func queryRecord(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	params := types.QueryRecordParams{}
	err := types.ModuleCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("failed to unmarshal JSON request: %s", err.Error()))
	}
	approver := k.GetIscnRecord(ctx, params.Id)

	res, err := codec.MarshalJSONIndent(ModuleCdc, approver)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("failed to JSON marshal result: %s", err.Error()))
	}

	return res, nil
}
