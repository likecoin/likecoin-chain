package keeper

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/likecoin/likechain/x/iscn/types"
)

func (k Keeper) InitGenesis(ctx sdk.Context, genesis *types.GenesisState) {
	k.SetParams(ctx, genesis.Params)
	for _, entry := range genesis.RecordEntries {
		owner, err := sdk.AccAddressFromBech32(entry.Owner)
		if err != nil {
			panic(err)
		}
		for _, record := range entry.Records {
			recordMap := map[string]interface{}{}
			err = json.Unmarshal(record, &recordMap)
			if err != nil {
				panic(err)
			}
			idStr := recordMap["@id"].(string)
			id, err := types.ParseIscnID(idStr)
			if err != nil {
				panic(err)
			}
			cid := types.ComputeRecordCid(record)
			k.SetCidBlock(ctx, cid, record)
			k.SetCidIscnId(ctx, cid, id)
			k.SetIscnIdCid(ctx, id, cid)
			k.SetIscnIdVersion(ctx, id, id.Version)
			k.SetIscnIdOwner(ctx, id, owner)
			fingerprints := recordMap["contentFingerprints"].([]string)
			for _, fingerprint := range fingerprints {
				k.AddFingerPrintCid(ctx, fingerprint, cid)
			}
		}
	}
}

func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	entries := []types.GenesisIscnEntry{}
	usedIscnIds := map[string]struct{}{}
	k.IterateIscnIds(ctx, func(iscnId IscnId, cid CID) bool {
		iscnId.Version = 0
		iscnIdStr := iscnId.String()
		_, used := usedIscnIds[iscnIdStr]
		if used {
			return false
		}
		owner := k.GetIscnIdOwner(ctx, iscnId)
		entry := types.GenesisIscnEntry{
			IscnId: iscnId.String(),
			Owner:  owner.String(),
		}
		maxVersion := k.GetIscnIdVersion(ctx, iscnId)
		records := make([][]byte, maxVersion)
		for version := uint64(1); version <= maxVersion; version++ {
			iscnId.Version = version
			cid := k.GetIscnIdCid(ctx, iscnId)
			record := k.GetCidBlock(ctx, *cid)
			records = append(records, record)
		}
		entry.Records = records
		entries = append(entries, entry)
		return false
	})
	return types.NewGenesisState(k.GetParams(ctx), entries)
}
