package types

type IscnPair struct {
	Id     []byte     `json:"id" yaml:"id"`
	Record IscnRecord `json:"record" yaml:"record"`
}

type GenesisState struct {
	Params      Params     `json:"params" yaml:"params"`
	Authors     []Author   `json:"authors" yaml:"authors"`
	RightTerms  []Right    `json:"rightTerms" yaml:"rightTerms"`
	IscnRecords []IscnPair `json:"iscnRecords" yaml:"iscnRecords"`
}

func DefaultGenesisState() GenesisState {
	return GenesisState{} // TODO: default param
}

func ValidateGenesis(data GenesisState) error {
	return nil // TODO
}
