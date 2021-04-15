package keeper

import (
	context "context"
	"fmt"

	gocid "github.com/ipfs/go-cid"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/likecoin/likechain/x/iscn/types"
)

var _ types.QueryServer = Keeper{}

func queryIscnRecord(ctx sdk.Context, k Keeper, iscnId IscnId) (*types.Record, error) {
	latestVersion := k.GetIscnIdVersion(ctx, iscnId)
	if latestVersion == 0 {
		return nil, sdkerrors.Wrapf(types.ErrRecordNotFound, "%s", iscnId.String())
	}
	if iscnId.Version == 0 {
		iscnId.Version = latestVersion
	}
	cid := k.GetIscnIdCid(ctx, iscnId)
	if cid == nil {
		return nil, sdkerrors.Wrapf(types.ErrRecordNotFound, "%s", iscnId.String())
	}
	owner := k.GetIscnIdOwner(ctx, iscnId)
	record := k.GetCidBlock(ctx, *cid)
	return &types.Record{
		IscnId:              iscnId.String(),
		Owner:               owner.String(),
		Ipld:                cid.String(),
		LatestRecordVersion: latestVersion,
		Record:              types.IscnInput(record),
	}, nil
}

func (k Keeper) IscnRecords(ctx context.Context, req *types.QueryIscnRecordsRequest) (*types.QueryIscnRecordsResponse, error) {
	// TODO: REMOVE debug fmt
	fmt.Printf("GETTING IscnRecords: '%s', '%s'\n", req.IscnId, req.Fingerprint)
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	if len(req.IscnId) > 0 {
		if len(req.Fingerprint) > 0 {
			return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "only one of iscn_id and fingerprint can exist in query parameters")
		}
		id, err := types.ParseIscnID(req.IscnId)
		if err != nil {
			return nil, sdkerrors.Wrapf(types.ErrInvalidIscnId, "%s", err.Error())
		}
		record, err := queryIscnRecord(sdkCtx, k, id)
		if err != nil {
			return nil, err
		}
		return &types.QueryIscnRecordsResponse{Records: []types.Record{*record}}, nil
	} else if len(req.Fingerprint) > 0 {
		records := []types.Record{}
		k.IterateFingerprintCids(sdkCtx, req.Fingerprint, func(cid CID) bool {
			// TODO: pagination?
			iscnId := k.GetCidIscnId(sdkCtx, cid)
			if iscnId == nil {
				// TODO: ???
				return false
			}
			result, err := queryIscnRecord(sdkCtx, k, *iscnId)
			if err != nil {
				// TODO: ???
				return false
			}
			records = append(records, *result)
			return false
		})
		return &types.QueryIscnRecordsResponse{Records: records}, nil
	} else {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "one of iscn_id and fingerprint must exist in query parameters")
	}
}

func (k Keeper) Params(ctx context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	params := k.GetParams(sdk.UnwrapSDKContext(ctx))
	return &types.QueryParamsResponse{
		Params: params,
	}, nil
}

func (k Keeper) GetCid(ctx context.Context, req *types.QueryGetCidRequest) (*types.QueryGetCidResponse, error) {
	cid, err := gocid.Decode(req.Cid)
	if err != nil {
		return nil, err
	}
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	record := k.GetCidBlock(sdkCtx, cid)
	return &types.QueryGetCidResponse{Data: record}, nil
}

func (k Keeper) HasCid(ctx context.Context, req *types.QueryHasCidRequest) (*types.QueryHasCidResponse, error) {
	cid, err := gocid.Decode(req.Cid)
	if err != nil {
		return nil, err
	}
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	exist := k.HasCidBlock(sdkCtx, cid)
	return &types.QueryHasCidResponse{Exist: exist}, nil
}

func (k Keeper) GetCidSize(ctx context.Context, req *types.QueryGetCidSizeRequest) (*types.QueryGetCidSizeResponse, error) {
	cid, err := gocid.Decode(req.Cid)
	if err != nil {
		return nil, err
	}
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	record := k.GetCidBlock(sdkCtx, cid)
	return &types.QueryGetCidSizeResponse{Size_: uint64(len(record))}, nil
}
