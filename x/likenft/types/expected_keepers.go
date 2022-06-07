//go:generate mockgen -source=$GOFILE -destination=../testutil/generated_mock_keepers.go -package=testutil
package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	nft "github.com/likecoin/likechain/backport/cosmos-sdk/v0.46.0-alpha2/x/nft"
	iscntypes "github.com/likecoin/likechain/x/iscn/types"
)

type IscnKeeper interface {
	// Methods imported from iscn should be defined here
	GetContentIdRecord(ctx sdk.Context, iscnIdPrefix iscntypes.IscnIdPrefix) *iscntypes.ContentIdRecord
	GetIscnIdSequence(ctx sdk.Context, iscnId iscntypes.IscnId) uint64
	GetStoreRecord(ctx sdk.Context, seq uint64) *iscntypes.StoreRecord
}

type NftKeeper interface {
	// Methods imported from nft should be defined here
	SaveClass(ctx sdk.Context, class nft.Class) error
	GetClass(ctx sdk.Context, classID string) (nft.Class, bool)
	GetTotalSupply(ctx sdk.Context, classID string) uint64
	UpdateClass(ctx sdk.Context, class nft.Class) error
	Mint(ctx sdk.Context, token nft.NFT, receiver sdk.AccAddress) error
	HasNFT(ctx sdk.Context, classID, id string) bool
	GetOwner(ctx sdk.Context, classID string, nftID string) sdk.AccAddress
	Burn(ctx sdk.Context, classID string, nftID string) error
	GetNFTsOfClass(ctx sdk.Context, classID string) (nfts []nft.NFT)
	Update(ctx sdk.Context, token nft.NFT) error
	Transfer(ctx sdk.Context, classID string, nftID string, receiver sdk.AccAddress) error
}

// AccountKeeper defines the expected account keeper used for simulations (noalias)
type AccountKeeper interface {
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) types.AccountI
	// Methods imported from account should be defined here
}

// BankKeeper defines the expected interface needed to retrieve account balances.
type BankKeeper interface {
	SpendableCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin
	SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	// Methods imported from bank should be defined here
}
