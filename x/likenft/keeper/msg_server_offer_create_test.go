package keeper_test

import (
	"fmt"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/golang/mock/gomock"
	"github.com/likecoin/likecoin-chain/v3/testutil/keeper"
	"github.com/likecoin/likecoin-chain/v3/x/likenft/testutil"
	"github.com/likecoin/likecoin-chain/v3/x/likenft/types"
	"github.com/stretchr/testify/require"
)

// Test normal
func TestCreateOfferNormal(t *testing.T) {
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
	userAddressBytes := []byte{0, 1, 0, 1, 0, 1, 0, 1}
	userAddress, _ := sdk.Bech32ifyAddressBytes("like", userAddressBytes)
	classId := "likenft1abcdef"
	nftId := "nft1"
	price := uint64(123456)
	expiration := time.Date(2022, 4, 1, 0, 0, 0, 0, time.UTC)

	// Mock
	nftKeeper.EXPECT().HasNFT(gomock.Any(), classId, nftId).Return(true)
	bankKeeper.EXPECT().GetBalance(gomock.Any(), userAddressBytes, "nanolike").Return(sdk.NewCoin("nanolike", sdk.NewInt(1000000)))
	bankKeeper.EXPECT().SendCoinsFromAccountToModule(gomock.Any(), userAddressBytes, types.ModuleName, sdk.NewCoins(sdk.NewCoin("nanolike", sdk.NewInt(int64(price))))).Return(nil)

	// Call
	res, err := msgServer.CreateOffer(goCtx, &types.MsgCreateOffer{
		Creator:    userAddress,
		ClassId:    classId,
		NftId:      nftId,
		Price:      price,
		Expiration: expiration,
	})
	require.NoError(t, err)
	expectedOffer := types.Offer{
		ClassId:    classId,
		NftId:      nftId,
		Buyer:      userAddress,
		Price:      price,
		Expiration: expiration,
	}
	require.Equal(t, &types.MsgCreateOfferResponse{
		Offer: expectedOffer,
	}, res)

	// Check state
	// expect new offer
	offer, isFound := k.GetOffer(
		ctx,
		classId,
		nftId,
		userAddressBytes,
	)
	require.True(t, isFound)
	require.Equal(t, expectedOffer.ToStoreRecord(), offer)
	// expect enqueued offer
	_, isFound = k.GetOfferExpireQueueEntry(
		ctx,
		expiration,
		types.OfferKey(classId, nftId, userAddressBytes),
	)
	require.True(t, isFound)

	ctrl.Finish()
}

// Already exists
func TestCreateOfferAlreadyExists(t *testing.T) {
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
	userAddressBytes := []byte{0, 1, 0, 1, 0, 1, 0, 1}
	userAddress, _ := sdk.Bech32ifyAddressBytes("like", userAddressBytes)
	classId := "likenft1abcdef"
	nftId := "nft1"
	price := uint64(123456)
	expiration := time.Date(2022, 4, 1, 0, 0, 0, 0, time.UTC)

	// Seed existing offer
	k.SetOffer(ctx, types.OfferStoreRecord{
		ClassId:    classId,
		NftId:      nftId,
		Buyer:      userAddressBytes,
		Price:      price,
		Expiration: expiration,
	})

	// Call
	res, err := msgServer.CreateOffer(goCtx, &types.MsgCreateOffer{
		Creator:    userAddress,
		ClassId:    classId,
		NftId:      nftId,
		Price:      price,
		Expiration: expiration,
	})
	require.Error(t, err)
	require.Nil(t, res)
	require.Contains(t, err.Error(), types.ErrOfferAlreadyExists.Error())

	ctrl.Finish()
}

// NFT not exist
func TestCreateOfferNFTNotExist(t *testing.T) {
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
	userAddressBytes := []byte{0, 1, 0, 1, 0, 1, 0, 1}
	userAddress, _ := sdk.Bech32ifyAddressBytes("like", userAddressBytes)
	classId := "likenft1abcdef"
	nftId := "nft1"
	price := uint64(123456)
	expiration := time.Date(2022, 4, 1, 0, 0, 0, 0, time.UTC)

	// Mock
	nftKeeper.EXPECT().HasNFT(gomock.Any(), classId, nftId).Return(false)

	// Call
	res, err := msgServer.CreateOffer(goCtx, &types.MsgCreateOffer{
		Creator:    userAddress,
		ClassId:    classId,
		NftId:      nftId,
		Price:      price,
		Expiration: expiration,
	})
	require.Error(t, err)
	require.Nil(t, res)
	require.Contains(t, err.Error(), types.ErrNftNotFound.Error())

	ctrl.Finish()
}

// expiration out of range
func TestCreateOfferExpirationInPast(t *testing.T) {
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
	userAddressBytes := []byte{0, 1, 0, 1, 0, 1, 0, 1}
	userAddress, _ := sdk.Bech32ifyAddressBytes("like", userAddressBytes)
	classId := "likenft1abcdef"
	nftId := "nft1"
	price := uint64(123456)
	expiration := time.Date(2019, 12, 31, 0, 0, 0, 0, time.UTC)

	// Mock
	nftKeeper.EXPECT().HasNFT(gomock.Any(), classId, nftId).Return(true)

	// Call
	res, err := msgServer.CreateOffer(goCtx, &types.MsgCreateOffer{
		Creator:    userAddress,
		ClassId:    classId,
		NftId:      nftId,
		Price:      price,
		Expiration: expiration,
	})
	require.Error(t, err)
	require.Nil(t, res)
	require.Contains(t, err.Error(), sdkerrors.ErrInvalidRequest.Error())

	ctrl.Finish()
}

func TestCreateOfferExpirationTooFarInFuture(t *testing.T) {
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
	userAddressBytes := []byte{0, 1, 0, 1, 0, 1, 0, 1}
	userAddress, _ := sdk.Bech32ifyAddressBytes("like", userAddressBytes)
	classId := "likenft1abcdef"
	nftId := "nft1"
	price := uint64(123456)
	expiration := time.Date(2022, 7, 1, 0, 0, 0, 0, time.UTC)

	// Mock
	nftKeeper.EXPECT().HasNFT(gomock.Any(), classId, nftId).Return(true)

	// Call
	res, err := msgServer.CreateOffer(goCtx, &types.MsgCreateOffer{
		Creator:    userAddress,
		ClassId:    classId,
		NftId:      nftId,
		Price:      price,
		Expiration: expiration,
	})
	require.Error(t, err)
	require.Nil(t, res)
	require.Contains(t, err.Error(), sdkerrors.ErrInvalidRequest.Error())

	ctrl.Finish()
}

// not enough balance
func TestCreateOfferNotEnoughBalance(t *testing.T) {
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
	userAddressBytes := []byte{0, 1, 0, 1, 0, 1, 0, 1}
	userAddress, _ := sdk.Bech32ifyAddressBytes("like", userAddressBytes)
	classId := "likenft1abcdef"
	nftId := "nft1"
	price := uint64(123456)
	expiration := time.Date(2022, 4, 1, 0, 0, 0, 0, time.UTC)

	// Mock
	nftKeeper.EXPECT().HasNFT(gomock.Any(), classId, nftId).Return(true)
	bankKeeper.EXPECT().GetBalance(gomock.Any(), userAddressBytes, "nanolike").Return(sdk.NewCoin("nanolike", sdk.NewInt(123455)))

	// Call
	res, err := msgServer.CreateOffer(goCtx, &types.MsgCreateOffer{
		Creator:    userAddress,
		ClassId:    classId,
		NftId:      nftId,
		Price:      price,
		Expiration: expiration,
	})
	require.Error(t, err)
	require.Nil(t, res)
	require.Contains(t, err.Error(), types.ErrInsufficientFunds.Error())

	ctrl.Finish()
}

// payment error
func TestCreateOfferPaymentError(t *testing.T) {
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
	userAddressBytes := []byte{0, 1, 0, 1, 0, 1, 0, 1}
	userAddress, _ := sdk.Bech32ifyAddressBytes("like", userAddressBytes)
	classId := "likenft1abcdef"
	nftId := "nft1"
	price := uint64(123456)
	expiration := time.Date(2022, 4, 1, 0, 0, 0, 0, time.UTC)

	// Mock
	nftKeeper.EXPECT().HasNFT(gomock.Any(), classId, nftId).Return(true)
	bankKeeper.EXPECT().GetBalance(gomock.Any(), userAddressBytes, "nanolike").Return(sdk.NewCoin("nanolike", sdk.NewInt(1000000)))
	bankKeeper.EXPECT().SendCoinsFromAccountToModule(gomock.Any(), userAddressBytes, types.ModuleName, sdk.NewCoins(sdk.NewCoin("nanolike", sdk.NewInt(int64(price))))).Return(fmt.Errorf("some error"))

	// Call
	res, err := msgServer.CreateOffer(goCtx, &types.MsgCreateOffer{
		Creator:    userAddress,
		ClassId:    classId,
		NftId:      nftId,
		Price:      price,
		Expiration: expiration,
	})
	require.Error(t, err)
	require.Nil(t, res)
	require.Contains(t, err.Error(), types.ErrFailedToCreateOffer.Error())

	ctrl.Finish()
}
