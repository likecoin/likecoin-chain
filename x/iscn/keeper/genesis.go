package keeper

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/likecoin/likechain/x/iscn/types"
)

func (k Keeper) InitGenesis(ctx sdk.Context, genesis *types.GenesisState) {
	k.SetParams(ctx, genesis.Params)
	for _, iscnRecord := range genesis.IscnRecords {
		iscnRecordMap := map[string]interface{}{}
		err := json.Unmarshal(iscnRecord, &iscnRecordMap)
		if err != nil {
			panic(err)
		}
		idStr := iscnRecordMap["@id"].(string)
		id, err := types.ParseIscnId(idStr)
		if err != nil {
			panic(err)
		}
		cid := types.ComputeDataCid(iscnRecord)
		seq := k.AddStoreRecord(ctx, StoreRecord{
			IscnId:   id,
			CidBytes: cid.Bytes(),
			Data:     iscnRecord,
		})
		fingerprints := iscnRecordMap["contentFingerprints"].([]string)
		for _, fingerprint := range fingerprints {
			k.AddFingerprintSequence(ctx, fingerprint, seq)
		}
	}
	for _, tracingIdRecord := range genesis.TracingIdRecords {
		iscnId, err := types.ParseIscnId(tracingIdRecord.IscnId)
		if err != nil {
			panic(err)
		}
		owner, err := sdk.AccAddressFromBech32(tracingIdRecord.Owner)
		if err != nil {
			panic(err)
		}
		k.SetTracingIdRecord(ctx, iscnId, &TracingIdRecord{
			OwnerAddressBytes: owner.Bytes(),
			LatestVersion:     tracingIdRecord.LatestVersion,
		})
	}
}

func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	params := k.GetParams(ctx)
	tracingIdRecords := []types.GenesisState_TracingIdRecord{}
	k.IterateTracingIdRecords(ctx, func(iscnId types.IscnId, tracingIdRecord types.TracingIdRecord) bool {
		iscnId.Version = 0
		tracingIdRecords = append(tracingIdRecords, types.GenesisState_TracingIdRecord{
			IscnId:        iscnId.String(),
			Owner:         tracingIdRecord.OwnerAddress().String(),
			LatestVersion: tracingIdRecord.LatestVersion,
		})
		return false
	})
	iscnRecords := []types.IscnInput{}
	k.IterateStoreRecords(ctx, func(_ uint64, record StoreRecord) bool {
		iscnRecords = append(iscnRecords, record.Data)
		return false
	})
	return types.NewGenesisState(params, tracingIdRecords, iscnRecords)
}
