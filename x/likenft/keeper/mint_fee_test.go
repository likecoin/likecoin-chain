package keeper_test

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/golang/mock/gomock"
	"github.com/likecoin/likecoin-chain/v3/testutil/keeper"
	"github.com/likecoin/likecoin-chain/v3/x/likenft/testutil"
	"github.com/likecoin/likecoin-chain/v3/x/likenft/types"
	"github.com/stretchr/testify/require"
)

type ICoinsMatcher struct {
	coins sdk.Coins
}

func coinsMatcher(coins sdk.Coins) gomock.Matcher {
	return ICoinsMatcher{
		coins,
	}
}

func (m ICoinsMatcher) Matches(x interface{}) bool {
	coins, ok := x.(sdk.Coins)
	if !ok {
		return false
	}

	return coins.IsEqual(m.coins)
}

func (m ICoinsMatcher) String() string {
	return fmt.Sprintf("data item coins is equal to %s", m.coins.String())
}

func TestMintFeeNormal(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	accountKeeper := testutil.NewMockAccountKeeper(ctrl)
	bankKeeper := testutil.NewMockBankKeeper(ctrl)
	iscnKeeper := testutil.NewMockIscnKeeper(ctrl)
	nftKeeper := testutil.NewMockNftKeeper(ctrl)
	_, goCtx, keeper := setupMsgServer(t, keeper.LikenftDependedKeepers{
		AccountKeeper: accountKeeper,
		BankKeeper:    bankKeeper,
		IscnKeeper:    iscnKeeper,
		NftKeeper:     nftKeeper,
	})
	ctx := sdk.UnwrapSDKContext(goCtx)
	feePerByte := sdk.NewDecCoin("nanoekil", sdk.NewInt(987654))
	params := types.DefaultParams()
	params.FeePerByte = feePerByte
	keeper.SetParams(ctx, params)

	ownerAddressBytes := []byte{0, 1, 0, 1, 0, 1, 0, 1}

	bytesLength := 123456

	// mock
	accountKeeper.EXPECT().GetAccount(gomock.Any(), ownerAddressBytes).Return(authtypes.NewBaseAccountWithAddress(ownerAddressBytes))

	expectedFee := feePerByte.Amount.MulInt64(int64(bytesLength))
	expectedCoins := sdk.NewCoins(sdk.NewCoin(feePerByte.Denom, expectedFee.Ceil().RoundInt()))

	bankKeeper.EXPECT().SendCoinsFromAccountToModule(gomock.Any(), ownerAddressBytes, authtypes.FeeCollectorName, coinsMatcher(expectedCoins)).Return(nil)

	// call
	err := keeper.DeductFeePerByte(ctx, ownerAddressBytes, bytesLength, nil)

	require.NoError(t, err)

	ctrl.Finish()
}

func TestMintFeeFreeForEmpty(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	accountKeeper := testutil.NewMockAccountKeeper(ctrl)
	bankKeeper := testutil.NewMockBankKeeper(ctrl)
	iscnKeeper := testutil.NewMockIscnKeeper(ctrl)
	nftKeeper := testutil.NewMockNftKeeper(ctrl)
	_, goCtx, keeper := setupMsgServer(t, keeper.LikenftDependedKeepers{
		AccountKeeper: accountKeeper,
		BankKeeper:    bankKeeper,
		IscnKeeper:    iscnKeeper,
		NftKeeper:     nftKeeper,
	})
	ctx := sdk.UnwrapSDKContext(goCtx)
	feePerByte := sdk.NewDecCoin("nanoekil", sdk.NewInt(987654))
	params := types.DefaultParams()
	params.FeePerByte = feePerByte
	keeper.SetParams(ctx, params)

	ownerAddressBytes := []byte{0, 1, 0, 1, 0, 1, 0, 1}

	bytesLength := 0

	// call
	err := keeper.DeductFeePerByte(ctx, ownerAddressBytes, bytesLength, nil)

	require.NoError(t, err)

	ctrl.Finish()
}

func TestMintFeeFreeForZeroRate(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	accountKeeper := testutil.NewMockAccountKeeper(ctrl)
	bankKeeper := testutil.NewMockBankKeeper(ctrl)
	iscnKeeper := testutil.NewMockIscnKeeper(ctrl)
	nftKeeper := testutil.NewMockNftKeeper(ctrl)
	_, goCtx, keeper := setupMsgServer(t, keeper.LikenftDependedKeepers{
		AccountKeeper: accountKeeper,
		BankKeeper:    bankKeeper,
		IscnKeeper:    iscnKeeper,
		NftKeeper:     nftKeeper,
	})
	ctx := sdk.UnwrapSDKContext(goCtx)
	feePerByte := sdk.NewDecCoin("nanoekil", sdk.NewInt(0))
	params := types.DefaultParams()
	params.FeePerByte = feePerByte
	keeper.SetParams(ctx, params)

	ownerAddressBytes := []byte{0, 1, 0, 1, 0, 1, 0, 1}

	bytesLength := 123456

	// call
	err := keeper.DeductFeePerByte(ctx, ownerAddressBytes, bytesLength, nil)

	require.NoError(t, err)

	ctrl.Finish()
}

func TestMintFeeAccNotFound(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	accountKeeper := testutil.NewMockAccountKeeper(ctrl)
	bankKeeper := testutil.NewMockBankKeeper(ctrl)
	iscnKeeper := testutil.NewMockIscnKeeper(ctrl)
	nftKeeper := testutil.NewMockNftKeeper(ctrl)
	_, goCtx, keeper := setupMsgServer(t, keeper.LikenftDependedKeepers{
		AccountKeeper: accountKeeper,
		BankKeeper:    bankKeeper,
		IscnKeeper:    iscnKeeper,
		NftKeeper:     nftKeeper,
	})
	ctx := sdk.UnwrapSDKContext(goCtx)
	feePerByte := sdk.NewDecCoin("nanoekil", sdk.NewInt(987654))
	params := types.DefaultParams()
	params.FeePerByte = feePerByte
	keeper.SetParams(ctx, params)

	ownerAddressBytes := []byte{0, 1, 0, 1, 0, 1, 0, 1}

	bytesLength := 123456

	// mock
	accountKeeper.EXPECT().GetAccount(gomock.Any(), ownerAddressBytes).Return(nil)

	// call
	err := keeper.DeductFeePerByte(ctx, ownerAddressBytes, bytesLength, nil)

	require.Error(t, err)

	ctrl.Finish()
}
