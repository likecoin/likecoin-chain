package types

import (
	"encoding/json"
	"errors"
	"net/url"
	"time"
)

func (record *IscnRecord) Validate() error {
	for _, fingerprint := range record.ContentFingerprints {
		u, err := url.ParseRequestURI(fingerprint)
		if err != nil {
			return err
		}
		if u.Scheme == "" {
			return errors.New("empty fingerprint URL scheme")
		}
	}
	for _, stakeholder := range record.Stakeholders {
		err := stakeholder.Validate()
		if err != nil {
			return err
		}
	}
	err := record.ContentMetadata.Validate()
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
		recordMap["recordParentIpld"] = map[string]string{
			"/": info.ParentIpld.String(),
		}
	}
	return json.Marshal(recordMap)
}
