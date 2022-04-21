package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DefaultIndex is the default capability global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default Capability genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		ClassesByISCNList:    []ClassesByISCN{},
		ClassesByAccountList: []ClassesByAccount{},
		ClaimableNFTList:     []ClaimableNFT{},
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
		acc, err := sdk.AccAddressFromBech32(elem.Account)
		if err != nil {
			return fmt.Errorf("Invalid account address: %s", err.Error())
		}

		index := string(ClassesByAccountKey(acc))

		if _, ok := classesByAccountIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for classesByAccount")
		}
		classesByAccountIndexMap[index] = struct{}{}
	}
	// Check for duplicated index in claimableNFT
	claimableNFTIndexMap := make(map[string]struct{})

	for _, elem := range gs.ClaimableNFTList {
		index := string(ClaimableNFTKey(elem.ClassId, elem.Id))
		if _, ok := claimableNFTIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for claimableNFT")
		}
		claimableNFTIndexMap[index] = struct{}{}
	}
	// this line is used by starport scaffolding # genesis/types/validate

	return gs.Params.Validate()
}
