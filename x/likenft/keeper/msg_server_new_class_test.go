package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/golang/mock/gomock"
	"github.com/likecoin/likechain/testutil/keeper"
	iscntypes "github.com/likecoin/likechain/x/iscn/types"
	"github.com/likecoin/likechain/x/likenft/testutil"
	"github.com/likecoin/likechain/x/likenft/types"
	"github.com/stretchr/testify/require"
)

func TestNewClassNormal(t *testing.T) {
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

	// Test Input
	ownerAddressBytes := []byte{0, 1, 0, 1, 0, 1, 0, 1}
	ownerAddress, _ := sdk.Bech32ifyAddressBytes("cosmos", ownerAddressBytes)
	iscnId := iscntypes.NewIscnId("likecoin-chain", "abcdef", 1)
	name := "Class Name"
	symbol := "ABC"
	description := "Testing Class 123"
	uri := "ipfs://abcdef"
	uriHash := "abcdef"
	metadata := types.JsonInput(
		`{
	"abc": "def",
	"qwerty": 1234,
	"bool": false,
	"null": null,
	"nested": {
		"object": {
			"abc": "def"
		}
	}
}`)
	burnable := true

	// Mock keeper calls
	iscnLatestVersion := uint64(2)
	iscnKeeper.
		EXPECT().
		GetContentIdRecord(gomock.Any(), gomock.Eq(iscnId.Prefix)).
		Return(&iscntypes.ContentIdRecord{
			OwnerAddressBytes: ownerAddressBytes,
			LatestVersion:     iscnLatestVersion,
		})

	nftKeeper.
		EXPECT().
		SaveClass(gomock.Any(), gomock.Any()).
		Return(nil)

	// Run
	res, err := msgServer.NewClass(goCtx, &types.MsgNewClass{
		Creator: ownerAddress,
		Parent: types.ClassParentInput{
			Type:         types.ClassParentType_ISCN,
			IscnIdPrefix: iscnId.Prefix.String(),
		},
		Name:        name,
		Symbol:      symbol,
		Description: description,
		Uri:         uri,
		UriHash:     uriHash,
		Metadata:    metadata,
		Burnable:    burnable,
	})

	// Check output
	require.NoError(t, err)
	expectedClassId, _ := types.NewClassId(iscnId.Prefix.String(), 0)
	require.Equal(t, *expectedClassId, res.Class.Id)
	require.Equal(t, name, res.Class.Name)
	require.Equal(t, symbol, res.Class.Symbol)
	require.Equal(t, description, res.Class.Description)
	require.Equal(t, uri, res.Class.Uri)
	require.Equal(t, uriHash, res.Class.UriHash)

	var classData types.ClassData
	err = classData.Unmarshal(res.Class.Data.Value)
	require.NoErrorf(t, err, "Error unmarshal class data")
	require.Equal(t, metadata, classData.Metadata)
	require.Equal(t, iscnId.Prefix.String(), classData.Parent.IscnIdPrefix)
	require.Equal(t, iscnLatestVersion, classData.Parent.IscnVersionAtMint)
	require.Equal(t, burnable, classData.Config.Burnable)

	// Check mock was called as expected
	ctrl.Finish()
}

func TestNewClassInvalidIscn(t *testing.T) {
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

	// Test Input
	ownerAddressBytes := []byte{0, 1, 0, 1, 0, 1, 0, 1}
	ownerAddress, _ := sdk.Bech32ifyAddressBytes("cosmos", ownerAddressBytes)
	iscnId := "not an iscn id"
	name := "Class Name"
	symbol := "ABC"
	description := "Testing Class 123"
	uri := "ipfs://abcdef"
	uriHash := "abcdef"
	metadata := types.JsonInput(
		`{
	"abc": "def",
	"qwerty": 1234,
	"bool": false,
	"null": null,
	"nested": {
		"object": {
			"abc": "def"
		}
	}
}`)
	burnable := true

	// Run
	res, err := msgServer.NewClass(goCtx, &types.MsgNewClass{
		Creator: ownerAddress,
		Parent: types.ClassParentInput{
			Type:         types.ClassParentType_ISCN,
			IscnIdPrefix: iscnId,
		},
		Name:        name,
		Symbol:      symbol,
		Description: description,
		Uri:         uri,
		UriHash:     uriHash,
		Metadata:    metadata,
		Burnable:    burnable,
	})

	// Check output
	require.Error(t, err)
	require.Contains(t, err.Error(), types.ErrInvalidIscnId.Error())
	require.Nil(t, res)

	// Check mock was called as expected
	ctrl.Finish()
}

func TestNewClassNonExistentIscn(t *testing.T) {
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

	// Test Input
	ownerAddressBytes := []byte{0, 1, 0, 1, 0, 1, 0, 1}
	ownerAddress, _ := sdk.Bech32ifyAddressBytes("cosmos", ownerAddressBytes)
	iscnId := iscntypes.NewIscnId("likecoin-chain", "abcdef", 1)
	name := "Class Name"
	symbol := "ABC"
	description := "Testing Class 123"
	uri := "ipfs://abcdef"
	uriHash := "abcdef"
	metadata := types.JsonInput(
		`{
	"abc": "def",
	"qwerty": 1234,
	"bool": false,
	"null": null,
	"nested": {
		"object": {
			"abc": "def"
		}
	}
}`)
	burnable := true

	// Mock keeper calls
	iscnKeeper.
		EXPECT().
		GetContentIdRecord(gomock.Any(), gomock.Eq(iscnId.Prefix)).
		Return(nil)

	// Run
	res, err := msgServer.NewClass(goCtx, &types.MsgNewClass{
		Creator: ownerAddress,
		Parent: types.ClassParentInput{
			Type:         types.ClassParentType_ISCN,
			IscnIdPrefix: iscnId.Prefix.String(),
		},
		Name:        name,
		Symbol:      symbol,
		Description: description,
		Uri:         uri,
		UriHash:     uriHash,
		Metadata:    metadata,
		Burnable:    burnable,
	})

	// Check output
	require.Error(t, err)
	require.Contains(t, err.Error(), types.ErrIscnRecordNotFound.Error())
	require.Nil(t, res)

	// Check mock was called as expected
	ctrl.Finish()
}

func TestNewClassInvalidUserAddress(t *testing.T) {
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

	// Test Input
	ownerAddressBytes := []byte{0, 1, 0, 1, 0, 1, 0, 1}
	iscnId := iscntypes.NewIscnId("likecoin-chain", "abcdef", 1)
	name := "Class Name"
	symbol := "ABC"
	description := "Testing Class 123"
	uri := "ipfs://abcdef"
	uriHash := "abcdef"
	metadata := types.JsonInput(
		`{
	"abc": "def",
	"qwerty": 1234,
	"bool": false,
	"null": null,
	"nested": {
		"object": {
			"abc": "def"
		}
	}
}`)
	burnable := true

	// Mock keeper calls
	iscnKeeper.
		EXPECT().
		GetContentIdRecord(gomock.Any(), gomock.Eq(iscnId.Prefix)).
		Return(&iscntypes.ContentIdRecord{
			OwnerAddressBytes: ownerAddressBytes,
			LatestVersion:     1,
		})

	// Run
	res, err := msgServer.NewClass(goCtx, &types.MsgNewClass{
		Creator: "invalid address",
		Parent: types.ClassParentInput{
			Type:         types.ClassParentType_ISCN,
			IscnIdPrefix: iscnId.Prefix.String(),
		},
		Name:        name,
		Symbol:      symbol,
		Description: description,
		Uri:         uri,
		UriHash:     uriHash,
		Metadata:    metadata,
		Burnable:    burnable,
	})

	// Check output
	require.Error(t, err)
	require.Contains(t, err.Error(), sdkerrors.ErrInvalidAddress.Error())
	require.Nil(t, res)

	// Check mock was called as expected
	ctrl.Finish()
}

func TestNewClassUserNotIscnOwner(t *testing.T) {
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

	// Test Input
	ownerAddressBytes := []byte{0, 1, 0, 1, 0, 1, 0, 1}
	ownerAddress, _ := sdk.Bech32ifyAddressBytes("cosmos", ownerAddressBytes)
	iscnId := iscntypes.NewIscnId("likecoin-chain", "abcdef", 1)
	name := "Class Name"
	symbol := "ABC"
	description := "Testing Class 123"
	uri := "ipfs://abcdef"
	uriHash := "abcdef"
	metadata := types.JsonInput(
		`{
	"abc": "def",
	"qwerty": 1234,
	"bool": false,
	"null": null,
	"nested": {
		"object": {
			"abc": "def"
		}
	}
}`)
	burnable := true

	// Mock keeper calls
	notOwnerAddressBytes := []byte{1, 1, 1, 1, 1, 1, 1, 1}
	iscnKeeper.
		EXPECT().
		GetContentIdRecord(gomock.Any(), gomock.Eq(iscnId.Prefix)).
		Return(&iscntypes.ContentIdRecord{
			OwnerAddressBytes: notOwnerAddressBytes,
			LatestVersion:     1,
		})

	// Run
	res, err := msgServer.NewClass(goCtx, &types.MsgNewClass{
		Creator: ownerAddress,
		Parent: types.ClassParentInput{
			Type:         types.ClassParentType_ISCN,
			IscnIdPrefix: iscnId.Prefix.String(),
		},
		Name:        name,
		Symbol:      symbol,
		Description: description,
		Uri:         uri,
		UriHash:     uriHash,
		Metadata:    metadata,
		Burnable:    burnable,
	})

	// Check output
	require.Error(t, err)
	require.Contains(t, err.Error(), sdkerrors.ErrUnauthorized.Error())
	require.Nil(t, res)

	// Check mock was called as expected
	ctrl.Finish()
}
