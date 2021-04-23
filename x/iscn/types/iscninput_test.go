package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIscnInputNormalize(t *testing.T) {
	var err error
	var input IscnInput
	var bz []byte

	err = input.Unmarshal([]byte(`{"b":2,"a":1}`))
	require.NoError(t, err)
	bz, err = input.Normalize()
	require.NoError(t, err)
	require.Equal(t, []byte(`{"a":1,"b":2}`), bz)

	err = input.Unmarshal([]byte(`{
		"title": "testing",
		"contents": [
			{
				"valid": true,
				"duration": 7200
			},
			{
				"valid": false,
				"duration": 7200
			}
		],
		"extra": null
	}`))
	require.NoError(t, err)
	bz, err = input.Normalize()
	require.NoError(t, err)
	require.Equal(t, []byte(`{"contents":[{"duration":7200,"valid":true},{"duration":7200,"valid":false}],"extra":null,"title":"testing"}`), bz)

	err = input.Unmarshal([]byte(`{
		"title": "testing",
		"contents": [
			{
				"valid": true,
				"duration": 7200
			},
			{
				"valid": false,
				"duration": 7200
			},
		],
		"extra": null
	}`))
	require.NoError(t, err)
	bz, err = input.Normalize()
	require.Error(t, err) // extra comma at the end of `contents`
}

func TestIscnInputValidate(t *testing.T) {
	var err error
	var input IscnInput

	err = input.Unmarshal([]byte(`{
		"title": "testing",
		"contents": [
			{
				"valid": true,
				"duration": 7200
			},
			{
				"valid": false,
				"duration": 7200
			}
		],
		"extra": null
	}`))
	require.NoError(t, err)
	err = input.Validate()
	require.NoError(t, err)

	err = input.Unmarshal([]byte(`{
		"title": "testing",
		"contents": [
			{
				"valid": true,
				"duration": 7200
			},
			{
				"valid": false,
				"duration": 7200
			},
		],
		"extra": null
	}`))
	require.NoError(t, err)
	err = input.Validate()
	require.Error(t, err)

	err = input.Unmarshal([]byte(`""`))
	require.NoError(t, err)
	err = input.Validate()
	require.NoError(t, err)

	err = input.Unmarshal([]byte(`"`))
	require.NoError(t, err)
	err = input.Validate()
	require.Error(t, err)

	err = input.Unmarshal(nil)
	require.NoError(t, err)
	err = input.Validate()
	require.Error(t, err)
}
