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
	id, err = ParseIscnID(idStr)
	require.NoError(t, err)
	require.Equal(t, id.RegistryId, "likecoin-chain")
	require.Equal(t, id.TracingId, "AQIDBAUGBwgJAA")
	require.Equal(t, id.Version, uint64(1))

	idStr = "iscn://likecoin-chain/AQIDBAUGBwgJAA"
	id, err = ParseIscnID(idStr)
	require.NoError(t, err, "should accept ISCN ID without version")
	require.Equal(t, id.RegistryId, "likecoin-chain")
	require.Equal(t, id.TracingId, "AQIDBAUGBwgJAA")
	require.Equal(t, id.Version, uint64(0))

	idStr = "iscn://likecoin chain/AQIDBAUGBwgJAA/1"
	id, err = ParseIscnID(idStr)
	require.Error(t, err, "should not accept ISCN ID with registry ID with invalid characters")

	idStr = "iscn://likecoin-chain!/AQIDBAUGBwgJAA/1"
	id, err = ParseIscnID(idStr)
	require.Error(t, err, "should not accept ISCN ID with registry ID with invalid characters")

	idStr = "iscn://likecoin~chain/AQIDBAUGBwgJAA/1"
	id, err = ParseIscnID(idStr)
	require.Error(t, err, "should not accept ISCN ID with registry ID with invalid characters")

	idStr = "iscn://likecoin-chain;/AQIDBAUGBwgJAA/1"
	id, err = ParseIscnID(idStr)
	require.Error(t, err, "should not accept ISCN ID with registry ID with invalid characters")

	idStr = "iscn://likecoin-chain?/AQIDBAUGBwgJAA/1"
	id, err = ParseIscnID(idStr)
	require.Error(t, err, "should not accept ISCN ID with registry ID with invalid characters")

	idStr = "iscn://like_coin-chain.is:good,1+1=2/record_id=123-a.b_c,d:e/1"
	id, err = ParseIscnID(idStr)
	require.NoError(t, err, "should accept ISCN ID with registry ID and tracing ID with ['.','-','_',',',':','+','=']")

	idStr = "iscn://likecoin-chain/AQIDBAUGBwgJAA /1"
	_, err = ParseIscnID(idStr)
	require.Error(t, err, "should not accept ISCN ID with tracing ID with invalid characters")

	idStr = "iscn://likecoin-chain/AQIDBAUGBwgJAA!/1"
	_, err = ParseIscnID(idStr)
	require.Error(t, err, "should not accept ISCN ID with tracing ID with invalid characters")

	idStr = "iscn://likecoin-chain/AQIDBAUGBwgJAA~/1"
	_, err = ParseIscnID(idStr)
	require.Error(t, err, "should not accept ISCN ID with tracing ID with invalid characters")

	idStr = "iscn://likecoin-chain/AQIDBAUGBwgJAA;/1"
	_, err = ParseIscnID(idStr)
	require.Error(t, err, "should not accept ISCN ID with tracing ID with invalid characters")

	idStr = "iscn://likecoin-chain/AQIDBAUGBwgJAA?/1"
	_, err = ParseIscnID(idStr)
	require.Error(t, err, "should not accept ISCN ID with tracing ID with invalid characters")

	idStr = "iscn://likecoin-chain/AQIDBAUGBwgJAA/1/2"
	id, err = ParseIscnID(idStr)
	require.Error(t, err, "should not accept ISCN ID with additional sub-path")

	idStr = "iscn://likecoin-chain/AQIDBAUGBwgJAA/1#2"
	id, err = ParseIscnID(idStr)
	require.Error(t, err, "should not accept URL with hash")

	idStr = "iscn://likecoin-chain/AQIDBAUGBwgJAA/1?a=1"
	id, err = ParseIscnID(idStr)
	require.Error(t, err, "should not accept URL with query")

	idStr = "iscn:///AQIDBAUGBwgJAA/1"
	id, err = ParseIscnID(idStr)
	require.Error(t, err, "should not accept ISCN ID with empty registry ID")

	idStr = "iscn://likecoin-chain//1"
	id, err = ParseIscnID(idStr)
	require.Error(t, err, "should not accept ISCN ID with empty tracing ID")

	idStr = "iscn://likecoin-chain/AQIDBAUGBwgJAA/0"
	id, err = ParseIscnID(idStr)
	require.NoError(t, err, "should accept 0 version")
	require.Equal(t, id.RegistryId, "likecoin-chain")
	require.Equal(t, id.TracingId, "AQIDBAUGBwgJAA")
	require.Equal(t, id.Version, uint64(0))

	idStr = "iscn://likecoin-chain/AQIDBAUGBwgJAA==/1"
	id, err = ParseIscnID(idStr)
	require.NoError(t, err, "should accept padded base64 as tracing ID")
	require.Equal(t, id.RegistryId, "likecoin-chain")
	require.Equal(t, id.TracingId, "AQIDBAUGBwgJAA==")
	require.Equal(t, id.Version, uint64(1))

	idStr = "iscn://likecoin-chain/++++/1"
	id, err = ParseIscnID(idStr)
	require.NoError(t, err, "should accept non-URL base64 as tracing ID")
	require.Equal(t, id.RegistryId, "likecoin-chain")
	require.Equal(t, id.TracingId, "++++")
	require.Equal(t, id.Version, uint64(1))

	idStr = "isbn://likecoin-chain/AQIDBAUGBwgJAA/1" // note that the scheme is "isbn://"
	id, err = ParseIscnID(idStr)
	require.Error(t, err, "should not accept URL with non-iscn scheme")

	idStr = "iscn://likecoin-chain/AQIDBAUGBwgJAA/a"
	id, err = ParseIscnID(idStr)
	require.Error(t, err, "should not accept non-numerical version")

	idStr = "iscn://likecoin-chain/AQIDBAUGBwgJAA/-1"
	id, err = ParseIscnID(idStr)
	require.Error(t, err, "should not accept negative number version")

	idStr = "iscn://likecoin-chain/AQIDBAUGBwgJAA/0.1"
	id, err = ParseIscnID(idStr)
	require.Error(t, err, "should not accept negative number version")

	idStr = "iscn://likecoin-chain/AQIDBAUGBwgJAA/1e9"
	id, err = ParseIscnID(idStr)
	require.Error(t, err, "should not accept number version with scientific notation")
}
