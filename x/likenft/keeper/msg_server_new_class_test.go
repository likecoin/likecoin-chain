package keeper_test

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/golang/mock/gomock"
	"github.com/likecoin/likechain/testutil/keeper"
	iscntypes "github.com/likecoin/likechain/x/iscn/types"
	"github.com/likecoin/likechain/x/likenft/testutil"
	"github.com/likecoin/likechain/x/likenft/types"
	"github.com/stretchr/testify/require"
)

func TestNewClassISCNNormal(t *testing.T) {
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
	maxSupply := uint64(5)
	mintPeriods := []types.MintPeriod{
		{
			StartTime:        *testutil.MustParseTime(time.RFC3339, "2020-01-01T00:00:00Z"),
			AllowedAddresses: make([]string, 0),
			MintPrice:        1000000000,
		},
	}
	revealTime := *testutil.MustParseTime(time.RFC3339, "2322-04-20T00:00:00Z")

	// Mock keeper calls
	iscnLatestVersion := uint64(2)
	iscnKeeper.
		EXPECT().
		GetContentIdRecord(gomock.Any(), gomock.Eq(iscnId.Prefix)).
		Return(&iscntypes.ContentIdRecord{
			OwnerAddressBytes: ownerAddressBytes,
			LatestVersion:     iscnLatestVersion,
		})
	accountKeeper.EXPECT().GetAccount(gomock.Any(), ownerAddressBytes).Return(authtypes.NewBaseAccountWithAddress(ownerAddressBytes))
	bankKeeper.EXPECT().SendCoinsFromAccountToModule(gomock.Any(), ownerAddressBytes, authtypes.FeeCollectorName, gomock.Any()).Return(nil)

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
		Input: types.ClassInput{
			Name:        name,
			Symbol:      symbol,
			Description: description,
			Uri:         uri,
			UriHash:     uriHash,
			Metadata:    metadata,
			Config: types.ClassConfig{
				Burnable:  burnable,
				MaxSupply: maxSupply,
				BlindBoxConfig: &types.BlindBoxConfig{
					MintPeriods: mintPeriods,
					RevealTime:  revealTime,
				},
			},
		},
	})

	// Check output
	require.NoError(t, err)
	expectedClassId, _ := types.NewClassIdForISCN(iscnId.Prefix.String(), 0)
	require.Equal(t, expectedClassId, res.Class.Id)
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
	require.Equal(t, maxSupply, classData.Config.MaxSupply)
	require.Equal(t, revealTime, classData.Config.BlindBoxConfig.RevealTime)

	require.Equal(t, len(mintPeriods), len(classData.Config.BlindBoxConfig.MintPeriods))
	for i, mintPeriod := range classData.Config.BlindBoxConfig.MintPeriods {
		require.Equal(t, mintPeriod.StartTime, mintPeriods[i].StartTime)
		require.ElementsMatch(t, mintPeriod.AllowedAddresses, mintPeriods[i].AllowedAddresses)
		require.Equal(t, mintPeriod.MintPrice, mintPeriods[i].MintPrice)
	}

	// Check mock was called as expected
	ctrl.Finish()
}

func TestNewClassAccountNormal(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	accountKeeper := testutil.NewMockAccountKeeper(ctrl)
	bankKeeper := testutil.NewMockBankKeeper(ctrl)
	nftKeeper := testutil.NewMockNftKeeper(ctrl)
	msgServer, goCtx, _ := setupMsgServer(t, keeper.LikenftDependedKeepers{
		AccountKeeper: accountKeeper,
		BankKeeper:    bankKeeper,
		NftKeeper:     nftKeeper,
	})

	// Test Input
	ownerAddressBytes := []byte{0, 1, 0, 1, 0, 1, 0, 1}
	ownerAddress, _ := sdk.Bech32ifyAddressBytes("cosmos", ownerAddressBytes)
	ownerAccAddress, _ := sdk.AccAddressFromBech32(ownerAddress)

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
	maxSupply := uint64(5)

	accountKeeper.EXPECT().GetAccount(gomock.Any(), ownerAddressBytes).Return(authtypes.NewBaseAccountWithAddress(ownerAddressBytes))
	bankKeeper.EXPECT().SendCoinsFromAccountToModule(gomock.Any(), ownerAddressBytes, authtypes.FeeCollectorName, gomock.Any()).Return(nil)

	nftKeeper.
		EXPECT().
		SaveClass(gomock.Any(), gomock.Any()).
		Return(nil)

	// Run
	res, err := msgServer.NewClass(goCtx, &types.MsgNewClass{
		Creator: ownerAddress,
		Parent: types.ClassParentInput{
			Type: types.ClassParentType_ACCOUNT,
		},
		Input: types.ClassInput{
			Name:        name,
			Symbol:      symbol,
			Description: description,
			Uri:         uri,
			UriHash:     uriHash,
			Metadata:    metadata,
			Config: types.ClassConfig{
				Burnable:  burnable,
				MaxSupply: maxSupply,
			},
		},
	})

	// Check output
	require.NoError(t, err)
	expectedClassId, _ := types.NewClassIdForAccount(ownerAccAddress, 0)
	require.Equal(t, expectedClassId, res.Class.Id)
	require.Equal(t, name, res.Class.Name)
	require.Equal(t, symbol, res.Class.Symbol)
	require.Equal(t, description, res.Class.Description)
	require.Equal(t, uri, res.Class.Uri)
	require.Equal(t, uriHash, res.Class.UriHash)

	var classData types.ClassData
	err = classData.Unmarshal(res.Class.Data.Value)
	require.NoErrorf(t, err, "Error unmarshal class data")
	require.Equal(t, metadata, classData.Metadata)
	require.Equal(t, ownerAccAddress.String(), classData.Parent.Account)
	require.Equal(t, types.ClassConfig{
		Burnable:  burnable,
		MaxSupply: maxSupply,
	}, classData.Config)

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
	maxSupply := uint64(5)
	mintPeriods := []types.MintPeriod{
		{
			StartTime:        *testutil.MustParseTime(time.RFC3339, "2020-01-01T00:00:00Z"),
			AllowedAddresses: make([]string, 0),
			MintPrice:        1000000000,
		},
	}
	revealTime := *testutil.MustParseTime(time.RFC3339, "2322-04-20T00:00:00Z")

	// Run
	res, err := msgServer.NewClass(goCtx, &types.MsgNewClass{
		Creator: ownerAddress,
		Parent: types.ClassParentInput{
			Type:         types.ClassParentType_ISCN,
			IscnIdPrefix: iscnId,
		},
		Input: types.ClassInput{
			Name:        name,
			Symbol:      symbol,
			Description: description,
			Uri:         uri,
			UriHash:     uriHash,
			Metadata:    metadata,
			Config: types.ClassConfig{
				Burnable:  burnable,
				MaxSupply: maxSupply,
				BlindBoxConfig: &types.BlindBoxConfig{
					MintPeriods: mintPeriods,
					RevealTime:  revealTime,
				},
			},
		},
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
	maxSupply := uint64(5)
	mintPeriods := []types.MintPeriod{
		{
			StartTime:        *testutil.MustParseTime(time.RFC3339, "2020-01-01T00:00:00Z"),
			AllowedAddresses: make([]string, 0),
			MintPrice:        1000000000,
		},
	}
	revealTime := *testutil.MustParseTime(time.RFC3339, "2322-04-20T00:00:00Z")

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
		Input: types.ClassInput{
			Name:        name,
			Symbol:      symbol,
			Description: description,
			Uri:         uri,
			UriHash:     uriHash,
			Metadata:    metadata,
			Config: types.ClassConfig{
				Burnable:  burnable,
				MaxSupply: maxSupply,
				BlindBoxConfig: &types.BlindBoxConfig{
					MintPeriods: mintPeriods,
					RevealTime:  revealTime,
				},
			},
		},
	})

	// Check output
	require.Error(t, err)
	require.Contains(t, err.Error(), types.ErrIscnRecordNotFound.Error())
	require.Nil(t, res)

	// Check mock was called as expected
	ctrl.Finish()
}

func TestNewClassISCNInvalidUserAddress(t *testing.T) {
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
	maxSupply := uint64(5)
	mintPeriods := []types.MintPeriod{
		{
			StartTime:        *testutil.MustParseTime(time.RFC3339, "2020-01-01T00:00:00Z"),
			AllowedAddresses: make([]string, 0),
			MintPrice:        1000000000,
		},
	}
	revealTime := *testutil.MustParseTime(time.RFC3339, "2322-04-20T00:00:00Z")

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
		Input: types.ClassInput{
			Name:        name,
			Symbol:      symbol,
			Description: description,
			Uri:         uri,
			UriHash:     uriHash,
			Metadata:    metadata,
			Config: types.ClassConfig{
				Burnable:  burnable,
				MaxSupply: maxSupply,
				BlindBoxConfig: &types.BlindBoxConfig{
					MintPeriods: mintPeriods,
					RevealTime:  revealTime,
				},
			},
		},
	})

	// Check output
	require.Error(t, err)
	require.Contains(t, err.Error(), sdkerrors.ErrInvalidAddress.Error())
	require.Nil(t, res)

	// Check mock was called as expected
	ctrl.Finish()
}

func TestNewClassAccountInvalidUserAddress(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	accountKeeper := testutil.NewMockAccountKeeper(ctrl)
	bankKeeper := testutil.NewMockBankKeeper(ctrl)
	nftKeeper := testutil.NewMockNftKeeper(ctrl)
	msgServer, goCtx, _ := setupMsgServer(t, keeper.LikenftDependedKeepers{
		AccountKeeper: accountKeeper,
		BankKeeper:    bankKeeper,
		NftKeeper:     nftKeeper,
	})

	// Test Input
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
	maxSupply := uint64(5)

	// Run
	res, err := msgServer.NewClass(goCtx, &types.MsgNewClass{
		Creator: "invalid address",
		Parent: types.ClassParentInput{
			Type: types.ClassParentType_ACCOUNT,
		},
		Input: types.ClassInput{
			Name:        name,
			Symbol:      symbol,
			Description: description,
			Uri:         uri,
			UriHash:     uriHash,
			Metadata:    metadata,
			Config: types.ClassConfig{
				Burnable:  burnable,
				MaxSupply: maxSupply,
			},
		},
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
	maxSupply := uint64(5)
	mintPeriods := []types.MintPeriod{
		{
			StartTime:        *testutil.MustParseTime(time.RFC3339, "2020-01-01T00:00:00Z"),
			AllowedAddresses: make([]string, 0),
			MintPrice:        1000000000,
		},
	}
	revealTime := *testutil.MustParseTime(time.RFC3339, "2322-04-20T00:00:00Z")

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
		Input: types.ClassInput{
			Name:        name,
			Symbol:      symbol,
			Description: description,
			Uri:         uri,
			UriHash:     uriHash,
			Metadata:    metadata,
			Config: types.ClassConfig{
				Burnable:  burnable,
				MaxSupply: maxSupply,
				BlindBoxConfig: &types.BlindBoxConfig{
					MintPeriods: mintPeriods,
					RevealTime:  revealTime,
				},
			},
		},
	})

	// Check output
	require.Error(t, err)
	require.Contains(t, err.Error(), sdkerrors.ErrUnauthorized.Error())
	require.Nil(t, res)

	// Check mock was called as expected
	ctrl.Finish()
}

func TestNewClassNormalMintPeriodConfig(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	accountKeeper := testutil.NewMockAccountKeeper(ctrl)
	bankKeeper := testutil.NewMockBankKeeper(ctrl)
	iscnKeeper := testutil.NewMockIscnKeeper(ctrl)
	nftKeeper := testutil.NewMockNftKeeper(ctrl)
	msgServer, goCtx, keeper := setupMsgServer(t, keeper.LikenftDependedKeepers{
		AccountKeeper: accountKeeper,
		BankKeeper:    bankKeeper,
		IscnKeeper:    iscnKeeper,
		NftKeeper:     nftKeeper,
	})
	ctx := sdk.UnwrapSDKContext(goCtx)

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
	maxSupply := uint64(5)

	mintPeriods := []types.MintPeriod{
		{
			StartTime:        *testutil.MustParseTime(time.RFC3339, "2022-04-19T00:00:00Z"),
			AllowedAddresses: []string{ownerAddress},
			MintPrice:        uint64(20000),
		},
		{
			StartTime:        *testutil.MustParseTime(time.RFC3339, "2022-04-20T00:00:00Z"),
			AllowedAddresses: []string{ownerAddress},
			MintPrice:        uint64(30000),
		},
		{
			StartTime:        *testutil.MustParseTime(time.RFC3339, "2022-04-21T00:00:00Z"),
			AllowedAddresses: make([]string, 0),
			MintPrice:        uint64(90000),
		},
	}
	revealTime := *testutil.MustParseTime(time.RFC3339, "2022-04-28T00:00:00Z")

	// Mock keeper calls
	iscnLatestVersion := uint64(2)
	iscnKeeper.
		EXPECT().
		GetContentIdRecord(gomock.Any(), gomock.Eq(iscnId.Prefix)).
		Return(&iscntypes.ContentIdRecord{
			OwnerAddressBytes: ownerAddressBytes,
			LatestVersion:     iscnLatestVersion,
		})

	accountKeeper.EXPECT().GetAccount(gomock.Any(), ownerAddressBytes).Return(authtypes.NewBaseAccountWithAddress(ownerAddressBytes))
	bankKeeper.EXPECT().SendCoinsFromAccountToModule(gomock.Any(), ownerAddressBytes, authtypes.FeeCollectorName, gomock.Any()).Return(nil)

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
		Input: types.ClassInput{
			Name:        name,
			Symbol:      symbol,
			Description: description,
			Uri:         uri,
			UriHash:     uriHash,
			Metadata:    metadata,
			Config: types.ClassConfig{
				Burnable:  burnable,
				MaxSupply: maxSupply,
				BlindBoxConfig: &types.BlindBoxConfig{
					MintPeriods: mintPeriods,
					RevealTime:  revealTime,
				},
			},
		},
	})

	// Check output
	require.NoError(t, err)
	expectedClassId, _ := types.NewClassIdForISCN(iscnId.Prefix.String(), 0)
	require.Equal(t, expectedClassId, res.Class.Id)
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
	require.Equal(t, maxSupply, classData.Config.MaxSupply)
	require.Equal(t, revealTime, classData.Config.BlindBoxConfig.RevealTime)

	require.Equal(t, len(mintPeriods), len(classData.Config.BlindBoxConfig.MintPeriods))
	for i, mintPeriod := range classData.Config.BlindBoxConfig.MintPeriods {
		require.Equal(t, mintPeriod.StartTime, mintPeriods[i].StartTime)
		require.ElementsMatch(t, mintPeriod.AllowedAddresses, mintPeriods[i].AllowedAddresses)
		require.Equal(t, mintPeriod.MintPrice, mintPeriods[i].MintPrice)
	}

	revealQueue := keeper.GetClassRevealQueue(ctx)
	require.Contains(t, revealQueue, types.ClassRevealQueueEntry{
		ClassId:    expectedClassId,
		RevealTime: revealTime,
	})

	// Check mock was called as expected
	ctrl.Finish()
}

func TestNewClassInvalidMintPeriod(t *testing.T) {
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
	maxSupply := uint64(5)

	mintPeriods := []types.MintPeriod{
		{
			StartTime:        *testutil.MustParseTime(time.RFC3339, "2922-04-21T00:00:00Z"),
			AllowedAddresses: make([]string, 0),
			MintPrice:        1000000000,
		},
	}
	revealTime := *testutil.MustParseTime(time.RFC3339, "2322-04-20T00:00:00Z")

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
		Creator: ownerAddress,
		Parent: types.ClassParentInput{
			Type:         types.ClassParentType_ISCN,
			IscnIdPrefix: iscnId.Prefix.String(),
		},
		Input: types.ClassInput{
			Name:        name,
			Symbol:      symbol,
			Description: description,
			Uri:         uri,
			UriHash:     uriHash,
			Metadata:    metadata,
			Config: types.ClassConfig{
				Burnable:  burnable,
				MaxSupply: maxSupply,
				BlindBoxConfig: &types.BlindBoxConfig{
					MintPeriods: mintPeriods,
					RevealTime:  revealTime,
				},
			},
		},
	})

	// Check output
	require.Error(t, err)
	require.Contains(t, err.Error(), types.ErrInvalidNftClassConfig.Error())
	require.Nil(t, res)

	// Check mock was called as expected
	ctrl.Finish()
}

func TestNewClassMintPeriodInvalidAllowListAddress(t *testing.T) {
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
	maxSupply := uint64(5)

	mintPeriods := []types.MintPeriod{
		{
			StartTime:        *testutil.MustParseTime(time.RFC3339, "2022-04-19T00:00:00Z"),
			AllowedAddresses: []string{"invalid address"},
			MintPrice:        1000000000,
		},
	}
	revealTime := *testutil.MustParseTime(time.RFC3339, "2322-04-20T00:00:00Z")

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
		Creator: ownerAddress,
		Parent: types.ClassParentInput{
			Type:         types.ClassParentType_ISCN,
			IscnIdPrefix: iscnId.Prefix.String(),
		},
		Input: types.ClassInput{
			Name:        name,
			Symbol:      symbol,
			Description: description,
			Uri:         uri,
			UriHash:     uriHash,
			Metadata:    metadata,
			Config: types.ClassConfig{
				Burnable:  burnable,
				MaxSupply: maxSupply,
				BlindBoxConfig: &types.BlindBoxConfig{
					MintPeriods: mintPeriods,
					RevealTime:  revealTime,
				},
			},
		},
	})

	// Check output
	require.Error(t, err)
	require.Contains(t, err.Error(), sdkerrors.ErrInvalidAddress.Error())
	require.Nil(t, res)

	// Check mock was called as expected
	ctrl.Finish()
}

func TestNewClassBlindBoxNoMintPeriod(t *testing.T) {
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
	maxSupply := uint64(5)
	revealTime := *testutil.MustParseTime(time.RFC3339, "2322-04-20T00:00:00Z")

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
		Creator: ownerAddress,
		Parent: types.ClassParentInput{
			Type:         types.ClassParentType_ISCN,
			IscnIdPrefix: iscnId.Prefix.String(),
		},
		Input: types.ClassInput{
			Name:        name,
			Symbol:      symbol,
			Description: description,
			Uri:         uri,
			UriHash:     uriHash,
			Metadata:    metadata,
			Config: types.ClassConfig{
				Burnable:  burnable,
				MaxSupply: maxSupply,
				BlindBoxConfig: &types.BlindBoxConfig{
					MintPeriods: nil,
					RevealTime:  revealTime,
				},
			},
		},
	})

	// Check output
	require.Error(t, err)
	require.Contains(t, err.Error(), types.ErrInvalidNftClassConfig.Error())
	require.Nil(t, res)

	// Check mock was called as expected
	ctrl.Finish()
}
