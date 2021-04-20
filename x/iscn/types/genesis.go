package types

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (genesis GenesisState) Validate() error {
	err := genesis.Params.Validate()
	if err != nil {
		return err
	}
	usedIscnIds := map[string]struct{}{}
	usedIplds := map[string]struct{}{}
	for _, entry := range genesis.RecordEntries {
		_, err := sdk.AccAddressFromBech32(entry.Owner)
		if err != nil {
			return err
		}
		_, ok := usedIscnIds[entry.IscnId]
		if ok {
			return fmt.Errorf("repeated ISCN ID: %s", entry.IscnId)
		}
		usedIscnIds[entry.IscnId] = struct{}{}
		id, err := ParseIscnId(entry.IscnId)
		if err != nil {
			return err
		}
		if id.Version != 0 {
			return fmt.Errorf("invalid ISCN ID version: %s", entry.IscnId)
		}
		for i, record := range entry.Records {
			recordMap := map[string]interface{}{}
			err = json.Unmarshal(record, &recordMap)
			if err != nil {
				return err
			}
			idAny, ok := recordMap["@id"]
			if !ok {
				return fmt.Errorf("entry has no \"@id\" field")
			}
			idStr, ok := idAny.(string)
			if !ok {
				return fmt.Errorf("invalid \"@id\" field type")
			}
			recordId, err := ParseIscnId(idStr)
			if err != nil {
				return err
			}
			expectedVersion := i + 1
			if recordId.Version != uint64(expectedVersion) {
				return fmt.Errorf("invalid ISCN ID version")
			}
			cid := ComputeRecordCid(record)
			_, cidExist := usedIplds[cid.String()]
			if cidExist {
				return fmt.Errorf("CID repeated: %s", cid.String())
			}
			fingerprintsAny, ok := recordMap["contentFingerprints"]
			if !ok {
				return fmt.Errorf("entry has no \"contentFingerprints\" field")
			}
			fingerprints, ok := fingerprintsAny.([]string)
			if !ok {
				return fmt.Errorf("invalid \"contentFingerprints\" field type")
			}
			for _, fingerprint := range fingerprints {
				u, err := url.ParseRequestURI(fingerprint)
				if err != nil {
					return err
				}
				if u.Scheme == "" {
					return fmt.Errorf("empty fingerprint URL scheme")
				}
			}
		}
	}
	return nil
}

func NewGenesisState(params Params, recordEntries []GenesisIscnEntry) *GenesisState {
	return &GenesisState{
		Params:        params,
		RecordEntries: recordEntries,
	}
}

func DefaultGenesisState() *GenesisState {
	return NewGenesisState(DefaultParams(), []GenesisIscnEntry{})
}

func GetGenesisStateFromAppState(cdc codec.JSONMarshaler, appState map[string]json.RawMessage) *GenesisState {
	var genesisState GenesisState

	if appState[ModuleName] != nil {
		cdc.MustUnmarshalJSON(appState[ModuleName], &genesisState)
	}

	return &genesisState
}
