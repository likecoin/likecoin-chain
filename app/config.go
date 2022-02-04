package app

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/std"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
)

const (
	appName = "LikeApp"

	HumanCoinUnit = "LIKE"
	BaseCoinUnit  = "nanolike"
	LikeExponent  = 9
)

func RegisterDenoms() {
	err := sdk.RegisterDenom(HumanCoinUnit, sdk.OneDec())
	if err != nil {
		panic(err)
	}
	err = sdk.RegisterDenom(BaseCoinUnit, sdk.NewDecWithPrec(1, LikeExponent))
	if err != nil {
		panic(err)
	}
}

func SetAddressPrefixes() {
	bech32PrefixesAccAddr := []string{"like", "cosmos"}
	bech32PrefixesAccPub := make([]string, 0, len(bech32PrefixesAccAddr))
	bech32PrefixesValAddr := make([]string, 0, len(bech32PrefixesAccAddr))
	bech32PrefixesValPub := make([]string, 0, len(bech32PrefixesAccAddr))
	bech32PrefixesConsAddr := make([]string, 0, len(bech32PrefixesAccAddr))
	bech32PrefixesConsPub := make([]string, 0, len(bech32PrefixesAccAddr))

	for _, prefix := range bech32PrefixesAccAddr {
		bech32PrefixesAccPub = append(bech32PrefixesAccPub, prefix+"pub")
		bech32PrefixesValAddr = append(bech32PrefixesValAddr, prefix+"valoper")
		bech32PrefixesValPub = append(bech32PrefixesValPub, prefix+"valoperpub")
		bech32PrefixesConsAddr = append(bech32PrefixesConsAddr, prefix+"valcons")
		bech32PrefixesConsPub = append(bech32PrefixesConsPub, prefix+"valconspub")
	}
	config := sdk.GetConfig()
	config.SetBech32PrefixesForAccount(bech32PrefixesAccAddr, bech32PrefixesAccPub)
	config.SetBech32PrefixesForValidator(bech32PrefixesValAddr, bech32PrefixesValPub)
	config.SetBech32PrefixesForConsensusNode(bech32PrefixesConsAddr, bech32PrefixesConsPub)
}

type EncodingConfig struct {
	InterfaceRegistry types.InterfaceRegistry
	Marshaler         codec.Marshaler
	TxConfig          client.TxConfig
	Amino             *codec.LegacyAmino
}

// MakeEncodingConfig creates an EncodingConfig for testing
func MakeEncodingConfig() EncodingConfig {
	cdc := codec.NewLegacyAmino()
	interfaceRegistry := types.NewInterfaceRegistry()
	marshaler := codec.NewProtoCodec(interfaceRegistry)
	std.RegisterLegacyAminoCodec(cdc)
	std.RegisterInterfaces(interfaceRegistry)
	ModuleBasics.RegisterLegacyAminoCodec(cdc)
	ModuleBasics.RegisterInterfaces(interfaceRegistry)

	return EncodingConfig{
		InterfaceRegistry: interfaceRegistry,
		Marshaler:         marshaler,
		TxConfig:          authtx.NewTxConfig(marshaler, authtx.DefaultSignModes),
		Amino:             cdc,
	}
}

func init() {
	SetAddressPrefixes()
	RegisterDenoms()
}
