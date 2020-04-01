package types

type IscnPair struct {
	Id     []byte     `json:"id" yaml:"id"`
	Record IscnRecord `json:"record" yaml:"record"`
}

type GenesisState struct {
	Params      Params     `json:"params" yaml:"params"`
	IscnRecords []IscnPair `json:"iscnRecords" yaml:"iscnRecords"`
}

func DefaultGenesisState() GenesisState {
	return GenesisState{}
}

func ValidateGenesis(data GenesisState) error {
	return nil
}
