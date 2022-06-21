package keeper_test

import (
	"testing"
	"time"

	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/golang/mock/gomock"
	"github.com/likecoin/likechain/backport/cosmos-sdk/v0.46.0-alpha2/x/nft"
	"github.com/likecoin/likechain/testutil/keeper"
	"github.com/likecoin/likechain/x/likenft/testutil"
	"github.com/likecoin/likechain/x/likenft/types"
	"github.com/stretchr/testify/require"
)

// normal royalty
func TestBuyNFTNormalRoyalty(t *testing.T) {
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
	creatorAddress, _ := sdk.Bech32ifyAddressBytes("like", creatorAddressBytes)
	sellerAddressBytes := []byte{0, 1, 0, 1, 0, 1, 0, 1}
	sellerAddress, _ := sdk.Bech32ifyAddressBytes("like", sellerAddressBytes)
	buyerAddressBytes := []byte{1, 0, 1, 0, 1, 0, 1, 0}
	buyerAddress, _ := sdk.Bech32ifyAddressBytes("like", buyerAddressBytes)
	classId := "likenft1abcdef"
	nftId := "nft1"
	price := uint64(123456)
	expiration := time.Date(2022, 4, 1, 0, 0, 0, 0, time.UTC)
	royaltyBasisPoints := uint64(234)
	finalPrice := uint64(200000)

	// Seed listing
	k.SetListing(ctx, types.ListingStoreRecord{
		ClassId:    classId,
		NftId:      nftId,
		Seller:     sellerAddressBytes,
		Price:      price,
		Expiration: expiration,
	})

	// Mock
	nftKeeper.EXPECT().GetOwner(gomock.Any(), classId, nftId).Return(sellerAddressBytes)
	bankKeeper.EXPECT().GetBalance(gomock.Any(), buyerAddressBytes, "nanolike").Return(sdk.NewCoin("nanolike", sdk.NewInt(1000000)))
	k.SetClassesByAccount(ctx, types.ClassesByAccount{
		Account:  creatorAddress,
		ClassIds: []string{classId},
	})
	classData := types.ClassData{
		Parent: types.ClassParent{
			Type:    types.ClassParentType_ACCOUNT,
			Account: creatorAddress,
		},
		Config: types.ClassConfig{
			RoyaltyBasisPoints: royaltyBasisPoints,
		},
	}
	classDataInAny, err := cdctypes.NewAnyWithValue(&classData)
	require.NoError(t, err)
	nftKeeper.EXPECT().GetClass(gomock.Any(), classId).Return(nft.Class{
		Id:   classId,
		Data: classDataInAny,
	}, true)
	royaltyAmount := finalPrice / 10000 * classData.Config.RoyaltyBasisPoints
	royaltyAmountCoins := sdk.NewCoins(sdk.NewCoin(k.GetParams(ctx).PriceDenom, sdk.NewInt(int64(royaltyAmount))))
	bankKeeper.EXPECT().SendCoins(gomock.Any(), buyerAddressBytes, creatorAddressBytes, royaltyAmountCoins).Return(nil)
	netAmount := finalPrice - royaltyAmount
	netAmountCoins := sdk.NewCoins(sdk.NewCoin(k.GetParams(ctx).PriceDenom, sdk.NewInt(int64(netAmount))))
	bankKeeper.EXPECT().SendCoins(gomock.Any(), buyerAddressBytes, sellerAddressBytes, netAmountCoins).Return(nil)
	nftKeeper.EXPECT().Transfer(gomock.Any(), classId, nftId, buyerAddressBytes).Return(nil)

	// Run
	res, err := msgServer.BuyNFT(goCtx, &types.MsgBuyNFT{
		Creator: buyerAddress,
		ClassId: classId,
		NftId:   nftId,
		Seller:  sellerAddress,
		Price:   finalPrice,
	})
	require.NoError(t, err)
	require.Equal(t, &types.MsgBuyNFTResponse{}, res)

	// Check state
	// Expect listing deleted
	_, found := k.GetListing(ctx, classId, nftId, sellerAddressBytes)
	require.False(t, found)

	ctrl.Finish()
}

// normal no royalty
func TestBuyNFTNormalNoRoyalty(t *testing.T) {
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
	creatorAddress, _ := sdk.Bech32ifyAddressBytes("like", creatorAddressBytes)
	sellerAddressBytes := []byte{0, 1, 0, 1, 0, 1, 0, 1}
	sellerAddress, _ := sdk.Bech32ifyAddressBytes("like", sellerAddressBytes)
	buyerAddressBytes := []byte{1, 0, 1, 0, 1, 0, 1, 0}
	buyerAddress, _ := sdk.Bech32ifyAddressBytes("like", buyerAddressBytes)
	classId := "likenft1abcdef"
	nftId := "nft1"
	price := uint64(123456)
	expiration := time.Date(2022, 4, 1, 0, 0, 0, 0, time.UTC)
	royaltyBasisPoints := uint64(0)
	finalPrice := uint64(200000)

	// Seed listing
	k.SetListing(ctx, types.ListingStoreRecord{
		ClassId:    classId,
		NftId:      nftId,
		Seller:     sellerAddressBytes,
		Price:      price,
		Expiration: expiration,
	})

	// Mock
	nftKeeper.EXPECT().GetOwner(gomock.Any(), classId, nftId).Return(sellerAddressBytes)
	bankKeeper.EXPECT().GetBalance(gomock.Any(), buyerAddressBytes, "nanolike").Return(sdk.NewCoin("nanolike", sdk.NewInt(1000000)))
	k.SetClassesByAccount(ctx, types.ClassesByAccount{
		Account:  creatorAddress,
		ClassIds: []string{classId},
	})
	classData := types.ClassData{
		Parent: types.ClassParent{
			Type:    types.ClassParentType_ACCOUNT,
			Account: creatorAddress,
		},
		Config: types.ClassConfig{
			RoyaltyBasisPoints: royaltyBasisPoints,
		},
	}
	classDataInAny, err := cdctypes.NewAnyWithValue(&classData)
	require.NoError(t, err)
	nftKeeper.EXPECT().GetClass(gomock.Any(), classId).Return(nft.Class{
		Id:   classId,
		Data: classDataInAny,
	}, true)
	netAmount := finalPrice
	netAmountCoins := sdk.NewCoins(sdk.NewCoin(k.GetParams(ctx).PriceDenom, sdk.NewInt(int64(netAmount))))
	bankKeeper.EXPECT().SendCoins(gomock.Any(), buyerAddressBytes, sellerAddressBytes, netAmountCoins).Return(nil)
	nftKeeper.EXPECT().Transfer(gomock.Any(), classId, nftId, buyerAddressBytes).Return(nil)

	// Run
	res, err := msgServer.BuyNFT(goCtx, &types.MsgBuyNFT{
		Creator: buyerAddress,
		ClassId: classId,
		NftId:   nftId,
		Seller:  sellerAddress,
		Price:   finalPrice,
	})
	require.NoError(t, err)
	require.Equal(t, &types.MsgBuyNFTResponse{}, res)

	// Check state
	// Expect listing deleted
	_, found := k.GetListing(ctx, classId, nftId, sellerAddressBytes)
	require.False(t, found)

	ctrl.Finish()
}

// listing not found
func TestBuyNFTListingNotFound(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	accountKeeper := testutil.NewMockAccountKeeper(ctrl)
	bankKeeper := testutil.NewMockBankKeeper(ctrl)
	iscnKeeper := testutil.NewMockIscnKeeper(ctrl)
	nftKeeper := testutil.NewMockNftKeeper(ctrl)
	msgServer, goCtx, _ := setupMsgServer(t, keeper.LikenftDependedKeepers{
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
	finalPrice := uint64(200000)

	// Run
	res, err := msgServer.BuyNFT(goCtx, &types.MsgBuyNFT{
		Creator: buyerAddress,
		ClassId: classId,
		NftId:   nftId,
		Seller:  sellerAddress,
		Price:   finalPrice,
	})
	require.Error(t, err)
	require.Nil(t, res)
	require.Contains(t, err.Error(), types.ErrListingNotFound.Error())

	ctrl.Finish()
}

// owner not valid
func TestBuyNFTListingOwnerInvalid(t *testing.T) {
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
	newOwnerAddressBytes := []byte{0, 0, 0, 0, 1, 1, 1, 1}
	buyerAddressBytes := []byte{1, 0, 1, 0, 1, 0, 1, 0}
	buyerAddress, _ := sdk.Bech32ifyAddressBytes("like", buyerAddressBytes)
	classId := "likenft1abcdef"
	nftId := "nft1"
	price := uint64(123456)
	expiration := time.Date(2022, 4, 1, 0, 0, 0, 0, time.UTC)
	finalPrice := uint64(200000)

	// Seed listing
	k.SetListing(ctx, types.ListingStoreRecord{
		ClassId:    classId,
		NftId:      nftId,
		Seller:     sellerAddressBytes,
		Price:      price,
		Expiration: expiration,
	})

	// Mock
	nftKeeper.EXPECT().GetOwner(gomock.Any(), classId, nftId).Return(newOwnerAddressBytes)

	// Run
	res, err := msgServer.BuyNFT(goCtx, &types.MsgBuyNFT{
		Creator: buyerAddress,
		ClassId: classId,
		NftId:   nftId,
		Seller:  sellerAddress,
		Price:   finalPrice,
	})
	require.Error(t, err)
	require.Nil(t, res)
	require.Contains(t, err.Error(), types.ErrListingExpired.Error())

	ctrl.Finish()
}

// listing expired
func TestBuyNFTListingExpired(t *testing.T) {
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
	expiration := time.Date(2021, 12, 31, 0, 0, 0, 0, time.UTC)
	finalPrice := uint64(200000)

	// Seed listing
	k.SetListing(ctx, types.ListingStoreRecord{
		ClassId:    classId,
		NftId:      nftId,
		Seller:     sellerAddressBytes,
		Price:      price,
		Expiration: expiration,
	})

	// Mock
	nftKeeper.EXPECT().GetOwner(gomock.Any(), classId, nftId).Return(sellerAddressBytes)

	// Run
	res, err := msgServer.BuyNFT(goCtx, &types.MsgBuyNFT{
		Creator: buyerAddress,
		ClassId: classId,
		NftId:   nftId,
		Seller:  sellerAddress,
		Price:   finalPrice,
	})
	require.Error(t, err)
	require.Nil(t, res)
	require.Contains(t, err.Error(), types.ErrListingExpired.Error())

	ctrl.Finish()
}

// price too low
func TestBuyNFTPriceTooLow(t *testing.T) {
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
	expiration := time.Date(2022, 3, 1, 0, 0, 0, 0, time.UTC)
	finalPrice := uint64(100000)

	// Seed listing
	k.SetListing(ctx, types.ListingStoreRecord{
		ClassId:    classId,
		NftId:      nftId,
		Seller:     sellerAddressBytes,
		Price:      price,
		Expiration: expiration,
	})

	// Mock
	nftKeeper.EXPECT().GetOwner(gomock.Any(), classId, nftId).Return(sellerAddressBytes)

	// Run
	res, err := msgServer.BuyNFT(goCtx, &types.MsgBuyNFT{
		Creator: buyerAddress,
		ClassId: classId,
		NftId:   nftId,
		Seller:  sellerAddress,
		Price:   finalPrice,
	})
	require.Error(t, err)
	require.Nil(t, res)
	require.Contains(t, err.Error(), types.ErrFailedToBuyNFT.Error())

	ctrl.Finish()
}

// not enough balance
func TestBuyNFTNotEnoughBalance(t *testing.T) {
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
	finalPrice := uint64(200000)

	// Seed listing
	k.SetListing(ctx, types.ListingStoreRecord{
		ClassId:    classId,
		NftId:      nftId,
		Seller:     sellerAddressBytes,
		Price:      price,
		Expiration: expiration,
	})

	// Mock
	nftKeeper.EXPECT().GetOwner(gomock.Any(), classId, nftId).Return(sellerAddressBytes)
	bankKeeper.EXPECT().GetBalance(gomock.Any(), buyerAddressBytes, "nanolike").Return(sdk.NewCoin("nanolike", sdk.NewInt(123456)))

	// Run
	res, err := msgServer.BuyNFT(goCtx, &types.MsgBuyNFT{
		Creator: buyerAddress,
		ClassId: classId,
		NftId:   nftId,
		Seller:  sellerAddress,
		Price:   finalPrice,
	})
	require.Error(t, err)
	require.Nil(t, res)
	require.Contains(t, err.Error(), types.ErrFailedToBuyNFT.Error())

	ctrl.Finish()
}
