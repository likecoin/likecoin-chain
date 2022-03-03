package keeper_test

import (
	"testing"

	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/golang/mock/gomock"
	"github.com/likecoin/likechain/backport/cosmos-sdk/v0.46.0-alpha2/x/nft"
	"github.com/likecoin/likechain/testutil/keeper"
	iscntypes "github.com/likecoin/likechain/x/iscn/types"
	"github.com/likecoin/likechain/x/likenft/testutil"
	"github.com/likecoin/likechain/x/likenft/types"
	"github.com/stretchr/testify/require"
)

func TestUpdateClassNormal(t *testing.T) {
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
	classId := "likenft1aabbccddeeff"
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
	oldClassData := types.ClassData{
		Metadata:     types.JsonInput(`{"aaaa": "bbbb"}`),
		IscnIdPrefix: iscnId.Prefix.String(),
		Config: types.ClassConfig{
			Burnable: false,
		},
	}
	oldClassDataInAny, _ := cdctypes.NewAnyWithValue(&oldClassData)
	nftKeeper.
		EXPECT().
		GetClass(gomock.Any(), classId).
		Return(nft.Class{
			Id:          classId,
			Name:        "Old Name",
			Symbol:      "OLD",
			Description: "Old Class 234",
			Uri:         "ipfs://11223344",
			UriHash:     "11223344",
			Data:        oldClassDataInAny,
		}, true)

	nftKeeper.
		EXPECT().
		GetTotalSupply(gomock.Any(), classId).
		Return(uint64(0))

	keeper.SetClassesByISCN(ctx, types.ClassesByISCN{
		IscnIdPrefix: iscnId.Prefix.String(),
		ClassIds:     []string{classId},
	})

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
		UpdateClass(gomock.Any(), gomock.Any()).
		Return(nil)

	// Run
	res, err := msgServer.UpdateClass(goCtx, &types.MsgUpdateClass{
		Creator:     ownerAddress,
		ClassId:     classId,
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
	require.Equal(t, classId, res.Class.Id)
	require.Equal(t, name, res.Class.Name)
	require.Equal(t, symbol, res.Class.Symbol)
	require.Equal(t, description, res.Class.Description)
	require.Equal(t, uri, res.Class.Uri)
	require.Equal(t, uriHash, res.Class.UriHash)

	var classData types.ClassData
	err = classData.Unmarshal(res.Class.Data.Value)
	require.NoErrorf(t, err, "Error unmarshal class data")
	require.Equal(t, metadata, classData.Metadata)
	require.Equal(t, iscnId.Prefix.String(), classData.IscnIdPrefix)
	require.Equal(t, iscnLatestVersion, classData.IscnVersionAtMint)
	require.Equal(t, burnable, classData.Config.Burnable)

	// Check mock was called as expected
	ctrl.Finish()
}

func TestUpdateClassNotFound(t *testing.T) {
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
	classId := "likenft1aabbccddeeff"
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
	nftKeeper.
		EXPECT().
		GetClass(gomock.Any(), classId).
		Return(nft.Class{}, false)

	// Run
	res, err := msgServer.UpdateClass(goCtx, &types.MsgUpdateClass{
		Creator:     ownerAddress,
		ClassId:     classId,
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
	require.Contains(t, err.Error(), types.ErrNftClassNotFound.Error())
	require.Nil(t, res)

	// Check mock was called as expected
	ctrl.Finish()
}

func TestUpdateClassExistingMintedTokens(t *testing.T) {
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
	classId := "likenft1aabbccddeeff"
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
	oldClassData := types.ClassData{
		Metadata:     types.JsonInput(`{"aaaa": "bbbb"}`),
		IscnIdPrefix: iscnId.Prefix.String(),
		Config: types.ClassConfig{
			Burnable: false,
		},
	}
	oldClassDataInAny, _ := cdctypes.NewAnyWithValue(&oldClassData)
	nftKeeper.
		EXPECT().
		GetClass(gomock.Any(), classId).
		Return(nft.Class{
			Id:          classId,
			Name:        "Old Name",
			Symbol:      "OLD",
			Description: "Old Class 234",
			Uri:         "ipfs://11223344",
			UriHash:     "11223344",
			Data:        oldClassDataInAny,
		}, true)

	nftKeeper.
		EXPECT().
		GetTotalSupply(gomock.Any(), classId).
		Return(uint64(1))

	// Run
	res, err := msgServer.UpdateClass(goCtx, &types.MsgUpdateClass{
		Creator:     ownerAddress,
		ClassId:     classId,
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
	require.Contains(t, err.Error(), types.ErrCannotUpdateClassWithMintedTokens.Error())
	require.Nil(t, res)

	// Check mock was called as expected
	ctrl.Finish()
}

func TestUpdateClassMissingRelationRecord(t *testing.T) {
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
	classId := "likenft1aabbccddeeff"
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
	oldClassData := types.ClassData{
		Metadata:     types.JsonInput(`{"aaaa": "bbbb"}`),
		IscnIdPrefix: iscnId.Prefix.String(),
		Config: types.ClassConfig{
			Burnable: false,
		},
	}
	oldClassDataInAny, _ := cdctypes.NewAnyWithValue(&oldClassData)
	nftKeeper.
		EXPECT().
		GetClass(gomock.Any(), classId).
		Return(nft.Class{
			Id:          classId,
			Name:        "Old Name",
			Symbol:      "OLD",
			Description: "Old Class 234",
			Uri:         "ipfs://11223344",
			UriHash:     "11223344",
			Data:        oldClassDataInAny,
		}, true)

	nftKeeper.
		EXPECT().
		GetTotalSupply(gomock.Any(), classId).
		Return(uint64(0))

	// Run
	res, err := msgServer.UpdateClass(goCtx, &types.MsgUpdateClass{
		Creator:     ownerAddress,
		ClassId:     classId,
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
	require.Contains(t, err.Error(), types.ErrNftClassNotRelatedToAnyIscn.Error())
	require.Nil(t, res)

	// Check mock was called as expected
	ctrl.Finish()
}

func TestUpdateClassNotRelatedToIscn(t *testing.T) {
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
	classId := "likenft1aabbccddeeff"
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
	oldClassData := types.ClassData{
		Metadata:     types.JsonInput(`{"aaaa": "bbbb"}`),
		IscnIdPrefix: iscnId.Prefix.String(),
		Config: types.ClassConfig{
			Burnable: false,
		},
	}
	oldClassDataInAny, _ := cdctypes.NewAnyWithValue(&oldClassData)
	nftKeeper.
		EXPECT().
		GetClass(gomock.Any(), classId).
		Return(nft.Class{
			Id:          classId,
			Name:        "Old Name",
			Symbol:      "OLD",
			Description: "Old Class 234",
			Uri:         "ipfs://11223344",
			UriHash:     "11223344",
			Data:        oldClassDataInAny,
		}, true)

	nftKeeper.
		EXPECT().
		GetTotalSupply(gomock.Any(), classId).
		Return(uint64(0))

	keeper.SetClassesByISCN(ctx, types.ClassesByISCN{
		IscnIdPrefix: iscnId.Prefix.String(),
		ClassIds:     []string{"likenft199887766"},
	})

	// Run
	res, err := msgServer.UpdateClass(goCtx, &types.MsgUpdateClass{
		Creator:     ownerAddress,
		ClassId:     classId,
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
	require.Contains(t, err.Error(), types.ErrNftClassNotRelatedToAnyIscn.Error())
	require.Nil(t, res)

	// Check mock was called as expected
	ctrl.Finish()
}

func TestUpdateClassIscnNotFound(t *testing.T) {
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
	classId := "likenft1aabbccddeeff"
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
	oldClassData := types.ClassData{
		Metadata:     types.JsonInput(`{"aaaa": "bbbb"}`),
		IscnIdPrefix: iscnId.Prefix.String(),
		Config: types.ClassConfig{
			Burnable: false,
		},
	}
	oldClassDataInAny, _ := cdctypes.NewAnyWithValue(&oldClassData)
	nftKeeper.
		EXPECT().
		GetClass(gomock.Any(), classId).
		Return(nft.Class{
			Id:          classId,
			Name:        "Old Name",
			Symbol:      "OLD",
			Description: "Old Class 234",
			Uri:         "ipfs://11223344",
			UriHash:     "11223344",
			Data:        oldClassDataInAny,
		}, true)

	nftKeeper.
		EXPECT().
		GetTotalSupply(gomock.Any(), classId).
		Return(uint64(0))

	keeper.SetClassesByISCN(ctx, types.ClassesByISCN{
		IscnIdPrefix: iscnId.Prefix.String(),
		ClassIds:     []string{classId},
	})

	iscnKeeper.
		EXPECT().
		GetContentIdRecord(gomock.Any(), gomock.Eq(iscnId.Prefix)).
		Return(nil)

	// Run
	res, err := msgServer.UpdateClass(goCtx, &types.MsgUpdateClass{
		Creator:     ownerAddress,
		ClassId:     classId,
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

func TestUpdateClassInvalidUserAddress(t *testing.T) {
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
	classId := "likenft1aabbccddeeff"
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
	oldClassData := types.ClassData{
		Metadata:     types.JsonInput(`{"aaaa": "bbbb"}`),
		IscnIdPrefix: iscnId.Prefix.String(),
		Config: types.ClassConfig{
			Burnable: false,
		},
	}
	oldClassDataInAny, _ := cdctypes.NewAnyWithValue(&oldClassData)
	nftKeeper.
		EXPECT().
		GetClass(gomock.Any(), classId).
		Return(nft.Class{
			Id:          classId,
			Name:        "Old Name",
			Symbol:      "OLD",
			Description: "Old Class 234",
			Uri:         "ipfs://11223344",
			UriHash:     "11223344",
			Data:        oldClassDataInAny,
		}, true)

	nftKeeper.
		EXPECT().
		GetTotalSupply(gomock.Any(), classId).
		Return(uint64(0))

	keeper.SetClassesByISCN(ctx, types.ClassesByISCN{
		IscnIdPrefix: iscnId.Prefix.String(),
		ClassIds:     []string{classId},
	})

	iscnKeeper.
		EXPECT().
		GetContentIdRecord(gomock.Any(), gomock.Eq(iscnId.Prefix)).
		Return(&iscntypes.ContentIdRecord{
			OwnerAddressBytes: ownerAddressBytes,
			LatestVersion:     1,
		})

	// Run
	res, err := msgServer.UpdateClass(goCtx, &types.MsgUpdateClass{
		Creator:     "not an address",
		ClassId:     classId,
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

func TestUpdateClassUserNotOwner(t *testing.T) {
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
	classId := "likenft1aabbccddeeff"
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
	oldClassData := types.ClassData{
		Metadata:     types.JsonInput(`{"aaaa": "bbbb"}`),
		IscnIdPrefix: iscnId.Prefix.String(),
		Config: types.ClassConfig{
			Burnable: false,
		},
	}
	oldClassDataInAny, _ := cdctypes.NewAnyWithValue(&oldClassData)
	nftKeeper.
		EXPECT().
		GetClass(gomock.Any(), classId).
		Return(nft.Class{
			Id:          classId,
			Name:        "Old Name",
			Symbol:      "OLD",
			Description: "Old Class 234",
			Uri:         "ipfs://11223344",
			UriHash:     "11223344",
			Data:        oldClassDataInAny,
		}, true)

	nftKeeper.
		EXPECT().
		GetTotalSupply(gomock.Any(), classId).
		Return(uint64(0))

	keeper.SetClassesByISCN(ctx, types.ClassesByISCN{
		IscnIdPrefix: iscnId.Prefix.String(),
		ClassIds:     []string{classId},
	})

	notOwnerAddressBytes := []byte{1, 1, 1, 1, 1, 1, 1, 1}
	iscnKeeper.
		EXPECT().
		GetContentIdRecord(gomock.Any(), gomock.Eq(iscnId.Prefix)).
		Return(&iscntypes.ContentIdRecord{
			OwnerAddressBytes: notOwnerAddressBytes,
			LatestVersion:     1,
		})

	// Run
	res, err := msgServer.UpdateClass(goCtx, &types.MsgUpdateClass{
		Creator:     ownerAddress,
		ClassId:     classId,
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
