package types

type GenesisState struct {
	Params Params `json:"params" yaml:"params"`
	// TODO: CIDs and ISCN related metadata
}

func DefaultGenesisState() GenesisState {
	return GenesisState{} // TODO: default param
}

func ValidateGenesis(data GenesisState) error {
	return nil // TODO
}
