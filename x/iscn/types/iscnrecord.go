package types

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"
)

func ValidateFingerprints(fingerprints []string) error {
	usedFingerprint := map[string]struct{}{}
	for _, fingerprint := range fingerprints {
		_, ok := usedFingerprint[fingerprint]
		if ok {
			return fmt.Errorf("repeated fingerprint entry")
		}
		usedFingerprint[fingerprint] = struct{}{}
		u, err := url.ParseRequestURI(fingerprint)
		if err != nil {
			return fmt.Errorf("invalid fingerprint URL %s: %w", fingerprint, err)
		}
		if u.Scheme == "" {
			return fmt.Errorf("empty fingerprint URL scheme in fingerprint %s", fingerprint)
		}
	}
	return nil
}

func (record *IscnRecord) Validate() error {
	err := ValidateFingerprints(record.ContentFingerprints)
	if err != nil {
		return fmt.Errorf("invalid content fingerprints: %w", err)
	}
	for _, stakeholder := range record.Stakeholders {
		err := stakeholder.Validate()
		if err != nil {
			return err
		}
	}
	err = record.ContentMetadata.Validate()
	if err != nil {
		return err
	}
	return nil
}

type IscnRecordJsonLdInfo struct {
	Id         IscnId
	Timestamp  time.Time
	ParentIpld *CID
}

func (record *IscnRecord) ToJsonLd(info *IscnRecordJsonLdInfo) ([]byte, error) {
	stakeholders := []interface{}{}
	for _, stakeholder := range record.Stakeholders {
		normalizedStakeholder, err := stakeholder.Normalize()
		if err != nil {
			return nil, err
		}
		stakeholders = append(stakeholders, normalizedStakeholder)
	}
	normalizedContentMetadata, err := record.ContentMetadata.Normalize()
	if err != nil {
		return nil, err
	}
	recordMap := map[string]interface{}{
		"@context": map[string]interface{}{
			"@vocab": "http://iscn.io/",
			"recordParentIPLD": map[string]interface{}{
				"@container": "@index",
			},
			"stakeholders": map[string]interface{}{
				"@context": map[string]interface{}{
					"@vocab":           "http://schema.org/",
					"entity":           "http://iscn.io/entity",
					"rewardProportion": "http://iscn.io/rewardProportion",
					"contributionType": "http://iscn.io/contributionType",
					"footprint":        "http://iscn.io/footprint",
				},
			},
			"contentMetadata": map[string]interface{}{
				"@context": nil,
			},
		},
		"@type":               "Record",
		"@id":                 info.Id.String(),
		"recordTimestamp":     info.Timestamp.UTC().Format("2006-01-02T15:04:05-07:00"),
		"recordVersion":       info.Id.Version,
		"recordNotes":         record.RecordNotes,
		"contentFingerprints": record.ContentFingerprints,
		"stakeholders":        stakeholders,
		"contentMetadata":     normalizedContentMetadata,
	}
	if info.ParentIpld != nil {
		recordMap["recordParentIPLD"] = map[string]string{
			"/": info.ParentIpld.String(),
		}
	}
	return json.Marshal(recordMap)
}
