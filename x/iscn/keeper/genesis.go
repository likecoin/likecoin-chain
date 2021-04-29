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
		normalizedRecord, err := iscnRecord.Normalize()
		if err != nil {
			panic(err)
		}
		cid := types.ComputeDataCid(normalizedRecord)
		seq := k.AddStoreRecord(ctx, StoreRecord{
			IscnId:   id,
			CidBytes: cid.Bytes(),
			Data:     IscnInput(normalizedRecord),
		})
		fingerprints := iscnRecordMap["contentFingerprints"].([]interface{})
		for _, fingerprint := range fingerprints {
			k.AddFingerprintSequence(ctx, fingerprint.(string), seq)
		}
	}
	for _, contentIdRecord := range genesis.ContentIdRecords {
		iscnId, err := types.ParseIscnId(contentIdRecord.IscnId)
		if err != nil {
			panic(err)
		}
		owner, err := sdk.AccAddressFromBech32(contentIdRecord.Owner)
		if err != nil {
			panic(err)
		}
		k.SetContentIdRecord(ctx, iscnId, &ContentIdRecord{
			OwnerAddressBytes: owner.Bytes(),
			LatestVersion:     contentIdRecord.LatestVersion,
		})
	}
}

func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	params := k.GetParams(ctx)
	contentIdRecords := []types.GenesisState_ContentIdRecord{}
	k.IterateContentIdRecords(ctx, func(iscnId types.IscnId, contentIdRecord types.ContentIdRecord) bool {
		contentIdRecords = append(contentIdRecords, types.GenesisState_ContentIdRecord{
			IscnId:        iscnId.Prefix.String(),
			Owner:         contentIdRecord.OwnerAddress().String(),
			LatestVersion: contentIdRecord.LatestVersion,
		})
		return false
	})
	iscnRecords := []types.IscnInput{}
	k.IterateStoreRecords(ctx, func(_ uint64, record StoreRecord) bool {
		iscnRecords = append(iscnRecords, IscnInput(record.Data))
		return false
	})
	return types.NewGenesisState(params, contentIdRecords, iscnRecords)
}
