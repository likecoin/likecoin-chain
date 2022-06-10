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
		ClassesByIscnList:    []ClassesByISCN{},
		ClassesByAccountList: []ClassesByAccount{},
		MintableNftList:      []MintableNFT{},
		ClassRevealQueue:     []ClassRevealQueueEntry{},
		OfferList:            []Offer{},
		ListingList:          []Listing{},
		OfferExpireQueue:     []OfferExpireQueueEntry{},
		ListingExpireQueue:   []ListingExpireQueueEntry{},
		// this line is used by starport scaffolding # genesis/types/default
		Params: DefaultParams(),
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// Check for duplicated index in classesByISCN
	classesByISCNIndexMap := make(map[string]struct{})

	for _, elem := range gs.ClassesByIscnList {
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
	// Check for duplicated index in mintableNFT
	mintableNFTIndexMap := make(map[string]struct{})

	for _, elem := range gs.MintableNftList {
		index := string(MintableNFTKey(elem.ClassId, elem.Id))
		if _, ok := mintableNFTIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for mintableNFT")
		}
		mintableNFTIndexMap[index] = struct{}{}
	}

	// Check for duplicated index in classRevealQueue
	classRevealQueueIndexMap := make(map[string]struct{})

	for _, elem := range gs.ClassRevealQueue {
		index := string(ClassRevealQueueKey(elem.RevealTime, elem.ClassId))
		if _, ok := classRevealQueueIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for classRevealQueueEntry")
		}
		classRevealQueueIndexMap[index] = struct{}{}
	}
	// Check for duplicated index in offer
	offerIndexMap := make(map[string]struct{})

	for _, elem := range gs.OfferList {
		acc, err := sdk.AccAddressFromBech32(elem.Buyer)
		if err != nil {
			return fmt.Errorf("Invalid account address: %s", err.Error())
		}
		index := string(OfferKey(elem.ClassId, elem.NftId, acc))
		if _, ok := offerIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for offer")
		}
		offerIndexMap[index] = struct{}{}
	}
	// Check for duplicated index in listing
	listingIndexMap := make(map[string]struct{})

	for _, elem := range gs.ListingList {
		acc, err := sdk.AccAddressFromBech32(elem.Seller)
		if err != nil {
			return fmt.Errorf("Invalid account address: %s", err.Error())
		}
		index := string(ListingKey(elem.ClassId, elem.NftId, acc))
		if _, ok := listingIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for listing")
		}
		listingIndexMap[index] = struct{}{}
	}
	// Check for duplicated index in offerExpireQueueEntry
	offerExpireQueueEntryIndexMap := make(map[string]struct{})

	for _, elem := range gs.OfferExpireQueue {
		index := string(OfferExpireQueueKey(elem.ExpireTime, elem.OfferKey))
		if _, ok := offerExpireQueueEntryIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for offerExpireQueueEntry")
		}
		offerExpireQueueEntryIndexMap[index] = struct{}{}
	}
	// Check for duplicated index in listingExpireQueueEntry
	listingExpireQueueEntryIndexMap := make(map[string]struct{})

	for _, elem := range gs.ListingExpireQueue {
		index := string(ListingExpireQueueKey(elem.ExpireTime, elem.ListingKey))
		if _, ok := listingExpireQueueEntryIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for listingExpireQueueEntry")
		}
		listingExpireQueueEntryIndexMap[index] = struct{}{}
	}
	// this line is used by starport scaffolding # genesis/types/validate

	return gs.Params.Validate()
}
