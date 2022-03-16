package dual_prefix_tests

import (
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	priv1       = secp256k1.GenPrivKey()
	legacyAddr1 = sdk.MustBech32ifyAddressBytes("cosmos", priv1.PubKey().Address())
	newAddr1    = sdk.MustBech32ifyAddressBytes("like", priv1.PubKey().Address())
	priv2       = secp256k1.GenPrivKey()
	legacyAddr2 = sdk.MustBech32ifyAddressBytes("cosmos", priv2.PubKey().Address())
	newAddr2    = sdk.MustBech32ifyAddressBytes("like", priv2.PubKey().Address())
)
