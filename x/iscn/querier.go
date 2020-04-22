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
		case QueryParams:
			return queryParams(ctx, req, k)
		case QueryIscnRecord:
			return queryRecord(ctx, req, k)
		case QueryAuthor:
			return queryAuthor(ctx, req, k)
		default:
			return nil, sdk.ErrUnknownRequest("unknown iscn query endpoint")
		}
	}
}

func queryParams(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	params := k.GetParams(ctx)

	res, err := codec.MarshalJSONIndent(ModuleCdc, params)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("failed to JSON marshal result: %s", err.Error()))
	}

	return res, nil
}

func queryRecord(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	params := types.QueryRecordParams{}
	err := types.ModuleCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("failed to unmarshal JSON request: %s", err.Error()))
	}
	record := k.GetIscnRecord(ctx, params.Id)

	res, err := codec.MarshalJSONIndent(ModuleCdc, record)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("failed to JSON marshal result: %s", err.Error()))
	}

	return res, nil
}

func queryAuthor(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	params := types.QueryAuthorParams{}
	err := types.ModuleCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("failed to unmarshal JSON request: %s", err.Error()))
	}
	author := k.GetAuthor(ctx, params.Cid)

	res, err := codec.MarshalJSONIndent(ModuleCdc, author)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("failed to JSON marshal result: %s", err.Error()))
	}

	return res, nil
}
