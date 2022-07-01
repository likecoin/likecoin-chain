package testutil

import (
	"testing"

	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/golang/mock/gomock"
	"github.com/likecoin/likecoin-chain/v3/backport/cosmos-sdk/v0.46.0-rc1/x/nft"
	keepertest "github.com/likecoin/likecoin-chain/v3/testutil/keeper"
	"github.com/likecoin/likecoin-chain/v3/x/likenft/keeper"
	"github.com/likecoin/likecoin-chain/v3/x/likenft/types"
)

func LikenftKeeperForBlindBoxTest(t *testing.T) (*keeper.Keeper, sdk.Context, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	accountKeeper := NewMockAccountKeeper(ctrl)
	bankKeeper := NewMockBankKeeper(ctrl)
	iscnKeeper := NewMockIscnKeeper(ctrl)
	nftKeeper := NewMockNftKeeper(ctrl)
	keeper, ctx := keepertest.LikenftKeeperOverrideDependedKeepers(t, keepertest.LikenftDependedKeepers{
		AccountKeeper: accountKeeper,
		BankKeeper:    bankKeeper,
		IscnKeeper:    iscnKeeper,
		NftKeeper:     nftKeeper,
	})
	classData := types.ClassData{
		BlindBoxState: types.BlindBoxState{
			ContentCount: 0,
		},
	}
	classDataInAny, _ := cdctypes.NewAnyWithValue(&classData)
	nftKeeper.
		EXPECT().
		GetClass(gomock.Any(), gomock.Any()).
		Return(nft.Class{
			Data: classDataInAny,
		}, true).
		AnyTimes()
	nftKeeper.
		EXPECT().
		UpdateClass(gomock.Any(), gomock.Any()).
		AnyTimes()

	return keeper, ctx, ctrl
}
