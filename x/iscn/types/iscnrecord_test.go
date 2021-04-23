package types

import (
	"testing"
	"time"

	gocid "github.com/ipfs/go-cid"
	"github.com/stretchr/testify/require"
)

func TestIscnRecordValidateFingerprints(t *testing.T) {
	var fingerprints []string
	var err error

	goodFingerprint1 := "hash://sha256/9564b85669d5e96ac969dd0161b8475bbced9e5999c6ec598da718a3045d6f2e"
	goodFingerprint2 := "ipfs://QmNrgEMcUygbKzZeZgYFosdd27VE9KnWbyUD73bKZJ3bGi"

	err = ValidateFingerprints(fingerprints) // empty fingerprints
	require.NoError(t, err)

	fingerprints = []string{goodFingerprint1, goodFingerprint2}
	err = ValidateFingerprints(fingerprints)
	require.NoError(t, err)

	fingerprints = []string{goodFingerprint1, goodFingerprint1}
	err = ValidateFingerprints(fingerprints)
	require.Error(t, err, "should not accept repeating fingerprints")

	fingerprints = []string{"://sha256/9564b85669d5e96ac969dd0161b8475bbced9e5999c6ec598da718a3045d6f2e"}
	err = ValidateFingerprints(fingerprints)
	require.Error(t, err, "should not accept fingerprint without scheme")

	fingerprints = []string{"sha256/9564b85669d5e96ac969dd0161b8475bbced9e5999c6ec598da718a3045d6f2e"}
	err = ValidateFingerprints(fingerprints)
	require.Error(t, err, "should not accept non-URI fingerprint")

	fingerprints = []string{"/9564b85669d5e96ac969dd0161b8475bbced9e5999c6ec598da718a3045d6f2e"}
	err = ValidateFingerprints(fingerprints)
	require.Error(t, err, "should not accept non-URI fingerprint")

	fingerprints = []string{"9564b85669d5e96ac969dd0161b8475bbced9e5999c6ec598da718a3045d6f2e"}
	err = ValidateFingerprints(fingerprints)
	require.Error(t, err, "should not accept non-URI fingerprint")
}
func TestIscnRecordValidate(t *testing.T) {
	goodIscnInput := IscnInput(`""`)
	badIscnInput := IscnInput(`"`)
	goodFingerprint := "hash://sha256/9564b85669d5e96ac969dd0161b8475bbced9e5999c6ec598da718a3045d6f2e"
	badFingerprint := "9564b85669d5e96ac969dd0161b8475bbced9e5999c6ec598da718a3045d6f2e"

	goodRecord := func() IscnRecord {
		return IscnRecord{
			RecordNotes:         "testing",
			ContentFingerprints: []string{goodFingerprint},
			Stakeholders:        []IscnInput{goodIscnInput},
			ContentMetadata:     goodIscnInput,
		}
	}

	var record IscnRecord
	var err error

	record = goodRecord()
	err = record.Validate()
	require.NoError(t, err)

	record = goodRecord()
	record.RecordNotes = ""
	err = record.Validate()
	require.NoError(t, err, "should accept record with no record notes")

	record = goodRecord()
	record.Stakeholders = nil
	err = record.Validate()
	require.NoError(t, err, "should accept record with no stakeholders")

	record = goodRecord()
	record.ContentFingerprints = nil
	err = record.Validate()
	require.NoError(t, err, "should accept record with no fingerprints")

	record = goodRecord()
	record.Stakeholders[0] = badIscnInput
	err = record.Validate()
	require.Error(t, err, "should not accept record with invalid IscnInput as stakeholders")

	record = goodRecord()
	record.Stakeholders[0] = nil
	err = record.Validate()
	require.Error(t, err, "should not accept record with nil IscnInput as stakeholders")

	record = goodRecord()
	record.ContentMetadata = badIscnInput
	err = record.Validate()
	require.Error(t, err, "should not accept record with invalid IscnInput as contentMetadata")

	record = goodRecord()
	record.ContentMetadata = nil
	err = record.Validate()
	require.Error(t, err, "should not accept record with nil IscnInput as contentMetadata")

	record = goodRecord()
	record.ContentFingerprints[0] = badFingerprint
	err = record.Validate()
	require.Error(t, err, "should not accept record with invalid fingerprints")
}

func TestIscnRecordToJsonLd(t *testing.T) {
	goodRecord := func() IscnRecord {
		return IscnRecord{
			RecordNotes:         "testing",
			ContentFingerprints: []string{"hash://sha256/9564b85669d5e96ac969dd0161b8475bbced9e5999c6ec598da718a3045d6f2e"},
			Stakeholders: []IscnInput{
				IscnInput(`{"name":"chung","description":"developer"}`),
			},
			ContentMetadata: IscnInput(`{"title": "iscn module", "description": "a Cosmos SDK module"}`),
		}
	}
	id1, _ := ParseIscnId("iscn://likecoin-chain/btC7CJvMm4WLj9Tau9LAPTfGK7sfymTJW7ORcFdruCU/1")
	id2, _ := ParseIscnId("iscn://likecoin-chain/btC7CJvMm4WLj9Tau9LAPTfGK7sfymTJW7ORcFdruCU/2")
	cid, _ := gocid.Decode("bahuaierav3bfvm4ytx7gvn4yqeu4piiocuvtvdpyyb5f6moxniwemae4tjyq")
	jsonLdInfo := IscnRecordJsonLdInfo{
		Id:         id2,
		Timestamp:  time.Unix(1234567890, 0),
		ParentIpld: &cid,
	}

	var record IscnRecord
	var err error
	var bz []byte
	var expected []byte

	record = goodRecord()
	bz, err = record.ToJsonLd(&jsonLdInfo)
	require.NoError(t, err)
	expected = []byte(`{"@context":{"@vocab":"http://iscn.io/","contentMetadata":{"@context":null},"recordParentIPLD":{"@container":"@index"},"stakeholders":{"@context":{"@vocab":"http://schema.org/","contributionType":"http://iscn.io/contributionType","entity":"http://iscn.io/entity","footprint":"http://iscn.io/footprint","rewardProportion":"http://iscn.io/rewardProportion"}}},"@id":"iscn://likecoin-chain/btC7CJvMm4WLj9Tau9LAPTfGK7sfymTJW7ORcFdruCU/2","@type":"Record","contentFingerprints":["hash://sha256/9564b85669d5e96ac969dd0161b8475bbced9e5999c6ec598da718a3045d6f2e"],"contentMetadata":{"description":"a Cosmos SDK module","title":"iscn module"},"recordNotes":"testing","recordParentIPLD":{"/":"bahuaierav3bfvm4ytx7gvn4yqeu4piiocuvtvdpyyb5f6moxniwemae4tjyq"},"recordTimestamp":"2009-02-13T23:31:30+00:00","recordVersion":2,"stakeholders":[{"description":"developer","name":"chung"}]}`)
	require.Equal(t, expected, bz)

	record = goodRecord()
	record.RecordNotes = ""
	bz, err = record.ToJsonLd(&jsonLdInfo)
	require.NoError(t, err, "should be able to convert record with empty record notes to JSON LD")
	expected = []byte(`{"@context":{"@vocab":"http://iscn.io/","contentMetadata":{"@context":null},"recordParentIPLD":{"@container":"@index"},"stakeholders":{"@context":{"@vocab":"http://schema.org/","contributionType":"http://iscn.io/contributionType","entity":"http://iscn.io/entity","footprint":"http://iscn.io/footprint","rewardProportion":"http://iscn.io/rewardProportion"}}},"@id":"iscn://likecoin-chain/btC7CJvMm4WLj9Tau9LAPTfGK7sfymTJW7ORcFdruCU/2","@type":"Record","contentFingerprints":["hash://sha256/9564b85669d5e96ac969dd0161b8475bbced9e5999c6ec598da718a3045d6f2e"],"contentMetadata":{"description":"a Cosmos SDK module","title":"iscn module"},"recordNotes":"","recordParentIPLD":{"/":"bahuaierav3bfvm4ytx7gvn4yqeu4piiocuvtvdpyyb5f6moxniwemae4tjyq"},"recordTimestamp":"2009-02-13T23:31:30+00:00","recordVersion":2,"stakeholders":[{"description":"developer","name":"chung"}]}`)
	require.Equal(t, expected, bz)

	record = goodRecord()
	record.ContentMetadata = nil
	_, err = record.ToJsonLd(&jsonLdInfo)
	require.Error(t, err, "should not be able to convert record with invalid content metadata to JSON LD")

	record = goodRecord()
	record.ContentMetadata = IscnInput(`null`)
	bz, err = record.ToJsonLd(&jsonLdInfo)
	require.NoError(t, err, "should be able to convert record with null content metadata to JSON LD")
	expected = []byte(`{"@context":{"@vocab":"http://iscn.io/","contentMetadata":{"@context":null},"recordParentIPLD":{"@container":"@index"},"stakeholders":{"@context":{"@vocab":"http://schema.org/","contributionType":"http://iscn.io/contributionType","entity":"http://iscn.io/entity","footprint":"http://iscn.io/footprint","rewardProportion":"http://iscn.io/rewardProportion"}}},"@id":"iscn://likecoin-chain/btC7CJvMm4WLj9Tau9LAPTfGK7sfymTJW7ORcFdruCU/2","@type":"Record","contentFingerprints":["hash://sha256/9564b85669d5e96ac969dd0161b8475bbced9e5999c6ec598da718a3045d6f2e"],"contentMetadata":null,"recordNotes":"testing","recordParentIPLD":{"/":"bahuaierav3bfvm4ytx7gvn4yqeu4piiocuvtvdpyyb5f6moxniwemae4tjyq"},"recordTimestamp":"2009-02-13T23:31:30+00:00","recordVersion":2,"stakeholders":[{"description":"developer","name":"chung"}]}`)
	require.Equal(t, expected, bz)

	record = goodRecord()
	record.Stakeholders = nil
	bz, err = record.ToJsonLd(&jsonLdInfo)
	require.NoError(t, err, "should be able to convert record with no stakeholders to JSON LD")
	expected = []byte(`{"@context":{"@vocab":"http://iscn.io/","contentMetadata":{"@context":null},"recordParentIPLD":{"@container":"@index"},"stakeholders":{"@context":{"@vocab":"http://schema.org/","contributionType":"http://iscn.io/contributionType","entity":"http://iscn.io/entity","footprint":"http://iscn.io/footprint","rewardProportion":"http://iscn.io/rewardProportion"}}},"@id":"iscn://likecoin-chain/btC7CJvMm4WLj9Tau9LAPTfGK7sfymTJW7ORcFdruCU/2","@type":"Record","contentFingerprints":["hash://sha256/9564b85669d5e96ac969dd0161b8475bbced9e5999c6ec598da718a3045d6f2e"],"contentMetadata":{"description":"a Cosmos SDK module","title":"iscn module"},"recordNotes":"testing","recordParentIPLD":{"/":"bahuaierav3bfvm4ytx7gvn4yqeu4piiocuvtvdpyyb5f6moxniwemae4tjyq"},"recordTimestamp":"2009-02-13T23:31:30+00:00","recordVersion":2,"stakeholders":[]}`)
	require.Equal(t, expected, bz)

	record = goodRecord()
	record.ContentFingerprints = nil
	bz, err = record.ToJsonLd(&jsonLdInfo)
	require.NoError(t, err, "should be able to convert record with no fingerprints to JSON LD")
	expected = []byte(`{"@context":{"@vocab":"http://iscn.io/","contentMetadata":{"@context":null},"recordParentIPLD":{"@container":"@index"},"stakeholders":{"@context":{"@vocab":"http://schema.org/","contributionType":"http://iscn.io/contributionType","entity":"http://iscn.io/entity","footprint":"http://iscn.io/footprint","rewardProportion":"http://iscn.io/rewardProportion"}}},"@id":"iscn://likecoin-chain/btC7CJvMm4WLj9Tau9LAPTfGK7sfymTJW7ORcFdruCU/2","@type":"Record","contentFingerprints":[],"contentMetadata":{"description":"a Cosmos SDK module","title":"iscn module"},"recordNotes":"testing","recordParentIPLD":{"/":"bahuaierav3bfvm4ytx7gvn4yqeu4piiocuvtvdpyyb5f6moxniwemae4tjyq"},"recordTimestamp":"2009-02-13T23:31:30+00:00","recordVersion":2,"stakeholders":[{"description":"developer","name":"chung"}]}`)
	require.Equal(t, expected, bz)

	record = goodRecord()
	bz, err = record.ToJsonLd(&IscnRecordJsonLdInfo{
		Id:         id1,
		Timestamp:  time.Unix(1234567890, 0),
		ParentIpld: nil,
	})
	require.NoError(t, err, "should be able to convert record with no parent IPLD to JSON LD")
	expected = []byte(`{"@context":{"@vocab":"http://iscn.io/","contentMetadata":{"@context":null},"recordParentIPLD":{"@container":"@index"},"stakeholders":{"@context":{"@vocab":"http://schema.org/","contributionType":"http://iscn.io/contributionType","entity":"http://iscn.io/entity","footprint":"http://iscn.io/footprint","rewardProportion":"http://iscn.io/rewardProportion"}}},"@id":"iscn://likecoin-chain/btC7CJvMm4WLj9Tau9LAPTfGK7sfymTJW7ORcFdruCU/1","@type":"Record","contentFingerprints":["hash://sha256/9564b85669d5e96ac969dd0161b8475bbced9e5999c6ec598da718a3045d6f2e"],"contentMetadata":{"description":"a Cosmos SDK module","title":"iscn module"},"recordNotes":"testing","recordTimestamp":"2009-02-13T23:31:30+00:00","recordVersion":1,"stakeholders":[{"description":"developer","name":"chung"}]}`)
	require.Equal(t, expected, bz)
}
