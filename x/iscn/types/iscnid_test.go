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
	require.Equal(t, id.TracingId, []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x00})
	require.Equal(t, id.Version, uint64(1))

	idStr = "iscn://likecoin-chain/AQIDBAUGBwgJAA"
	id, err = ParseIscnID(idStr)
	require.NoError(t, err, "should accept ISCN ID without version")
	require.Equal(t, id.RegistryId, "likecoin-chain")
	require.Equal(t, id.TracingId, []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x00})
	require.Equal(t, id.Version, uint64(0))

	idStr = "iscn://; drop table iscn; --/AQIDBAUGBwgJAA/1"
	id, err = ParseIscnID(idStr)
	require.Error(t, err, "should not accept ISCN ID with registry ID with invalid characters")

	idStr = "iscn://likecoin chain/AQIDBAUGBwgJAA/1"
	id, err = ParseIscnID(idStr)
	require.Error(t, err, "should not accept ISCN ID with registry ID with invalid characters")

	idStr = "iscn://likecoin-chain!/AQIDBAUGBwgJAA/1"
	id, err = ParseIscnID(idStr)
	require.Error(t, err, "should not accept ISCN ID with registry ID with invalid characters")

	idStr = "iscn://likecoin~chain/AQIDBAUGBwgJAA/1"
	id, err = ParseIscnID(idStr)
	require.Error(t, err, "should not accept ISCN ID with registry ID with invalid characters")

	idStr = "iscn://likecoin=chain/AQIDBAUGBwgJAA/1"
	id, err = ParseIscnID(idStr)
	require.Error(t, err, "should not accept ISCN ID with registry ID with invalid characters")

	idStr = "iscn://likecoin-chain;/AQIDBAUGBwgJAA/1"
	id, err = ParseIscnID(idStr)
	require.Error(t, err, "should not accept ISCN ID with registry ID with invalid characters")

	idStr = "iscn://like_coin-chain.is:good/AQIDBAUGBwgJAA/1"
	id, err = ParseIscnID(idStr)
	require.NoError(t, err, "should accept ISCN ID with registry ID with ['.','-','_',':']")

	idStr = "iscn://likecoin-chain/AQIDBAUGBwgJAA/1/2"
	id, err = ParseIscnID(idStr)
	require.Error(t, err, "should not accept ISCN ID with additional sub-path")

	idStr = "iscn://likecoin-chain/AQIDBAUGBwgJAA/1#2"
	_, err = ParseIscnID(idStr)
	require.Error(t, err, "should not accept URL with hash")

	idStr = "iscn://likecoin-chain/AQIDBAUGBwgJAA/1?a=1"
	_, err = ParseIscnID(idStr)
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
	require.Equal(t, id.TracingId, []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x00})
	require.Equal(t, id.Version, uint64(0))

	idStr = "iscn://likecoin-chain/?/1"
	_, err = ParseIscnID(idStr)
	require.Error(t, err, "should not accept non-base64 as tracing ID")

	idStr = "iscn://likecoin-chain/AQIDBAUGBwgJAA==/1"
	_, err = ParseIscnID(idStr)
	require.Error(t, err, "should not accept padded base64 as tracing ID")

	idStr = "iscn://likecoin-chain/++++/1"
	_, err = ParseIscnID(idStr)
	require.Error(t, err, "should not accept non-URL base64 as tracing ID")

	idStr = "isbn://likecoin-chain/AQIDBAUGBwgJAA/1" // note that the scheme is "isbn://"
	_, err = ParseIscnID(idStr)
	require.Error(t, err, "should not accept URL with non-iscn scheme")

	idStr = "iscn://likecoin-chain/AQIDBAUGBwgJAA/a"
	_, err = ParseIscnID(idStr)
	require.Error(t, err, "should not accept non-numerical version")

	idStr = "iscn://likecoin-chain/AQIDBAUGBwgJAA/-1"
	_, err = ParseIscnID(idStr)
	require.Error(t, err, "should not accept negative number version")

	idStr = "iscn://likecoin-chain/AQIDBAUGBwgJAA/0.1"
	_, err = ParseIscnID(idStr)
	require.Error(t, err, "should not accept negative number version")
}
