package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestParamsValidate(t *testing.T) {
	var err error
	var params Params

	params = DefaultParams()
	err = params.Validate()
	require.NoError(t, err)

	params = Params{
		RegistryName: "likecoin-chain",
		FeePerByte:   sdk.NewDecCoin("nanolike", sdk.NewInt(10000)),
	}
	err = params.Validate()
	require.NoError(t, err)

	params = Params{
		RegistryName: "",
		FeePerByte:   sdk.NewDecCoin("nanolike", sdk.NewInt(10000)),
	}
	err = params.Validate()
	require.Error(t, err, "should not accept empty registry ID")

	params = Params{
		RegistryName: "Like_coin-cha.in:1337",
		FeePerByte:   sdk.NewDecCoin("nanolike", sdk.NewInt(10000)),
	}
	err = params.Validate()
	require.NoError(t, err, "should accept registry name with `_`, `-`, `.`, `:`")

	params = Params{
		RegistryName: "likecoin!chain",
		FeePerByte:   sdk.NewDecCoin("nanolike", sdk.NewInt(10000)),
	}
	err = params.Validate()
	require.Error(t, err, "should not accept registry ID with `!`")

	params = Params{
		RegistryName: "likecoin?chain",
		FeePerByte:   sdk.NewDecCoin("nanolike", sdk.NewInt(10000)),
	}
	err = params.Validate()
	require.Error(t, err, "should not accept registry ID with `?`")

	params = Params{
		RegistryName: "likecoin=chain",
		FeePerByte:   sdk.NewDecCoin("nanolike", sdk.NewInt(10000)),
	}
	err = params.Validate()
	require.Error(t, err, "should not accept registry ID with `=`")

	params = Params{
		RegistryName: "likecoin-chain",
	}
	err = params.Validate()
	require.Error(t, err, "should not accept empty fee per byte parameter")

	params = Params{
		RegistryName: "likecoin-chain",
		FeePerByte:   sdk.NewDecCoin("nanolike", sdk.NewInt(0)), // cannot construct coin with negative amount directly
	}
	params.FeePerByte.Amount = sdk.NewDec(-1)
	err = params.Validate()
	require.Error(t, err, "should not accept negative fee per byte parameter")

	params = Params{
		RegistryName: "likecoin-chain",
		FeePerByte:   sdk.NewDecCoin("nanolike", sdk.NewInt(0)), // cannot construct coin with empty denom directly
	}
	params.FeePerByte.Denom = ""
	err = params.Validate()
	require.Error(t, err, "should not accept empty denom fee per byte parameter")

	var dec sdk.Dec
	dec, err = sdk.NewDecFromStr("0.123")
	require.NoError(t, err)
	params = Params{
		RegistryName: "likecoin-chain",
		FeePerByte:   sdk.NewDecCoinFromDec("nanolike", dec),
	}
	err = params.Validate()
	require.NoError(t, err, "should accept decimal fee per byte parameter")
}
