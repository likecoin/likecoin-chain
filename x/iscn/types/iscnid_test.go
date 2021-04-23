package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseIscnID(t *testing.T) {
	var idStr string
	var id IscnId
	var err error

	idStr = "iscn://likecoin-chain/AQIDBAUGBwgJAA/1"
	id, err = ParseIscnId(idStr)
	require.NoError(t, err)
	require.Equal(t, "likecoin-chain", id.Prefix.RegistryName)
	require.Equal(t, "AQIDBAUGBwgJAA", id.Prefix.ContentId)
	require.Equal(t, uint64(1), id.Version)

	idStr = "iscn://likecoin-chain/AQIDBAUGBwgJAA"
	id, err = ParseIscnId(idStr)
	require.NoError(t, err, "should accept ISCN ID without version")
	require.Equal(t, "likecoin-chain", id.Prefix.RegistryName)
	require.Equal(t, "AQIDBAUGBwgJAA", id.Prefix.ContentId)
	require.Equal(t, uint64(0), id.Version)

	idStr = "iscn://likecoin chain/AQIDBAUGBwgJAA/1"
	id, err = ParseIscnId(idStr)
	require.Error(t, err, "should not accept ISCN ID with registry name with invalid characters")

	idStr = "iscn://likecoin-chain!/AQIDBAUGBwgJAA/1"
	id, err = ParseIscnId(idStr)
	require.Error(t, err, "should not accept ISCN ID with registry name with invalid characters")

	idStr = "iscn://likecoin~chain/AQIDBAUGBwgJAA/1"
	id, err = ParseIscnId(idStr)
	require.Error(t, err, "should not accept ISCN ID with registry name with invalid characters")

	idStr = "iscn://likecoin-chain;/AQIDBAUGBwgJAA/1"
	id, err = ParseIscnId(idStr)
	require.Error(t, err, "should not accept ISCN ID with registry name with invalid characters")

	idStr = "iscn://likecoin-chain?/AQIDBAUGBwgJAA/1"
	id, err = ParseIscnId(idStr)
	require.Error(t, err, "should not accept ISCN ID with registry name with invalid characters")

	idStr = "iscn://like_coin-chain.is:good,1+1=2/record_id=123-a.b_c,d:e/1"
	id, err = ParseIscnId(idStr)
	require.NoError(t, err, "should accept ISCN ID with registry name and content ID with ['.','-','_',',',':','+','=']")

	idStr = "iscn://likecoin-chain/AQIDBAUGBwgJAA /1"
	_, err = ParseIscnId(idStr)
	require.Error(t, err, "should not accept ISCN ID with content ID with invalid characters")

	idStr = "iscn://likecoin-chain/AQIDBAUGBwgJAA!/1"
	_, err = ParseIscnId(idStr)
	require.Error(t, err, "should not accept ISCN ID with content ID with invalid characters")

	idStr = "iscn://likecoin-chain/AQIDBAUGBwgJAA~/1"
	_, err = ParseIscnId(idStr)
	require.Error(t, err, "should not accept ISCN ID with content ID with invalid characters")

	idStr = "iscn://likecoin-chain/AQIDBAUGBwgJAA;/1"
	_, err = ParseIscnId(idStr)
	require.Error(t, err, "should not accept ISCN ID with content ID with invalid characters")

	idStr = "iscn://likecoin-chain/AQIDBAUGBwgJAA?/1"
	_, err = ParseIscnId(idStr)
	require.Error(t, err, "should not accept ISCN ID with content ID with invalid characters")

	idStr = "iscn://likecoin-chain/AQIDBAUGBwgJAA/1/2"
	id, err = ParseIscnId(idStr)
	require.Error(t, err, "should not accept ISCN ID with additional sub-path")

	idStr = "iscn://likecoin-chain/AQIDBAUGBwgJAA/1#2"
	id, err = ParseIscnId(idStr)
	require.Error(t, err, "should not accept URL with hash")

	idStr = "iscn://likecoin-chain/AQIDBAUGBwgJAA/1?a=1"
	id, err = ParseIscnId(idStr)
	require.Error(t, err, "should not accept URL with query")

	idStr = "iscn:///AQIDBAUGBwgJAA/1"
	id, err = ParseIscnId(idStr)
	require.Error(t, err, "should not accept ISCN ID with empty registry name")

	idStr = "iscn://likecoin-chain//1"
	id, err = ParseIscnId(idStr)
	require.Error(t, err, "should not accept ISCN ID with empty content ID")

	idStr = "iscn://likecoin-chain/AQIDBAUGBwgJAA/0"
	id, err = ParseIscnId(idStr)
	require.NoError(t, err, "should accept 0 version")
	require.Equal(t, "likecoin-chain", id.Prefix.RegistryName)
	require.Equal(t, "AQIDBAUGBwgJAA", id.Prefix.ContentId)
	require.Equal(t, uint64(0), id.Version)

	idStr = "iscn://likecoin-chain/AQIDBAUGBwgJAA==/1"
	id, err = ParseIscnId(idStr)
	require.NoError(t, err, "should accept padded base64 as content ID")
	require.Equal(t, "likecoin-chain", id.Prefix.RegistryName)
	require.Equal(t, "AQIDBAUGBwgJAA==", id.Prefix.ContentId)
	require.Equal(t, uint64(1), id.Version)

	idStr = "iscn://likecoin-chain/++++/1"
	id, err = ParseIscnId(idStr)
	require.NoError(t, err, "should accept non-URL base64 as content ID")
	require.Equal(t, "likecoin-chain", id.Prefix.RegistryName)
	require.Equal(t, "++++", id.Prefix.ContentId)
	require.Equal(t, uint64(1), id.Version)

	idStr = "isbn://likecoin-chain/AQIDBAUGBwgJAA/1" // note that the scheme is "isbn://"
	id, err = ParseIscnId(idStr)
	require.Error(t, err, "should not accept URL with non-iscn scheme")

	idStr = "iscn://likecoin-chain/AQIDBAUGBwgJAA/a"
	id, err = ParseIscnId(idStr)
	require.Error(t, err, "should not accept non-numerical version")

	idStr = "iscn://likecoin-chain/AQIDBAUGBwgJAA/-1"
	id, err = ParseIscnId(idStr)
	require.Error(t, err, "should not accept negative number version")

	idStr = "iscn://likecoin-chain/AQIDBAUGBwgJAA/0.1"
	id, err = ParseIscnId(idStr)
	require.Error(t, err, "should not accept negative number version")

	idStr = "iscn://likecoin-chain/AQIDBAUGBwgJAA/1e9"
	id, err = ParseIscnId(idStr)
	require.Error(t, err, "should not accept number version with scientific notation")
}

func TestIscnPrefixId(t *testing.T) {
	idStr := "iscn://likecoin-chain/AQIDBAUGBwgJAA/1"
	id, err := ParseIscnId(idStr)
	require.NoError(t, err)
	require.Equal(t, "likecoin-chain", id.Prefix.RegistryName)
	require.Equal(t, "AQIDBAUGBwgJAA", id.Prefix.ContentId)
	require.Equal(t, uint64(1), id.Version)

	idPrefixStr := id.Prefix.String()
	require.Equal(t, "iscn://likecoin-chain/AQIDBAUGBwgJAA", idPrefixStr)

	prefixId := id.PrefixId()
	require.Equal(t, "likecoin-chain", prefixId.Prefix.RegistryName)
	require.Equal(t, "AQIDBAUGBwgJAA", prefixId.Prefix.ContentId)
	require.Equal(t, uint64(0), prefixId.Version)

	id2Str := "iscn://likecoin-chain/AQIDBAUGBwgJAA/2"
	id2, err := ParseIscnId(id2Str)
	require.NoError(t, err)
	require.Equal(t, "likecoin-chain", id.Prefix.RegistryName)
	require.Equal(t, "AQIDBAUGBwgJAA", id.Prefix.ContentId)
	require.Equal(t, uint64(2), id.Version)

	require.True(t, id.PrefixEqual(&id2))
	require.True(t, id2.PrefixEqual(&id))
}

func TestIscnIdJson(t *testing.T) {
	idStr := "iscn://likecoin-chain/AQIDBAUGBwgJAA/1"
	id, err := ParseIscnId(idStr)
	require.NoError(t, err)
	require.Equal(t, "likecoin-chain", id.Prefix.RegistryName)
	require.Equal(t, "AQIDBAUGBwgJAA", id.Prefix.ContentId)
	require.Equal(t, uint64(1), id.Version)

	json, err := id.MarshalJSON()
	require.NoError(t, err)
	require.Equal(t, []byte(`"iscn://likecoin-chain/AQIDBAUGBwgJAA/1"`), json)

	unmarshalledId := IscnId{}
	err = unmarshalledId.UnmarshalJSON(json)
	require.NoError(t, err)
	require.Equal(t, id, unmarshalledId)

	invalidIdJson1 := []byte(`"iscn://likecoin-chain/AQ!DBAUGBwgJAA/1"`)
	err = unmarshalledId.UnmarshalJSON(invalidIdJson1)
	require.Error(t, err)
	invalidIdJson2 := []byte(`iscn://likecoin-chain/AQIDBAUGBwgJAA/1`) // no quotation
	err = unmarshalledId.UnmarshalJSON(invalidIdJson2)
	require.Error(t, err)

	err = unmarshalledId.UnmarshalJSON(nil)
	require.Error(t, err)
}
