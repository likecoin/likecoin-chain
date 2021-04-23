package keeper

import (
	context "context"
	"fmt"

	gocid "github.com/ipfs/go-cid"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/likecoin/likechain/x/iscn/types"
)

const maxLimit = 100

var _ types.QueryServer = Keeper{}

func (k Keeper) RecordsById(ctx context.Context, req *types.QueryRecordsByIdRequest) (*types.QueryRecordsByIdResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	fromVersion := req.FromVersion
	toVersion := req.ToVersion
	if toVersion != 0 && toVersion < fromVersion {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid version range")
	}
	iscnId, err := types.ParseIscnId(req.IscnId)
	if err != nil {
		return nil, sdkerrors.Wrapf(types.ErrInvalidIscnId, "%s", err.Error())
	}
	if iscnId.Version != 0 {
		fromVersion = iscnId.Version
		toVersion = iscnId.Version
	}
	contentIdRecord := k.GetContentIdRecord(sdkCtx, iscnId)
	if contentIdRecord == nil {
		return nil, sdkerrors.Wrapf(types.ErrRecordNotFound, "%s", iscnId.String())
	}
	latestVersion := contentIdRecord.LatestVersion
	if fromVersion == 0 {
		fromVersion = latestVersion
	}
	if toVersion == 0 {
		toVersion = latestVersion
	}
	if fromVersion > latestVersion || toVersion > latestVersion {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "version range exceed current maximum version")
	}
	records := make([]types.QueryResponseRecord, 0, toVersion-fromVersion+1)
	for version := fromVersion; version <= toVersion; version++ {
		iscnId.Version = version
		seq := k.GetIscnIdSequence(sdkCtx, iscnId)
		storeRecord := k.GetStoreRecord(sdkCtx, seq)
		records = append(records, types.QueryResponseRecord{
			Ipld: storeRecord.Cid().String(),
			Data: storeRecord.Data,
		})
	}
	return &types.QueryRecordsByIdResponse{
		Owner:         contentIdRecord.OwnerAddress().String(),
		LatestVersion: latestVersion,
		Records:       records,
	}, nil
}

func (k Keeper) RecordsByFingerprint(ctx context.Context, req *types.QueryRecordsByFingerprintRequest) (*types.QueryRecordsByFingerprintResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	records := []types.QueryResponseRecord{}
	nextSeq := uint64(0)
	count := 0
	k.IterateFingerprintSequencesWithStartingSequence(sdkCtx, req.Fingerprint, req.FromSequence, func(seq uint64) bool {
		if count >= maxLimit {
			nextSeq = seq
			return true
		}
		count++
		storeRecord := k.GetStoreRecord(sdkCtx, seq)
		records = append(records, types.QueryResponseRecord{
			Ipld: storeRecord.Cid().String(),
			Data: storeRecord.Data,
		})
		return false
	})
	return &types.QueryRecordsByFingerprintResponse{
		Records:      records,
		NextSequence: nextSeq,
	}, nil
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
