package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type GenesisState struct {
	Whitelist []sdk.ValAddress `json:"whitelist" yaml:"whitelist"`
	Params    Params           `json:"params" yaml:"params"`
}

func DefaultGenesisState() GenesisState {
	return GenesisState{}
}

func ValidateGenesis(data GenesisState) error {
	return nil
}
