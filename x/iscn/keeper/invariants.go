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
	ir.RegisterRoute(types.ModuleName, "iscn-records", IscnRecordsInvariant(k))
	ir.RegisterRoute(types.ModuleName, "iscn-fingerprints", IscnFingerprintsInvariant(k))
}

func IscnRecordsInvariant(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		problemCount := uint64(0)
		msgBuf := strings.Builder{}

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
			parentSeq := k.GetCidSequence(ctx, cid)
			if parentSeq == 0 {
				logProblem(fmt.Sprintf("no record for parent CID %s for ISCN ID %s", cid.String(), id.String()))
				return
			}
			parentStoreRecord := k.GetStoreRecord(ctx, parentSeq)
			if parentStoreRecord == nil {
				logProblem(fmt.Sprintf("no store record for parent CID %s for ISCN ID %s", cid.String(), id.String()))
				return
			}
			parentIscnId := parentStoreRecord.IscnId
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

		// 1. check all records are valid
		// 2. check every conntent ID record has the corresponding ISCN ID records
		k.IterateContentIdRecords(ctx, func(iscnPrefixId IscnId, contentIdRecord ContentIdRecord) bool {
			if contentIdRecord.LatestVersion == 0 {
				logProblem(fmt.Sprintf("content ID %s has 0 as latest version record", contentIdRecord.String()))
				return false
			}
			for version := uint64(1); version <= contentIdRecord.LatestVersion; version++ {
				id := iscnPrefixId
				id.Version = version
				idStr := id.String()
				seq := k.GetIscnIdSequence(ctx, id)
				if seq == 0 {
					logProblem(fmt.Sprintf("ISCN ID %s has latest version %d, but sequence returns 0", idStr, contentIdRecord.LatestVersion))
					continue
				}
				storeRecord := k.GetStoreRecord(ctx, seq)
				if storeRecord == nil {
					logProblem(fmt.Sprintf("ISCN ID %s has sequence record %d, but store record not found", idStr, seq))
					continue
				}
				if !storeRecord.IscnId.Equal(&id) {
					logProblem(fmt.Sprintf("ISCN ID %s has sequence record %d, but store record for sequence %d has another ISCN ID %s", idStr, seq, seq, storeRecord.IscnId.String()))
				}
				cid := storeRecord.Cid()
				cidStr := cid.String()
				computedCid := types.ComputeDataCid(storeRecord.Data)
				if !cid.Equals(computedCid) {
					logProblem(fmt.Sprintf("ISCN ID has CID record %s, but the computed CID for the record is %s", cidStr, computedCid.String()))
				}
				checkRecord(id, storeRecord.Data)
			}
			return false
		})

		// 3. check all ISCN ID has content ID record
		// 4. check all ISCN ID and CID can reverse lookup sequence
		// 5. check contiguous sequence
		prevSeq := uint64(0)
		k.IterateStoreRecords(ctx, func(seq uint64, storeRecord StoreRecord) bool {
			if seq != prevSeq+1 {
				logProblem(fmt.Sprintf("discontiguous sequence (%d to %d)", prevSeq, seq))
			}
			prevSeq = seq
			contentIdRecord := k.GetContentIdRecord(ctx, storeRecord.IscnId)
			if contentIdRecord == nil {
				logProblem(fmt.Sprintf("store record sequence %d has ISCN ID %s, but the content ID record does not exist", seq, storeRecord.IscnId.String()))
			} else if contentIdRecord.LatestVersion < storeRecord.IscnId.Version {
				logProblem(fmt.Sprintf("ISCN ID %s has content ID record with smaller latest version %d", storeRecord.IscnId.String(), contentIdRecord.LatestVersion))
			}
			iscnIdSeq := k.GetIscnIdSequence(ctx, storeRecord.IscnId)
			if iscnIdSeq != seq {
				logProblem(fmt.Sprintf("store record sequence %d has ISCN ID %s, but reverse lookup record points to sequence %d", seq, storeRecord.IscnId.String(), iscnIdSeq))
			}
			cidSeq := k.GetCidSequence(ctx, storeRecord.Cid())
			if cidSeq != seq {
				logProblem(fmt.Sprintf("store record sequence %d has CID %s, but reverse lookup record points to sequence %d", seq, storeRecord.Cid().String(), cidSeq))
			}
			return false
		})
		seqCount := k.GetSequenceCount(ctx)
		if prevSeq != seqCount {
			logProblem(fmt.Sprintf("max sequence (%d) does not equal to sequence count (%d)", prevSeq, seqCount))
		}

		// 5. check all ISCN ID and CID reverse lookup sequence actually exist
		cidIter := k.prefixStore(ctx, CidToSequencePrefix).Iterator(nil, nil)
		defer cidIter.Close()
		for ; cidIter.Valid(); cidIter.Next() {
			seq := types.DecodeUint64(cidIter.Value())
			if seq == 0 || seq > seqCount {
				cid := types.MustCidFromBytes(cidIter.Key())
				logProblem(fmt.Sprintf("CID %s has CID-sequence reverse lookup (sequence %d)", cid.String(), seq))
			}
		}

		iscnIdIter := k.prefixStore(ctx, IscnIdToSequencePrefix).Iterator(nil, nil)
		defer iscnIdIter.Close()
		for ; iscnIdIter.Valid(); iscnIdIter.Next() {
			seq := types.DecodeUint64(iscnIdIter.Value())
			if seq == 0 || seq > seqCount {
				iscnId := k.MustUnmarshalIscnId(iscnIdIter.Key())
				logProblem(fmt.Sprintf("ISCN ID %s has ISCN-ID-sequence reverse lookup (sequence %d)", iscnId.String(), seq))
			}
		}

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
		k.IterateAllFingerprints(ctx, func(fingerprint string, seq uint64) bool {
			storeRecord := k.GetStoreRecord(ctx, seq)
			if storeRecord == nil {
				logProblem(fmt.Sprintf("fingerprint %s has sequence record %d, but store record does not exist", fingerprint, seq))
				return false
			}
			recordMap := map[string]interface{}{}
			err := json.Unmarshal(storeRecord.Data, &recordMap)
			if err != nil {
				logProblem(fmt.Sprintf("cannot unmarshal record for fingerprint %s (sequence %d) as JSON", fingerprint, seq))
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
		k.IterateStoreRecords(ctx, func(seq uint64, storeRecord StoreRecord) bool {
			recordMap := map[string]interface{}{}
			err := json.Unmarshal(storeRecord.Data, &recordMap)
			if err != nil {
				logProblem(fmt.Sprintf("cannot unmarshal record for store record sequence %d as JSON", seq))
				return false
			}
			field, ok := recordMap["contentFingerprints"]
			if !ok {
				logProblem(fmt.Sprintf("record for store record sequence %d has no contentFingerprints field", seq))
				return false
			}
			arr, ok := field.([]interface{})
			if !ok {
				logProblem(fmt.Sprintf("record for store record sequence %d has wrong type for contentFingerprints field", seq))
				return false
			}
			for _, v := range arr {
				fingerprint, ok := v.(string)
				if !ok {
					logProblem(fmt.Sprintf("record for store record sequence %d has value with wrong type in contentFingerprints field", seq))
					return false
				}
				if !k.HasFingerprintSequence(ctx, fingerprint, seq) {
					logProblem(fmt.Sprintf("dangling fingerprint value %s in sequence %d", fingerprint, seq))
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
