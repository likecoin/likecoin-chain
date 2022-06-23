package keeper_test

import (
	"math"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/golang/mock/gomock"
	"github.com/likecoin/likechain/testutil/keeper"
	"github.com/likecoin/likechain/x/likenft/testutil"
	"github.com/likecoin/likechain/x/likenft/types"
	"github.com/stretchr/testify/require"
)

// normal royalty
func TestSellNFTNormalRoyalty(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	accountKeeper := testutil.NewMockAccountKeeper(ctrl)
	bankKeeper := testutil.NewMockBankKeeper(ctrl)
	iscnKeeper := testutil.NewMockIscnKeeper(ctrl)
	nftKeeper := testutil.NewMockNftKeeper(ctrl)
	msgServer, goCtx, k := setupMsgServer(t, keeper.LikenftDependedKeepers{
		AccountKeeper: accountKeeper,
		BankKeeper:    bankKeeper,
		IscnKeeper:    iscnKeeper,
		NftKeeper:     nftKeeper,
	})
	ctx := sdk.UnwrapSDKContext(goCtx)
	ctx = ctx.WithBlockTime(time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC))
	goCtx = sdk.WrapSDKContext(ctx)

	// Data
	creatorAddressBytes := []byte{1, 1, 1, 1, 0, 0, 0, 0}
	sellerAddressBytes := []byte{0, 1, 0, 1, 0, 1, 0, 1}
	sellerAddress, _ := sdk.Bech32ifyAddressBytes("like", sellerAddressBytes)
	buyerAddressBytes := []byte{1, 0, 1, 0, 1, 0, 1, 0}
	buyerAddress, _ := sdk.Bech32ifyAddressBytes("like", buyerAddressBytes)
	classId := "likenft1abcdef"
	nftId := "nft1"
	price := uint64(123456)
	expiration := time.Date(2022, 4, 1, 0, 0, 0, 0, time.UTC)
	royaltyBasisPoints := uint64(234)
	finalPrice := uint64(100000)

	// Seed listing to test deletion after txn
	k.SetListing(ctx, types.ListingStoreRecord{
		ClassId:    classId,
		NftId:      nftId,
		Seller:     sellerAddressBytes,
		Price:      uint64(999999),
		Expiration: time.Date(2022, 2, 1, 0, 0, 0, 0, time.UTC),
	})

	// Seed offer
	k.SetOffer(ctx, types.OfferStoreRecord{
		ClassId:    classId,
		NftId:      nftId,
		Buyer:      buyerAddressBytes,
		Price:      price,
		Expiration: expiration,
	})

	// Seed royalty config
	k.SetRoyaltyConfig(ctx, types.RoyaltyConfigByClass{
		ClassId: classId,
		RoyaltyConfig: types.RoyaltyConfig{
			RateBasisPoints: royaltyBasisPoints,
			Stakeholders: []types.RoyaltyStakeholder{
				{
					Account: creatorAddressBytes,
					Weight:  uint64(1),
				},
			},
		},
	})
	royaltyAmount := uint64(math.Floor(float64(finalPrice) / 10000 * float64(royaltyBasisPoints)))
	royaltyAmountCoins := sdk.NewCoins(sdk.NewCoin(k.GetParams(ctx).PriceDenom, sdk.NewInt(int64(royaltyAmount))))

	// Mock
	nftKeeper.EXPECT().GetOwner(gomock.Any(), classId, nftId).Return(sellerAddressBytes)
	bankKeeper.EXPECT().SendCoinsFromModuleToAccount(gomock.Any(), types.ModuleName, creatorAddressBytes, royaltyAmountCoins).Return(nil)
	netAmount := finalPrice - royaltyAmount
	netAmountCoins := sdk.NewCoins(sdk.NewCoin(k.GetParams(ctx).PriceDenom, sdk.NewInt(int64(netAmount))))
	bankKeeper.EXPECT().SendCoinsFromModuleToAccount(gomock.Any(), types.ModuleName, sellerAddressBytes, netAmountCoins).Return(nil)
	refundAmount := price - finalPrice
	refundAmountCoins := sdk.NewCoins(sdk.NewCoin(k.GetParams(ctx).PriceDenom, sdk.NewInt(int64(refundAmount))))
	bankKeeper.EXPECT().SendCoinsFromModuleToAccount(gomock.Any(), types.ModuleName, buyerAddressBytes, refundAmountCoins).Return(nil)
	nftKeeper.EXPECT().Transfer(gomock.Any(), classId, nftId, buyerAddressBytes).Return(nil)

	// Call
	res, err := msgServer.SellNFT(goCtx, &types.MsgSellNFT{
		Creator: sellerAddress,
		ClassId: classId,
		NftId:   nftId,
		Buyer:   buyerAddress,
		Price:   finalPrice,
	})
	require.NoError(t, err)
	require.Equal(t, &types.MsgSellNFTResponse{}, res)

	// Check state
	// Expect offer removed
	_, found := k.GetOffer(ctx, classId, nftId, buyerAddressBytes)
	require.False(t, found)
	// Expect listing removed
	_, found = k.GetListing(ctx, classId, nftId, sellerAddressBytes)
	require.False(t, found)

	ctrl.Finish()
}

// normal no royalty
func TestSellNFTNormalNoRoyalty(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	accountKeeper := testutil.NewMockAccountKeeper(ctrl)
	bankKeeper := testutil.NewMockBankKeeper(ctrl)
	iscnKeeper := testutil.NewMockIscnKeeper(ctrl)
	nftKeeper := testutil.NewMockNftKeeper(ctrl)
	msgServer, goCtx, k := setupMsgServer(t, keeper.LikenftDependedKeepers{
		AccountKeeper: accountKeeper,
		BankKeeper:    bankKeeper,
		IscnKeeper:    iscnKeeper,
		NftKeeper:     nftKeeper,
	})
	ctx := sdk.UnwrapSDKContext(goCtx)
	ctx = ctx.WithBlockTime(time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC))
	goCtx = sdk.WrapSDKContext(ctx)

	// Data
	sellerAddressBytes := []byte{0, 1, 0, 1, 0, 1, 0, 1}
	sellerAddress, _ := sdk.Bech32ifyAddressBytes("like", sellerAddressBytes)
	buyerAddressBytes := []byte{1, 0, 1, 0, 1, 0, 1, 0}
	buyerAddress, _ := sdk.Bech32ifyAddressBytes("like", buyerAddressBytes)
	classId := "likenft1abcdef"
	nftId := "nft1"
	price := uint64(123456)
	expiration := time.Date(2022, 4, 1, 0, 0, 0, 0, time.UTC)
	finalPrice := uint64(100000)

	// Seed listing to test deletion after txn
	k.SetListing(ctx, types.ListingStoreRecord{
		ClassId:    classId,
		NftId:      nftId,
		Seller:     sellerAddressBytes,
		Price:      uint64(999999),
		Expiration: time.Date(2022, 2, 1, 0, 0, 0, 0, time.UTC),
	})

	// Seed offer
	k.SetOffer(ctx, types.OfferStoreRecord{
		ClassId:    classId,
		NftId:      nftId,
		Buyer:      buyerAddressBytes,
		Price:      price,
		Expiration: expiration,
	})

	// no royalty config

	// Mock
	nftKeeper.EXPECT().GetOwner(gomock.Any(), classId, nftId).Return(sellerAddressBytes)
	netAmount := finalPrice
	netAmountCoins := sdk.NewCoins(sdk.NewCoin(k.GetParams(ctx).PriceDenom, sdk.NewInt(int64(netAmount))))
	bankKeeper.EXPECT().SendCoinsFromModuleToAccount(gomock.Any(), types.ModuleName, sellerAddressBytes, netAmountCoins).Return(nil)
	refundAmount := price - finalPrice
	refundAmountCoins := sdk.NewCoins(sdk.NewCoin(k.GetParams(ctx).PriceDenom, sdk.NewInt(int64(refundAmount))))
	bankKeeper.EXPECT().SendCoinsFromModuleToAccount(gomock.Any(), types.ModuleName, buyerAddressBytes, refundAmountCoins).Return(nil)
	nftKeeper.EXPECT().Transfer(gomock.Any(), classId, nftId, buyerAddressBytes).Return(nil)

	// Call
	res, err := msgServer.SellNFT(goCtx, &types.MsgSellNFT{
		Creator: sellerAddress,
		ClassId: classId,
		NftId:   nftId,
		Buyer:   buyerAddress,
		Price:   finalPrice,
	})
	require.NoError(t, err)
	require.Equal(t, &types.MsgSellNFTResponse{}, res)

	// Check state
	// Expect offer removed
	_, found := k.GetOffer(ctx, classId, nftId, buyerAddressBytes)
	require.False(t, found)
	// Expect listing removed
	_, found = k.GetListing(ctx, classId, nftId, sellerAddressBytes)
	require.False(t, found)

	ctrl.Finish()
}

// normal no refund
func TestSellNFTNoRefund(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	accountKeeper := testutil.NewMockAccountKeeper(ctrl)
	bankKeeper := testutil.NewMockBankKeeper(ctrl)
	iscnKeeper := testutil.NewMockIscnKeeper(ctrl)
	nftKeeper := testutil.NewMockNftKeeper(ctrl)
	msgServer, goCtx, k := setupMsgServer(t, keeper.LikenftDependedKeepers{
		AccountKeeper: accountKeeper,
		BankKeeper:    bankKeeper,
		IscnKeeper:    iscnKeeper,
		NftKeeper:     nftKeeper,
	})
	ctx := sdk.UnwrapSDKContext(goCtx)
	ctx = ctx.WithBlockTime(time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC))
	goCtx = sdk.WrapSDKContext(ctx)

	// Data
	creatorAddressBytes := []byte{1, 1, 1, 1, 0, 0, 0, 0}
	sellerAddressBytes := []byte{0, 1, 0, 1, 0, 1, 0, 1}
	sellerAddress, _ := sdk.Bech32ifyAddressBytes("like", sellerAddressBytes)
	buyerAddressBytes := []byte{1, 0, 1, 0, 1, 0, 1, 0}
	buyerAddress, _ := sdk.Bech32ifyAddressBytes("like", buyerAddressBytes)
	classId := "likenft1abcdef"
	nftId := "nft1"
	price := uint64(123456)
	expiration := time.Date(2022, 4, 1, 0, 0, 0, 0, time.UTC)
	royaltyBasisPoints := uint64(234)
	finalPrice := uint64(123456)

	// Seed listing to test deletion after txn
	k.SetListing(ctx, types.ListingStoreRecord{
		ClassId:    classId,
		NftId:      nftId,
		Seller:     sellerAddressBytes,
		Price:      uint64(999999),
		Expiration: time.Date(2022, 2, 1, 0, 0, 0, 0, time.UTC),
	})

	// Seed offer
	k.SetOffer(ctx, types.OfferStoreRecord{
		ClassId:    classId,
		NftId:      nftId,
		Buyer:      buyerAddressBytes,
		Price:      price,
		Expiration: expiration,
	})

	// Seed royalty config
	k.SetRoyaltyConfig(ctx, types.RoyaltyConfigByClass{
		ClassId: classId,
		RoyaltyConfig: types.RoyaltyConfig{
			RateBasisPoints: royaltyBasisPoints,
			Stakeholders: []types.RoyaltyStakeholder{
				{
					Account: creatorAddressBytes,
					Weight:  uint64(1),
				},
			},
		},
	})
	royaltyAmount := uint64(math.Floor(float64(finalPrice) / 10000 * float64(royaltyBasisPoints)))
	royaltyAmountCoins := sdk.NewCoins(sdk.NewCoin(k.GetParams(ctx).PriceDenom, sdk.NewInt(int64(royaltyAmount))))

	// Mock
	nftKeeper.EXPECT().GetOwner(gomock.Any(), classId, nftId).Return(sellerAddressBytes)
	bankKeeper.EXPECT().SendCoinsFromModuleToAccount(gomock.Any(), types.ModuleName, creatorAddressBytes, royaltyAmountCoins).Return(nil)
	netAmount := finalPrice - royaltyAmount
	netAmountCoins := sdk.NewCoins(sdk.NewCoin(k.GetParams(ctx).PriceDenom, sdk.NewInt(int64(netAmount))))
	bankKeeper.EXPECT().SendCoinsFromModuleToAccount(gomock.Any(), types.ModuleName, sellerAddressBytes, netAmountCoins).Return(nil)
	nftKeeper.EXPECT().Transfer(gomock.Any(), classId, nftId, buyerAddressBytes).Return(nil)

	// Call
	res, err := msgServer.SellNFT(goCtx, &types.MsgSellNFT{
		Creator: sellerAddress,
		ClassId: classId,
		NftId:   nftId,
		Buyer:   buyerAddress,
		Price:   finalPrice,
	})
	require.NoError(t, err)
	require.Equal(t, &types.MsgSellNFTResponse{}, res)

	// Check state
	// Expect offer removed
	_, found := k.GetOffer(ctx, classId, nftId, buyerAddressBytes)
	require.False(t, found)
	// Expect listing removed
	_, found = k.GetListing(ctx, classId, nftId, sellerAddressBytes)
	require.False(t, found)

	ctrl.Finish()
}

// user not owner
func TestSellNFTUserNotOwner(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	accountKeeper := testutil.NewMockAccountKeeper(ctrl)
	bankKeeper := testutil.NewMockBankKeeper(ctrl)
	iscnKeeper := testutil.NewMockIscnKeeper(ctrl)
	nftKeeper := testutil.NewMockNftKeeper(ctrl)
	msgServer, goCtx, k := setupMsgServer(t, keeper.LikenftDependedKeepers{
		AccountKeeper: accountKeeper,
		BankKeeper:    bankKeeper,
		IscnKeeper:    iscnKeeper,
		NftKeeper:     nftKeeper,
	})
	ctx := sdk.UnwrapSDKContext(goCtx)
	ctx = ctx.WithBlockTime(time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC))
	goCtx = sdk.WrapSDKContext(ctx)

	// Data
	sellerAddressBytes := []byte{0, 1, 0, 1, 0, 1, 0, 1}
	sellerAddress, _ := sdk.Bech32ifyAddressBytes("like", sellerAddressBytes)
	notSellerAddressBytes := []byte{0, 0, 0, 0, 1, 1, 1, 1}
	buyerAddressBytes := []byte{1, 0, 1, 0, 1, 0, 1, 0}
	buyerAddress, _ := sdk.Bech32ifyAddressBytes("like", buyerAddressBytes)
	classId := "likenft1abcdef"
	nftId := "nft1"
	price := uint64(123456)
	expiration := time.Date(2022, 4, 1, 0, 0, 0, 0, time.UTC)
	finalPrice := uint64(100000)

	// Seed listing to test deletion after txn
	k.SetListing(ctx, types.ListingStoreRecord{
		ClassId:    classId,
		NftId:      nftId,
		Seller:     sellerAddressBytes,
		Price:      uint64(999999),
		Expiration: time.Date(2022, 2, 1, 0, 0, 0, 0, time.UTC),
	})

	// Seed offer
	k.SetOffer(ctx, types.OfferStoreRecord{
		ClassId:    classId,
		NftId:      nftId,
		Buyer:      buyerAddressBytes,
		Price:      price,
		Expiration: expiration,
	})

	// Mock
	nftKeeper.EXPECT().GetOwner(gomock.Any(), classId, nftId).Return(notSellerAddressBytes)

	// Call
	res, err := msgServer.SellNFT(goCtx, &types.MsgSellNFT{
		Creator: sellerAddress,
		ClassId: classId,
		NftId:   nftId,
		Buyer:   buyerAddress,
		Price:   finalPrice,
	})
	require.Error(t, err)
	require.Nil(t, res)
	require.Contains(t, err.Error(), sdkerrors.ErrUnauthorized.Error())

	// Check state
	// Expect offer not removed
	_, found := k.GetOffer(ctx, classId, nftId, buyerAddressBytes)
	require.True(t, found)
	// Expect listing not removed
	_, found = k.GetListing(ctx, classId, nftId, sellerAddressBytes)
	require.True(t, found)

	ctrl.Finish()
}

// offer not found
func TestSellNFTOfferNotFound(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	accountKeeper := testutil.NewMockAccountKeeper(ctrl)
	bankKeeper := testutil.NewMockBankKeeper(ctrl)
	iscnKeeper := testutil.NewMockIscnKeeper(ctrl)
	nftKeeper := testutil.NewMockNftKeeper(ctrl)
	msgServer, goCtx, k := setupMsgServer(t, keeper.LikenftDependedKeepers{
		AccountKeeper: accountKeeper,
		BankKeeper:    bankKeeper,
		IscnKeeper:    iscnKeeper,
		NftKeeper:     nftKeeper,
	})
	ctx := sdk.UnwrapSDKContext(goCtx)
	ctx = ctx.WithBlockTime(time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC))
	goCtx = sdk.WrapSDKContext(ctx)

	// Data)
	sellerAddressBytes := []byte{0, 1, 0, 1, 0, 1, 0, 1}
	sellerAddress, _ := sdk.Bech32ifyAddressBytes("like", sellerAddressBytes)
	buyerAddressBytes := []byte{1, 0, 1, 0, 1, 0, 1, 0}
	buyerAddress, _ := sdk.Bech32ifyAddressBytes("like", buyerAddressBytes)
	classId := "likenft1abcdef"
	nftId := "nft1"
	finalPrice := uint64(100000)

	// Seed listing to test deletion after txn
	k.SetListing(ctx, types.ListingStoreRecord{
		ClassId:    classId,
		NftId:      nftId,
		Seller:     sellerAddressBytes,
		Price:      uint64(999999),
		Expiration: time.Date(2022, 2, 1, 0, 0, 0, 0, time.UTC),
	})

	// Mock
	nftKeeper.EXPECT().GetOwner(gomock.Any(), classId, nftId).Return(sellerAddressBytes)

	// Call
	res, err := msgServer.SellNFT(goCtx, &types.MsgSellNFT{
		Creator: sellerAddress,
		ClassId: classId,
		NftId:   nftId,
		Buyer:   buyerAddress,
		Price:   finalPrice,
	})
	require.Error(t, err)
	require.Nil(t, res)
	require.Contains(t, err.Error(), types.ErrOfferNotFound.Error())

	// Check state
	// Expect listing not removed
	_, found := k.GetListing(ctx, classId, nftId, sellerAddressBytes)
	require.True(t, found)

	ctrl.Finish()
}

// offer expired
func TestSellNFTOfferExpired(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	accountKeeper := testutil.NewMockAccountKeeper(ctrl)
	bankKeeper := testutil.NewMockBankKeeper(ctrl)
	iscnKeeper := testutil.NewMockIscnKeeper(ctrl)
	nftKeeper := testutil.NewMockNftKeeper(ctrl)
	msgServer, goCtx, k := setupMsgServer(t, keeper.LikenftDependedKeepers{
		AccountKeeper: accountKeeper,
		BankKeeper:    bankKeeper,
		IscnKeeper:    iscnKeeper,
		NftKeeper:     nftKeeper,
	})
	ctx := sdk.UnwrapSDKContext(goCtx)
	ctx = ctx.WithBlockTime(time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC))
	goCtx = sdk.WrapSDKContext(ctx)

	// Data)
	sellerAddressBytes := []byte{0, 1, 0, 1, 0, 1, 0, 1}
	sellerAddress, _ := sdk.Bech32ifyAddressBytes("like", sellerAddressBytes)
	buyerAddressBytes := []byte{1, 0, 1, 0, 1, 0, 1, 0}
	buyerAddress, _ := sdk.Bech32ifyAddressBytes("like", buyerAddressBytes)
	classId := "likenft1abcdef"
	nftId := "nft1"
	price := uint64(123456)
	expiration := time.Date(2021, 12, 31, 0, 0, 0, 0, time.UTC)
	finalPrice := uint64(100000)

	// Seed listing to test deletion after txn
	k.SetListing(ctx, types.ListingStoreRecord{
		ClassId:    classId,
		NftId:      nftId,
		Seller:     sellerAddressBytes,
		Price:      uint64(999999),
		Expiration: time.Date(2022, 2, 1, 0, 0, 0, 0, time.UTC),
	})

	// Seed offer
	k.SetOffer(ctx, types.OfferStoreRecord{
		ClassId:    classId,
		NftId:      nftId,
		Buyer:      buyerAddressBytes,
		Price:      price,
		Expiration: expiration,
	})

	// Mock
	nftKeeper.EXPECT().GetOwner(gomock.Any(), classId, nftId).Return(sellerAddressBytes)

	// Call
	res, err := msgServer.SellNFT(goCtx, &types.MsgSellNFT{
		Creator: sellerAddress,
		ClassId: classId,
		NftId:   nftId,
		Buyer:   buyerAddress,
		Price:   finalPrice,
	})
	require.Error(t, err)
	require.Nil(t, res)
	require.Contains(t, err.Error(), types.ErrOfferExpired.Error())

	// Check state
	// Expect offer not removed
	_, found := k.GetOffer(ctx, classId, nftId, buyerAddressBytes)
	require.True(t, found)
	// Expect listing not removed
	_, found = k.GetListing(ctx, classId, nftId, sellerAddressBytes)
	require.True(t, found)

	ctrl.Finish()
}

// price too high
func TestSellNFTPriceTooHigh(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	accountKeeper := testutil.NewMockAccountKeeper(ctrl)
	bankKeeper := testutil.NewMockBankKeeper(ctrl)
	iscnKeeper := testutil.NewMockIscnKeeper(ctrl)
	nftKeeper := testutil.NewMockNftKeeper(ctrl)
	msgServer, goCtx, k := setupMsgServer(t, keeper.LikenftDependedKeepers{
		AccountKeeper: accountKeeper,
		BankKeeper:    bankKeeper,
		IscnKeeper:    iscnKeeper,
		NftKeeper:     nftKeeper,
	})
	ctx := sdk.UnwrapSDKContext(goCtx)
	ctx = ctx.WithBlockTime(time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC))
	goCtx = sdk.WrapSDKContext(ctx)

	// Data)
	sellerAddressBytes := []byte{0, 1, 0, 1, 0, 1, 0, 1}
	sellerAddress, _ := sdk.Bech32ifyAddressBytes("like", sellerAddressBytes)
	buyerAddressBytes := []byte{1, 0, 1, 0, 1, 0, 1, 0}
	buyerAddress, _ := sdk.Bech32ifyAddressBytes("like", buyerAddressBytes)
	classId := "likenft1abcdef"
	nftId := "nft1"
	price := uint64(123456)
	expiration := time.Date(2022, 4, 1, 0, 0, 0, 0, time.UTC)
	finalPrice := uint64(200000)

	// Seed listing to test deletion after txn
	k.SetListing(ctx, types.ListingStoreRecord{
		ClassId:    classId,
		NftId:      nftId,
		Seller:     sellerAddressBytes,
		Price:      uint64(999999),
		Expiration: time.Date(2022, 2, 1, 0, 0, 0, 0, time.UTC),
	})

	// Seed offer
	k.SetOffer(ctx, types.OfferStoreRecord{
		ClassId:    classId,
		NftId:      nftId,
		Buyer:      buyerAddressBytes,
		Price:      price,
		Expiration: expiration,
	})

	// Mock
	nftKeeper.EXPECT().GetOwner(gomock.Any(), classId, nftId).Return(sellerAddressBytes)

	// Call
	res, err := msgServer.SellNFT(goCtx, &types.MsgSellNFT{
		Creator: sellerAddress,
		ClassId: classId,
		NftId:   nftId,
		Buyer:   buyerAddress,
		Price:   finalPrice,
	})
	require.Error(t, err)
	require.Nil(t, res)
	require.Contains(t, err.Error(), types.ErrFailedToSellNFT.Error())

	// Check state
	// Expect offer not removed
	_, found := k.GetOffer(ctx, classId, nftId, buyerAddressBytes)
	require.True(t, found)
	// Expect listing not removed
	_, found = k.GetListing(ctx, classId, nftId, sellerAddressBytes)
	require.True(t, found)

	ctrl.Finish()
}
