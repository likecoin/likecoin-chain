package keeper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	gocid "github.com/ipfs/go-cid"

	"github.com/likecoin/likechain/x/iscn/types"
)

func RegisterInvariants(ir sdk.InvariantRegistry, k Keeper) {
	ir.RegisterRoute(types.ModuleName, "iscn-records", IscnRecordsInvariant(k))
	ir.RegisterRoute(types.ModuleName, "iscn-fingerprints", IscnFingerprintsInvariant(k))
}

func IscnRecordsInvariant(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		problemCount := uint64(0)
		msgBuf := strings.Builder{}
		seenId := map[string]string{}
		cidUsedByIscnId := map[string]string{}

		logProblem := func(msg string) {
			msgBuf.WriteString(" - ")
			msgBuf.WriteString(msg)
			msgBuf.WriteString("\n")
			problemCount++
		}

		checkRecordId := func(id IscnId, recordMap map[string]interface{}) {
			field, ok := recordMap["@id"]
			if !ok {
				logProblem(fmt.Sprintf("record for ISCN ID %s has no @id field", id.String()))
				return
			}
			s, ok := field.(string)
			if !ok {
				logProblem(fmt.Sprintf("record for ISCN ID %s has wrong type for @id field", id.String()))
				return
			}
			idStr := id.String()
			if s != idStr {
				logProblem(fmt.Sprintf("record for ISCN ID %s has @id field %s, not equal to the ISCN ID", idStr, s))
				return
			}
		}

		checkRecordVersion := func(id IscnId, recordMap map[string]interface{}) {
			field, ok := recordMap["recordVersion"]
			if !ok {
				logProblem(fmt.Sprintf("record for ISCN ID %s has no recordVersion field", id.String()))
				return
			}
			num, ok := field.(json.Number)
			if !ok {
				logProblem(fmt.Sprintf("record for ISCN ID %s has wrong type for recordVersion field", id.String()))
				return
			}
			version, err := num.Int64()
			if err != nil {
				logProblem(fmt.Sprintf("record for ISCN ID %s has non-integer recordVersion field", id.String()))
				return
			}
			if version < 1 {
				logProblem(fmt.Sprintf("record for ISCN ID %s has invalid value (%d) for recordVersion field", id.String(), version))
				return
			}
			if uint64(version) != id.Version {
				logProblem(fmt.Sprintf("record for ISCN ID %s has recordVersion field %d, not equal to the ISCN ID version", id.String(), version))
				return
			}
		}

		checkRecordParent := func(id IscnId, recordMap map[string]interface{}) {
			field, ok := recordMap["recordParentIPLD"]
			if id.Version == 1 {
				if ok {
					logProblem(fmt.Sprintf("record for ISCN ID %s has recordParentIPLD field, which should not exist for version 1", id.String()))
				}
				return
			}
			if !ok {
				logProblem(fmt.Sprintf("record for ISCN ID %s has no recordParentIPLD field", id.String()))
				return
			}
			fieldMap, ok := field.(map[string]interface{})
			if !ok {
				logProblem(fmt.Sprintf("record for ISCN ID %s has wrong type for recordParentIPLD field", id.String()))
				return
			}
			subField, ok := fieldMap["/"]
			if !ok {
				logProblem(fmt.Sprintf("record for ISCN ID %s has no recordParentIPLD.\"/\" sub-field", id.String()))
				return
			}
			ipldStr, ok := subField.(string)
			if !ok {
				logProblem(fmt.Sprintf("record for ISCN ID %s has wrong type for reocrdParentIPLD.\"/\" sub-field", id.String()))
				return
			}
			cid, err := gocid.Decode(ipldStr)
			if err != nil {
				logProblem(fmt.Sprintf("record for ISCN ID %s has invalid CID for reocrdParentIPLD field", id.String()))
				return
			}
			parentIscnId := k.GetCidIscnId(ctx, cid)
			if parentIscnId == nil {
				logProblem(fmt.Sprintf("no record for parent CID %s for ISCN ID %s", cid.String(), id.String()))
				return
			}
			if !parentIscnId.PrefixEqual(&id) || parentIscnId.Version != id.Version-1 {
				logProblem(fmt.Sprintf("record for parent CID %s for ISCN ID %s has ISCN ID %s, which is not the ID of the parent version", cid.String(), id.String(), parentIscnId.String()))
				return
			}
		}

		checkRecordTimestamp := func(id IscnId, recordMap map[string]interface{}) {
			field, ok := recordMap["recordTimestamp"]
			if !ok {
				logProblem(fmt.Sprintf("record for ISCN ID %s has no recordTimestamp field", id.String()))
				return
			}
			s, ok := field.(string)
			if !ok {
				logProblem(fmt.Sprintf("record for ISCN ID %s has wrong type for recordTimestamp field", id.String()))
				return
			}
			t, err := time.Parse("2006-01-02T15:04:05-07:00", s)
			if err != nil {
				logProblem(fmt.Sprintf("cannot parse field recordTimestamp as time for ISCN ID %s", id.String()))
				return
			}
			zoneName, offset := t.Zone()
			if offset != 0 {
				logProblem(fmt.Sprintf("record for ISCN ID %s has recordTimestamp with non-UTC timezone (%s, %d)", id.String(), zoneName, offset))
				return
			}
		}

		checkRecord := func(id IscnId, record []byte) {
			recordMap := map[string]interface{}{}
			decoder := json.NewDecoder(bytes.NewReader(record))
			decoder.UseNumber()
			err := decoder.Decode(&recordMap)
			if err != nil {
				logProblem(fmt.Sprintf("cannot unmarshal record for %s as JSON", id.String()))
				return
			}
			checkRecordId(id, recordMap)
			checkRecordVersion(id, recordMap)
			checkRecordParent(id, recordMap)
			checkRecordTimestamp(id, recordMap)
		}

		k.IterateIscnIds(ctx, func(id IscnId, cid CID) bool {
			_, ok := seenId[id.String()]
			if ok {
				return false
			}
			latestVersion := k.GetIscnIdVersion(ctx, id)
			if latestVersion == 0 {
				logProblem(fmt.Sprintf("ISCN ID %s has no latest version record", id.String()))
				return false
			}
			for version := uint64(1); version <= latestVersion; version++ {
				id.Version = version
				idStr := id.String()
				cid := k.GetIscnIdCid(ctx, id)
				if cid == nil {
					logProblem(fmt.Sprintf("ISCN ID %s has latest version %d, but CID record does not exist", idStr, latestVersion))
					continue
				}
				cidStr := cid.String()
				seenId[idStr] = cidStr
				cidUsedByIscnId[cidStr] = idStr
				cidReverseIscnId := k.GetCidIscnId(ctx, *cid)
				if cidReverseIscnId == nil || !cidReverseIscnId.Equal(&id) {
					logProblem(fmt.Sprintf("ISCN ID %s has CID record %s, but CID record of %s points to %s", idStr, cidStr, cidStr, cidReverseIscnId.String()))
				}
				record := k.GetCidBlock(ctx, *cid)
				if record == nil {
					logProblem(fmt.Sprintf("ISCN ID %s has CID record %s, but CID block does not exist", idStr, cidStr))
					continue
				}
				computedCid := types.ComputeRecordCid(record)
				if !cid.Equals(computedCid) {
					logProblem(fmt.Sprintf("ISCN ID has CID record %s, but the computed CID for the record is %s", cidStr, computedCid.String()))
					// we can still check other fields, so go on instead of skipping the remaining parts of the loop
				}
				checkRecord(id, record)
			}
			return false
		})

		k.IterateCidBlocks(ctx, func(cid CID, bz []byte) bool {
			cidStr := cid.String()
			_, ok := cidUsedByIscnId[cidStr]
			if ok {
				return false
			}
			logProblem(fmt.Sprintf("dangling CID record %s", cidStr))
			computedCid := types.ComputeRecordCid(bz)
			if !cid.Equals(computedCid) {
				logProblem(fmt.Sprintf("dangling CID record %s has computed CID %s", cidStr, computedCid.String()))
			}
			return false
		})

		broken := problemCount > 0
		msg := sdk.FormatInvariant(
			types.ModuleName, "iscn-records",
			fmt.Sprintf("Total number of problems found: %d\n%s", problemCount, msgBuf.String()),
		)
		// TODO: better logging
		fmt.Println(msg)
		return msg, broken
	}
}

func IscnFingerprintsInvariant(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		problemCount := uint64(0)
		msgBuf := strings.Builder{}

		logProblem := func(msg string) {
			msgBuf.WriteString(" - ")
			msgBuf.WriteString(msg)
			msgBuf.WriteString("\n")
			problemCount++
		}

		// 1. to check each fingerprint record actually points to a record with that fingerprint
		k.IterateFingerprints(ctx, func(fingerprint string, cid CID) bool {
			record := k.GetCidBlock(ctx, cid)
			if record == nil {
				logProblem(fmt.Sprintf("fingerprint %s has CID record %s, but CID block does not exist", fingerprint, cid.String()))
				return false
			}
			recordMap := map[string]interface{}{}
			err := json.Unmarshal(record, &recordMap)
			if err != nil {
				logProblem(fmt.Sprintf("cannot unmarshal record for fingerprint %s (CID %s) as JSON", fingerprint, cid.String()))
				return false
			}
			field, ok := recordMap["contentFingerprints"]
			if !ok {
				logProblem(fmt.Sprintf("record for fingerprint %s has no contentFingerprints field", fingerprint))
				return false
			}
			arr, ok := field.([]interface{})
			if !ok {
				logProblem(fmt.Sprintf("record for fingerprint %s has wrong type for contentFingerprints field", fingerprint))
				return false
			}
			found := false
			for _, v := range arr {
				recordFingerprint, ok := v.(string)
				if !ok {
					logProblem(fmt.Sprintf("record for fingerprint %s has value with wrong type in contentFingerprints field", fingerprint))
					return false
				}
				if recordFingerprint == fingerprint {
					found = true
				}
			}
			if !found {
				logProblem(fmt.Sprintf("record for fingerprint %s has no fingerprint value in contentFingerprints field matching the fingerprint", fingerprint))
				return false
			}
			return false
		})

		// 2. to check each record actually has a fingerprint record points to that record
		k.IterateCidBlocks(ctx, func(cid CID, bz []byte) bool {
			record := k.GetCidBlock(ctx, cid)
			if record == nil {
				logProblem(fmt.Sprintf("block for CID %s does not exist", cid.String()))
				return false
			}
			recordMap := map[string]interface{}{}
			err := json.Unmarshal(record, &recordMap)
			if err != nil {
				logProblem(fmt.Sprintf("cannot unmarshal record for CID %s as JSON", cid.String()))
				return false
			}
			field, ok := recordMap["contentFingerprints"]
			if !ok {
				logProblem(fmt.Sprintf("record for CID %s has no contentFingerprints field", cid.String()))
				return false
			}
			arr, ok := field.([]interface{})
			if !ok {
				logProblem(fmt.Sprintf("record for CID %s has wrong type for contentFingerprints field", cid.String()))
				return false
			}
			for _, v := range arr {
				fingerprint, ok := v.(string)
				if !ok {
					logProblem(fmt.Sprintf("record for CID %s has value with wrong type in contentFingerprints field", cid.String()))
					return false
				}
				if !k.HasFingerprintCid(ctx, fingerprint, cid) {
					logProblem(fmt.Sprintf("dangling fingerprint value %s in CID %s", fingerprint, cid.String()))
					return false
				}
			}
			return false
		})

		broken := problemCount > 0
		msg := sdk.FormatInvariant(
			types.ModuleName, "iscn-fingerprints",
			fmt.Sprintf("Total number of problems found: %d\n%s", problemCount, msgBuf.String()),
		)
		// TODO: better logging
		fmt.Println(msg)
		return msg, broken
	}
}
