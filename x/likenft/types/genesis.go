package types

import (
	"fmt"
)

// DefaultIndex is the default capability global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default Capability genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		ClassesByISCNList:    []ClassesByISCN{},
		ClassesByAccountList: []ClassesByAccount{},
		// this line is used by starport scaffolding # genesis/types/default
		Params: DefaultParams(),
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// Check for duplicated index in classesByISCN
	classesByISCNIndexMap := make(map[string]struct{})

	for _, elem := range gs.ClassesByISCNList {
		index := string(ClassesByISCNKey(elem.IscnIdPrefix))
		if _, ok := classesByISCNIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for classesByISCN")
		}
		classesByISCNIndexMap[index] = struct{}{}
	}
	// Check for duplicated index in classesByAccount
	classesByAccountIndexMap := make(map[string]struct{})

	for _, elem := range gs.ClassesByAccountList {
		index := string(ClassesByAccountKey(elem.Account))
		if _, ok := classesByAccountIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for classesByAccount")
		}
		classesByAccountIndexMap[index] = struct{}{}
	}
	// this line is used by starport scaffolding # genesis/types/validate

	return gs.Params.Validate()
}
