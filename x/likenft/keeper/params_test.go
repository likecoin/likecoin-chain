package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	testkeeper "github.com/likecoin/likecoin-chain/v3/testutil/keeper"
	"github.com/likecoin/likecoin-chain/v3/x/likenft/types"
	"github.com/stretchr/testify/require"
)

func TestParamsValidate(t *testing.T) {
	var err error
	var params types.Params

	params = types.DefaultParams()
	err = params.Validate()
	require.NoError(t, err)

	params = types.Params{
		PriceDenom:             "nanolike",
		FeePerByte:             sdk.NewDecCoin("nanolike", sdk.NewInt(123456)),
		MaxOfferDurationDays:   180,
		MaxListingDurationDays: 180,
		MaxRoyaltyBasisPoints:  1000,
	}
	err = params.Validate()
	require.NoError(t, err)

	params = types.Params{
		PriceDenom:             "",
		FeePerByte:             sdk.NewDecCoin("nanolike", sdk.NewInt(123456)),
		MaxOfferDurationDays:   180,
		MaxListingDurationDays: 180,
		MaxRoyaltyBasisPoints:  1000,
	}
	err = params.Validate()
	require.Error(t, err, "should not accept empty price denom")

	params = types.Params{
		PriceDenom:             "nanolike",
		FeePerByte:             sdk.DecCoin{},
		MaxOfferDurationDays:   180,
		MaxListingDurationDays: 180,
		MaxRoyaltyBasisPoints:  1000,
	}
	err = params.Validate()
	require.Error(t, err, "should not accept empty fee per byte")

	params = types.Params{
		PriceDenom:             "nanolike",
		FeePerByte:             sdk.NewDecCoin("nanolike", sdk.NewInt(123456)),
		MaxOfferDurationDays:   0,
		MaxListingDurationDays: 180,
		MaxRoyaltyBasisPoints:  1000,
	}
	err = params.Validate()
	require.Error(t, err, "should not accept 0 max offer duration")

	params = types.Params{
		PriceDenom:             "nanolike",
		FeePerByte:             sdk.NewDecCoin("nanolike", sdk.NewInt(123456)),
		MaxOfferDurationDays:   180,
		MaxListingDurationDays: 0,
		MaxRoyaltyBasisPoints:  1000,
	}
	err = params.Validate()
	require.Error(t, err, "should not accept 0 max listing duration")

	params = types.Params{
		PriceDenom:             "nanolike",
		FeePerByte:             sdk.NewDecCoin("nanolike", sdk.NewInt(123456)),
		MaxOfferDurationDays:   180,
		MaxListingDurationDays: 180,
		MaxRoyaltyBasisPoints:  10001,
	}
	err = params.Validate()
	require.Error(t, err, "should not accept royalty basis points > 100000")

	params = types.Params{
		PriceDenom:             "nanolike123!!!??",
		FeePerByte:             sdk.NewDecCoin("nanolike", sdk.NewInt(123456)),
		MaxOfferDurationDays:   180,
		MaxListingDurationDays: 180,
		MaxRoyaltyBasisPoints:  1000,
	}
	err = params.Validate()
	require.Error(t, err, "should not accept price denom with invalid characters")

	params = types.Params{}
	err = params.Validate()
	require.Error(t, err, "should not accept empty params")
}

func TestGetParams(t *testing.T) {
	k, ctx := testkeeper.LikenftKeeper(t)
	params := types.DefaultParams()

	k.SetParams(ctx, params)

	require.EqualValues(t, params, k.GetParams(ctx))
}
