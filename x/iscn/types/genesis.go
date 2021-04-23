package types

import (
	"encoding/json"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (genesis GenesisState) Validate() error {
	err := genesis.Params.Validate()
	if err != nil {
		return fmt.Errorf("invalid ISCN parameters: %w", err)
	}
	iscnVersionMap := map[string]uint64{}
	for i, record := range genesis.IscnRecords {
		recordMap := map[string]interface{}{}
		err := json.Unmarshal(record, &recordMap)
		if err != nil {
			return fmt.Errorf("cannot unmarshal record at index %d as JSON: %v", i, err)
		}
		idAny, ok := recordMap["@id"]
		if !ok {
			return fmt.Errorf("record at index %d has no \"@id\" field", i)
		}
		idStr, ok := idAny.(string)
		if !ok {
			return fmt.Errorf("record at index %d has invalid \"@id\" type", i)
		}
		iscnId, err := ParseIscnId(idStr)
		if err != nil {
			return fmt.Errorf("record at index %d has invalid ISCN ID: %w", i, err)
		}
		if iscnId.Version == 0 {
			return fmt.Errorf("record at index %d has 0 as ISCN ID version", i)
		}
		iscnPrefix := iscnId.Prefix.String()
		prevVersion := iscnVersionMap[iscnPrefix]
		if iscnId.Version != prevVersion+1 {
			return fmt.Errorf("record at index %d (ISCN ID %s) has non-contiguous version (previous version %d, current version %d)", i, iscnId.String(), prevVersion, iscnId.Version)
		}
		iscnVersionMap[iscnPrefix] = iscnId.Version
		// not checking repeated CID, since CID bases from the hash of content, CID repeated -> hash repeated -> content repeated -> "@id" field repeated -> invalid version
		fingerprintsAny, ok := recordMap["contentFingerprints"]
		if !ok {
			return fmt.Errorf("record at index %d (ISCN ID %s) has no \"contentFingerprints\" field", i, iscnId.String())
		}
		fingerprints, ok := fingerprintsAny.([]string)
		if !ok {
			return fmt.Errorf("record at index %d (ISCN ID %s) has invalid \"contentFingerprints\" field type", i, iscnId.String())
		}
		err = ValidateFingerprints(fingerprints)
		if err != nil {
			return fmt.Errorf("record at index %d (ISCN ID %s) has invalid \"contentFingerprints\" entries: %w", i, iscnId.String(), err)
		}
	}
	for _, contentIdRecord := range genesis.ContentIdRecords {
		_, err := sdk.AccAddressFromBech32(contentIdRecord.Owner)
		if err != nil {
			return fmt.Errorf("invalid owner address %s in content ID record entries: %w", contentIdRecord.Owner, err)
		}
		iscnId, err := ParseIscnId(contentIdRecord.IscnId)
		if err != nil {
			return fmt.Errorf("cannot parse ISCN ID %s in content ID record entries: %w", contentIdRecord.IscnId, err)
		}
		if iscnId.Version != 0 {
			return fmt.Errorf("invalid version in ISCN ID %s in content ID record entries, expect version 0", iscnId.String())
		}
		idPrefixStr := iscnId.Prefix.String()
		if iscnVersionMap[idPrefixStr] != contentIdRecord.LatestVersion {
			return fmt.Errorf("ISCN ID prefix %s latest version does not match the content ID record entry", iscnId.String())
		}
		delete(iscnVersionMap, idPrefixStr)
		iscnVersionMap[idPrefixStr] = 0
	}
	for prefixStr := range iscnVersionMap {
		return fmt.Errorf("ISCN ID prefix %s has related ISCN record but no content ID record", prefixStr)
	}
	return nil
}

func NewGenesisState(params Params, contentIdRecords []GenesisState_ContentIdRecord, iscnRecords []IscnInput) *GenesisState {
	return &GenesisState{
		Params:           params,
		ContentIdRecords: contentIdRecords,
		IscnRecords:      iscnRecords,
	}
}

func DefaultGenesisState() *GenesisState {
	return NewGenesisState(DefaultParams(), nil, nil)
}

func GetGenesisStateFromAppState(cdc codec.JSONMarshaler, appState map[string]json.RawMessage) *GenesisState {
	var genesisState GenesisState

	if appState[ModuleName] != nil {
		cdc.MustUnmarshalJSON(appState[ModuleName], &genesisState)
	}

	return &genesisState
}
