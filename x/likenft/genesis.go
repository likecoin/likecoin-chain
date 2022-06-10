package likenft

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/likecoin/likechain/x/likenft/keeper"
	"github.com/likecoin/likechain/x/likenft/types"
)

// InitGenesis initializes the capability module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// Set all the classesByISCN
	for _, elem := range genState.ClassesByIscnList {
		k.SetClassesByISCN(ctx, elem)
	}
	// Set all the classesByAccount
	for _, elem := range genState.ClassesByAccountList {
		k.SetClassesByAccount(ctx, elem)
	}
	// Set all the mintableNFT
	for _, elem := range genState.MintableNftList {
		k.SetMintableNFT(ctx, elem)
	}

	// Set all the classRevealQueueEntry
	for _, elem := range genState.ClassRevealQueue {
		k.SetClassRevealQueueEntry(ctx, elem)
	}
	// Set all the offer
	for _, elem := range genState.OfferList {
		k.SetOffer(ctx, elem.ToStoreRecord())
	}
	// Set all the listing
	for _, elem := range genState.ListingList {
		k.SetListing(ctx, elem.ToStoreRecord())
	}
	// Set all the offerExpireQueueEntry
	for _, elem := range genState.OfferExpireQueue {
		k.SetOfferExpireQueueEntry(ctx, elem)
	}
	// Set all the listingExpireQueueEntry
	for _, elem := range genState.ListingExpireQueue {
		k.SetListingExpireQueueEntry(ctx, elem)
	}
	// this line is used by starport scaffolding # genesis/module/init
	k.SetParams(ctx, genState.Params)
}

// ExportGenesis returns the capability module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

	genesis.ClassesByIscnList = k.GetAllClassesByISCN(ctx)
	genesis.ClassesByAccountList = k.GetAllClassesByAccount(ctx)
	genesis.MintableNftList = k.GetAllMintableNFT(ctx)
	genesis.ClassRevealQueue = k.GetClassRevealQueue(ctx)
	genesis.OfferList = types.MapOffersToPublicRecords(k.GetAllOffer(ctx))
	genesis.ListingList = types.MapListingsToPublicRecords(k.GetAllListing(ctx))
	genesis.OfferExpireQueue = k.GetOfferExpireQueue(ctx)
	genesis.ListingExpireQueue = k.GetListingExpireQueue(ctx)
	// this line is used by starport scaffolding # genesis/module/export

	return genesis
}
