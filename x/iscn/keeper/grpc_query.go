package keeper

import (
	context "context"

	gocid "github.com/ipfs/go-cid"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/likecoin/likechain/x/iscn/types"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) queryIscnRecordsByIscnId(ctx sdk.Context, iscnId IscnId) (*types.QueryIscnRecordsResponse, error) {
	latestVersion := k.GetIscnIdVersion(ctx, iscnId)
	if latestVersion == 0 || iscnId.Version > latestVersion {
		return nil, sdkerrors.Wrapf(types.ErrRecordNotFound, "%s", iscnId.String())
	}
	if iscnId.Version == 0 {
		iscnId.Version = latestVersion
	}
	owner := k.GetIscnIdOwner(ctx, iscnId)
	cid := k.GetIscnIdCid(ctx, iscnId)
	record := k.GetCidBlock(ctx, *cid)
	records := []types.Record{{
		IscnId:              iscnId.String(),
		Owner:               owner.String(),
		Ipld:                cid.String(),
		LatestRecordVersion: latestVersion,
		Record:              types.IscnInput(record),
	}}
	return &types.QueryIscnRecordsResponse{Records: records}, nil
}

func (k Keeper) queryIscnRecordsByFingerprint(ctx sdk.Context, fingerprint string) (*types.QueryIscnRecordsResponse, error) {
	records := []types.Record{}
	// TODO: pagination?
	k.IterateFingerprintCids(ctx, fingerprint, func(cid CID) bool {
		iscnId := *k.GetCidIscnId(ctx, cid)
		owner := k.GetIscnIdOwner(ctx, iscnId)
		latestVersion := k.GetIscnIdVersion(ctx, iscnId)
		record := k.GetCidBlock(ctx, cid)
		records = append(records, types.Record{
			IscnId:              iscnId.String(),
			Owner:               owner.String(),
			Ipld:                cid.String(),
			LatestRecordVersion: latestVersion,
			Record:              types.IscnInput(record),
		})
		return false
	})
	return &types.QueryIscnRecordsResponse{Records: records}, nil
}

func (k Keeper) IscnRecords(ctx context.Context, req *types.QueryIscnRecordsRequest) (*types.QueryIscnRecordsResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	if len(req.IscnId) > 0 {
		if len(req.Fingerprint) > 0 {
			return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "only one of iscn_id and fingerprint can exist in query parameters")
		}
		iscnId, err := types.ParseIscnId(req.IscnId)
		if err != nil {
			return nil, sdkerrors.Wrapf(types.ErrInvalidIscnId, "%s", err.Error())
		}
		return k.queryIscnRecordsByIscnId(sdkCtx, iscnId)
	} else if len(req.Fingerprint) > 0 {
		return k.queryIscnRecordsByFingerprint(sdkCtx, req.Fingerprint)
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
