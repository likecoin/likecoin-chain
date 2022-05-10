package dual_prefix_tests

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/likecoin/likecoin-chain/v2/testutil"
	"github.com/stretchr/testify/require"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

func testSendWithBech32(t *testing.T, from string, to string) {
	app := testutil.SetupTestApp([]testutil.GenesisBalance{
		{
			Address: from,
			Coin:    "1000nanolike",
		},
	})
	app.NextHeader(1234567890)
	app.SetForTx()

	fromAddr, err := sdk.AccAddressFromBech32(from)
	require.NoError(t, err)
	toAddr, err := sdk.AccAddressFromBech32(to)
	require.NoError(t, err)

	msg := banktypes.NewMsgSend(fromAddr, toAddr, sdk.NewCoins(sdk.NewCoin("nanolike", sdk.NewInt(1))))
	app.DeliverMsgNoError(t, msg, priv1)

	ctx := app.SetForQuery()
	fromAccBalance := app.BankKeeper.GetBalance(ctx, fromAddr, "nanolike")
	require.Equal(t, sdk.NewCoin("nanolike", sdk.NewInt(999)), fromAccBalance)
	toAccBalance := app.BankKeeper.GetBalance(ctx, toAddr, "nanolike")
	require.Equal(t, sdk.NewCoin("nanolike", sdk.NewInt(1)), toAccBalance)
}

func TestSendFromLegacyPrefixToNew(t *testing.T) {
	testSendWithBech32(t, legacyAddr1, newAddr2)
}

func TestSendFromLegacyPrefixToLegacy(t *testing.T) {
	testSendWithBech32(t, legacyAddr1, legacyAddr2)
}

func TestSendFromNewPrefixToLegacy(t *testing.T) {
	testSendWithBech32(t, newAddr1, legacyAddr2)
}

func TestSendFromNewPrefixToNew(t *testing.T) {
	testSendWithBech32(t, newAddr1, newAddr2)
}
