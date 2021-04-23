package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenesisValidate(t *testing.T) {
	var state *GenesisState
	var err error

	state = DefaultGenesisState()
	err = state.Validate()
	require.NoError(t, err)

	state.Params.RegistryName = ""
	err = state.Validate()
	require.Error(t, err, "should not accpet genesis state with invalid parameters")

	goodState := func() *GenesisState {
		return &GenesisState{
			Params: DefaultParams(),
			ContentIdRecords: []GenesisState_ContentIdRecord{
				{
					IscnId:        "iscn://likecoin-chain/btC7CJvMm4WLj9Tau9LAPTfGK7sfymTJW7ORcFdruCU",
					Owner:         "cosmos1r623mw6k77g6s3t67fy3042u9nshdl49fgvtex",
					LatestVersion: 2,
				},
				{
					IscnId:        "iscn://likecoin-chain/pZWQk7vER3nkA8wCF4E4sJ9AOi3O-p-6kXxl2JkWviM",
					Owner:         "cosmos1r623mw6k77g6s3t67fy3042u9nshdl49fgvtex",
					LatestVersion: 1,
				},
			},
			IscnRecords: []IscnInput{
				IscnInput(`{"@id":"iscn://likecoin-chain/btC7CJvMm4WLj9Tau9LAPTfGK7sfymTJW7ORcFdruCU/1","contentFingerprints":["hash://sha256/9564b85669d5e96ac969dd0161b8475bbced9e5999c6ec598da718a3045d6f2e","ipfs://QmNrgEMcUygbKzZeZgYFosdd27VE9KnWbyUD73bKZJ3bGi"]}`),
				IscnInput(`{"@id":"iscn://likecoin-chain/pZWQk7vER3nkA8wCF4E4sJ9AOi3O-p-6kXxl2JkWviM/1","contentFingerprints":["hash://sha256/9564b85669d5e96ac969dd0161b8475bbced9e5999c6ec598da718a3045d6f2e","ipfs://QmNrgEMcUygbKzZeZgYFosdd27VE9KnWbyUD73bKZJ3bGi"]}`),
				IscnInput(`{"@id":"iscn://likecoin-chain/btC7CJvMm4WLj9Tau9LAPTfGK7sfymTJW7ORcFdruCU/2","contentFingerprints":["hash://sha256/9564b85669d5e96ac969dd0161b8475bbced9e5999c6ec598da718a3045d6f2e","ipfs://QmNrgEMcUygbKzZeZgYFosdd27VE9KnWbyUD73bKZJ3bGi"]}`),
			},
		}
	}

	state = goodState()
	err = state.Validate()
	require.NoError(t, err)

	state = goodState()
	state.ContentIdRecords[0].LatestVersion = 1
	err = state.Validate()
	require.Error(t, err, "should not accept genesis state with wrong latest version in content ID reocrds")

	state = goodState()
	state.ContentIdRecords[0].LatestVersion = 3
	err = state.Validate()
	require.Error(t, err, "should not accept genesis state with wrong latest version in content ID reocrds")

	state = goodState()
	state.ContentIdRecords = state.ContentIdRecords[:1]
	err = state.Validate()
	require.Error(t, err, "should not accept genesis state with missing content ID reocrd")

	state = goodState()
	state.IscnRecords = []IscnInput{state.IscnRecords[0], state.IscnRecords[2]}
	err = state.Validate()
	require.Error(t, err, "should not accept genesis state with dangling content ID reocrd")

	state = goodState()
	r0 := state.IscnRecords[0]
	r2 := state.IscnRecords[2]
	state.IscnRecords[0] = r2
	state.IscnRecords[2] = r0
	err = state.Validate()
	require.Error(t, err, "should not accept records with wrong order")

	state = goodState()
	state.ContentIdRecords[0].Owner = "cosmos1r623mw6k77g6s3t67fy3042u9nshdl49fgvtey" // invalid checksum
	err = state.Validate()
	require.Error(t, err, "should not accept content ID record with invalid owner address")

	state = goodState()
	state.IscnRecords[0] = IscnInput(`{"@id":"iscn://likecoin-chain/btC7CJvMm4WLj9Tau9LAPTfGK7sfymTJW7ORcFdruCU/1","contentFingerprints":"hash://sha256/9564b85669d5e96ac969dd0161b8475bbced9e5999c6ec598da718a3045d6f2e"}`)
	err = state.Validate()
	require.Error(t, err, "should not accept record with invalid fingerprint type")

	state = goodState()
	state.IscnRecords[0] = IscnInput(`{"@id":"iscn://likecoin-chain/btC7CJvMm4WLj9Tau9LAPTfGK7sfymTJW7ORcFdruCU/1","contentFingerprints":["9564b85669d5e96ac969dd0161b8475bbced9e5999c6ec598da718a3045d6f2e","ipfs://QmNrgEMcUygbKzZeZgYFosdd27VE9KnWbyUD73bKZJ3bGi"]}`)
	err = state.Validate()
	require.Error(t, err, "should not accept record with invalid fingerprint format")

	state = goodState()
	state.IscnRecords[0] = IscnInput(`{"@id":"iscn://likecoin-chain/btC7CJvMm4WLj9Tau9LAPTfGK7sfymTJW7ORcFdruCU/1","contentFingerprints":["sha256/9564b85669d5e96ac969dd0161b8475bbced9e5999c6ec598da718a3045d6f2e","ipfs://QmNrgEMcUygbKzZeZgYFosdd27VE9KnWbyUD73bKZJ3bGi"]}`)
	err = state.Validate()
	require.Error(t, err, "should not accept record with invalid fingerprint format")

	state = goodState()
	state.IscnRecords[0] = IscnInput(`{"@id":"iscn://likecoin-chain/btC7CJvMm4WLj9Tau9LAPTfGK7sfymTJW7ORcFdruCU/1","contentFingerprints":["://sha256/9564b85669d5e96ac969dd0161b8475bbced9e5999c6ec598da718a3045d6f2e","ipfs://QmNrgEMcUygbKzZeZgYFosdd27VE9KnWbyUD73bKZJ3bGi"]}`)
	err = state.Validate()
	require.Error(t, err, "should not accept record with invalid fingerprint format")

	state = goodState()
	state.ContentIdRecords[1].IscnId = "iscn://likecoin?chain/pZWQk7vER3nkA8wCF4E4sJ9AOi3O-p-6kXxl2JkWviM"
	state.IscnRecords[1] = IscnInput(`{"@id":"iscn://likecoin?chain/pZWQk7vER3nkA8wCF4E4sJ9AOi3O-p-6kXxl2JkWviM/1","contentFingerprints":["hash://sha256/9564b85669d5e96ac969dd0161b8475bbced9e5999c6ec598da718a3045d6f2e","ipfs://QmNrgEMcUygbKzZeZgYFosdd27VE9KnWbyUD73bKZJ3bGi"]}`)
	err = state.Validate()
	require.Error(t, err, "should not accept record with invalid ISCN ID")

	// iscn://likecoin-chain/btC7CJvMm4WLj9Tau9LAPTfGK7sfymTJW7ORcFdruCU/1
	// iscn://likecoin-chain/pZWQk7vER3nkA8wCF4E4sJ9AOi3O-p-6kXxl2JkWviM/1
	// iscn://likecoin-chain/Mgd7LH0aAAwyEUYW_rU9EKp9J5cb0598PlHSzN4cQiU/
	// cosmos172nhdqasd2t9e8vvqw4cxfnnutt98q7elzluk9
	// cosmos17dj36xsnaszfwpmv92ct6hfkc2m88nqyls2pvd
	// cosmos1r623mw6k77g6s3t67fy3042u9nshdl49fgvtex
}
