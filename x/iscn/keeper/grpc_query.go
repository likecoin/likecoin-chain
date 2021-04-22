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

func (k Keeper) queryIscnRecordsByIscnId(ctx sdk.Context, iscnId IscnId) (*types.QueryIscnRecordsResponse, error) {
	tracingIdRecord := k.GetTracingIdRecord(ctx, iscnId)
	if tracingIdRecord == nil || iscnId.Version > tracingIdRecord.LatestVersion {
		return nil, sdkerrors.Wrapf(types.ErrRecordNotFound, "%s", iscnId.String())
	}
	if iscnId.Version == 0 {
		iscnId.Version = tracingIdRecord.LatestVersion
	}
	seq := k.GetIscnIdSequence(ctx, iscnId)
	storeRecord := k.GetStoreRecord(ctx, seq)
	records := []types.Record{{
		IscnId:              iscnId.String(),
		Owner:               tracingIdRecord.OwnerAddress().String(),
		Ipld:                storeRecord.Cid().String(),
		LatestRecordVersion: tracingIdRecord.LatestVersion,
		Record:              storeRecord.Data,
	}}
	return &types.QueryIscnRecordsResponse{Records: records}, nil
}

func (k Keeper) queryIscnRecordsByFingerprint(ctx sdk.Context, fingerprint string) (*types.QueryIscnRecordsResponse, error) {
	records := []types.Record{}
	// TODO: pagination?
	k.IterateFingerprintSequences(ctx, fingerprint, func(seq uint64) bool {
		storeRecord := k.GetStoreRecord(ctx, seq)
		tracingIdRecord := k.GetTracingIdRecord(ctx, storeRecord.IscnId)
		records = append(records, types.Record{
			IscnId:              storeRecord.String(),
			Owner:               tracingIdRecord.OwnerAddress().String(),
			Ipld:                storeRecord.Cid().String(),
			LatestRecordVersion: tracingIdRecord.LatestVersion,
			Record:              storeRecord.Data,
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
	seq := k.GetCidSequence(sdkCtx, cid)
	var data []byte
	if seq != 0 {
		storeRecord := k.GetStoreRecord(sdkCtx, seq)
		data = storeRecord.Data
	}
	return &types.QueryGetCidResponse{Data: data}, nil
}

func (k Keeper) HasCid(ctx context.Context, req *types.QueryHasCidRequest) (*types.QueryHasCidResponse, error) {
	cid, err := gocid.Decode(req.Cid)
	if err != nil {
		return nil, err
	}
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	seq := k.GetCidSequence(sdkCtx, cid)
	exist := seq != 0
	return &types.QueryHasCidResponse{Exist: exist}, nil
}

func (k Keeper) GetCidSize(ctx context.Context, req *types.QueryGetCidSizeRequest) (*types.QueryGetCidSizeResponse, error) {
	cid, err := gocid.Decode(req.Cid)
	if err != nil {
		return nil, err
	}
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	seq := k.GetCidSequence(sdkCtx, cid)
	if seq == 0 {
		return nil, fmt.Errorf("CID %s not found", cid.String())
	}
	size := uint64(0)
	storeRecord := k.GetStoreRecord(sdkCtx, seq)
	if storeRecord != nil {
		size = uint64(len(storeRecord.Data))
	}
	return &types.QueryGetCidSizeResponse{Size_: size}, nil
}
