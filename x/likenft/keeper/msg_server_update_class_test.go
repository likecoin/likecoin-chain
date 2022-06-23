package keeper_test

import (
	"testing"
	"time"

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

func TestUpdateClassISCNNormal(t *testing.T) {
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
	maxSupply := uint64(5)

	// Mock keeper calls
	oldClassData := types.ClassData{
		Metadata: types.JsonInput(`{"aaaa": "bbbb"}`),
		Parent: types.ClassParent{
			Type:         types.ClassParentType_ISCN,
			IscnIdPrefix: iscnId.Prefix.String(),
		},
		Config: types.ClassConfig{
			Burnable:  false,
			MaxSupply: uint64(500),
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
		Creator: ownerAddress,
		ClassId: classId,
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
	require.Equal(t, iscnId.Prefix.String(), classData.Parent.IscnIdPrefix)
	require.Equal(t, iscnLatestVersion, classData.Parent.IscnVersionAtMint)
	require.Equal(t, types.ClassConfig{
		Burnable:  burnable,
		MaxSupply: maxSupply,
	}, classData.Config)

	// Check mock was called as expected
	ctrl.Finish()
}

func TestUpdateClassAccountNormal(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	accountKeeper := testutil.NewMockAccountKeeper(ctrl)
	bankKeeper := testutil.NewMockBankKeeper(ctrl)
	nftKeeper := testutil.NewMockNftKeeper(ctrl)
	msgServer, goCtx, keeper := setupMsgServer(t, keeper.LikenftDependedKeepers{
		AccountKeeper: accountKeeper,
		BankKeeper:    bankKeeper,
		NftKeeper:     nftKeeper,
	})
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Test Input
	ownerAddressBytes := []byte{0, 1, 0, 1, 0, 1, 0, 1}
	ownerAddress, _ := sdk.Bech32ifyAddressBytes("cosmos", ownerAddressBytes)
	ownerAccAddress, _ := sdk.AccAddressFromBech32(ownerAddress)

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
	maxSupply := uint64(5)

	// Mock keeper calls
	oldClassData := types.ClassData{
		Metadata: types.JsonInput(`{"aaaa": "bbbb"}`),
		Parent: types.ClassParent{
			Type:    types.ClassParentType_ACCOUNT,
			Account: ownerAddress,
		},
		Config: types.ClassConfig{
			Burnable:  false,
			MaxSupply: uint64(500),
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

	keeper.SetClassesByAccount(ctx, types.ClassesByAccount{
		Account:  ownerAddress,
		ClassIds: []string{classId},
	})

	nftKeeper.
		EXPECT().
		UpdateClass(gomock.Any(), gomock.Any()).
		Return(nil)

	// Run
	res, err := msgServer.UpdateClass(goCtx, &types.MsgUpdateClass{
		Creator: ownerAddress,
		ClassId: classId,
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
	require.Equal(t, ownerAccAddress.String(), classData.Parent.Account)
	require.Equal(t, types.ClassConfig{
		Burnable:  burnable,
		MaxSupply: maxSupply,
	}, classData.Config)

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
	maxSupply := uint64(5)

	// Mock keeper calls
	nftKeeper.
		EXPECT().
		GetClass(gomock.Any(), classId).
		Return(nft.Class{}, false)

	// Run
	res, err := msgServer.UpdateClass(goCtx, &types.MsgUpdateClass{
		Creator: ownerAddress,
		ClassId: classId,
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
	maxSupply := uint64(5)

	// Mock keeper calls
	oldClassData := types.ClassData{
		Metadata: types.JsonInput(`{"aaaa": "bbbb"}`),
		Parent: types.ClassParent{
			Type:         types.ClassParentType_ISCN,
			IscnIdPrefix: iscnId.Prefix.String(),
		},
		Config: types.ClassConfig{
			Burnable:  false,
			MaxSupply: uint64(500),
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
		Creator: ownerAddress,
		ClassId: classId,
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
	maxSupply := uint64(5)

	// Mock keeper calls
	oldClassData := types.ClassData{
		Metadata: types.JsonInput(`{"aaaa": "bbbb"}`),
		Parent: types.ClassParent{
			Type:         types.ClassParentType_ISCN,
			IscnIdPrefix: iscnId.Prefix.String(),
		},
		Config: types.ClassConfig{
			Burnable:  false,
			MaxSupply: uint64(500),
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
		Creator: ownerAddress,
		ClassId: classId,
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
	maxSupply := uint64(5)

	// Mock keeper calls
	oldClassData := types.ClassData{
		Metadata: types.JsonInput(`{"aaaa": "bbbb"}`),
		Parent: types.ClassParent{
			Type:         types.ClassParentType_ISCN,
			IscnIdPrefix: iscnId.Prefix.String(),
		},
		Config: types.ClassConfig{
			Burnable:  false,
			MaxSupply: uint64(500),
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
		Creator: ownerAddress,
		ClassId: classId,
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
	require.Contains(t, err.Error(), types.ErrNftClassNotRelatedToAnyIscn.Error())
	require.Nil(t, res)

	// Check mock was called as expected
	ctrl.Finish()
}

func TestUpdateClassNotRelatedToAccount(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	accountKeeper := testutil.NewMockAccountKeeper(ctrl)
	bankKeeper := testutil.NewMockBankKeeper(ctrl)
	nftKeeper := testutil.NewMockNftKeeper(ctrl)
	msgServer, goCtx, keeper := setupMsgServer(t, keeper.LikenftDependedKeepers{
		AccountKeeper: accountKeeper,
		BankKeeper:    bankKeeper,
		NftKeeper:     nftKeeper,
	})
	ctx := sdk.UnwrapSDKContext(goCtx)

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
	maxSupply := uint64(5)

	// Mock keeper calls
	oldClassData := types.ClassData{
		Metadata: types.JsonInput(`{"aaaa": "bbbb"}`),
		Parent: types.ClassParent{
			Type:    types.ClassParentType_ACCOUNT,
			Account: ownerAddress,
		},
		Config: types.ClassConfig{
			Burnable:  false,
			MaxSupply: uint64(500),
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

	keeper.SetClassesByAccount(ctx, types.ClassesByAccount{
		Account:  ownerAddress,
		ClassIds: []string{"likenft199887766"},
	})

	// Run
	res, err := msgServer.UpdateClass(goCtx, &types.MsgUpdateClass{
		Creator: ownerAddress,
		ClassId: classId,
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
	require.Contains(t, err.Error(), types.ErrNftClassNotRelatedToAnyAccount.Error())
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
	maxSupply := uint64(5)

	// Mock keeper calls
	oldClassData := types.ClassData{
		Metadata: types.JsonInput(`{"aaaa": "bbbb"}`),
		Parent: types.ClassParent{
			Type:         types.ClassParentType_ISCN,
			IscnIdPrefix: iscnId.Prefix.String(),
		},
		Config: types.ClassConfig{
			Burnable:  false,
			MaxSupply: uint64(500),
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
		Creator: ownerAddress,
		ClassId: classId,
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
	maxSupply := uint64(5)

	// Mock keeper calls
	oldClassData := types.ClassData{
		Metadata: types.JsonInput(`{"aaaa": "bbbb"}`),
		Parent: types.ClassParent{
			Type:         types.ClassParentType_ISCN,
			IscnIdPrefix: iscnId.Prefix.String(),
		},
		Config: types.ClassConfig{
			Burnable:  false,
			MaxSupply: uint64(500),
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
		Creator: "not an address",
		ClassId: classId,
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

func TestUpdateClassUserNotISCNOwner(t *testing.T) {
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
	maxSupply := uint64(5)

	// Mock keeper calls
	oldClassData := types.ClassData{
		Metadata: types.JsonInput(`{"aaaa": "bbbb"}`),
		Parent: types.ClassParent{
			Type:         types.ClassParentType_ISCN,
			IscnIdPrefix: iscnId.Prefix.String(),
		},
		Config: types.ClassConfig{
			Burnable:  false,
			MaxSupply: uint64(500),
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
		Creator: ownerAddress,
		ClassId: classId,
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
	require.Contains(t, err.Error(), sdkerrors.ErrUnauthorized.Error())
	require.Nil(t, res)

	// Check mock was called as expected
	ctrl.Finish()
}

func TestUpdateClassUserNotAccountOwner(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	accountKeeper := testutil.NewMockAccountKeeper(ctrl)
	bankKeeper := testutil.NewMockBankKeeper(ctrl)
	nftKeeper := testutil.NewMockNftKeeper(ctrl)
	msgServer, goCtx, keeper := setupMsgServer(t, keeper.LikenftDependedKeepers{
		AccountKeeper: accountKeeper,
		BankKeeper:    bankKeeper,
		NftKeeper:     nftKeeper,
	})
	ctx := sdk.UnwrapSDKContext(goCtx)

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
	maxSupply := uint64(5)

	// Mock keeper calls
	oldClassData := types.ClassData{
		Metadata: types.JsonInput(`{"aaaa": "bbbb"}`),
		Parent: types.ClassParent{
			Type:    types.ClassParentType_ACCOUNT,
			Account: ownerAddress,
		},
		Config: types.ClassConfig{
			Burnable:  false,
			MaxSupply: uint64(500),
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

	keeper.SetClassesByAccount(ctx, types.ClassesByAccount{
		Account:  ownerAddress,
		ClassIds: []string{classId},
	})

	notOwnerAddressBytes := []byte{1, 1, 1, 1, 1, 1, 1, 1}
	notOwnerAddress, _ := sdk.Bech32ifyAddressBytes("cosmos", notOwnerAddressBytes)

	// Run
	res, err := msgServer.UpdateClass(goCtx, &types.MsgUpdateClass{
		Creator: notOwnerAddress,
		ClassId: classId,
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
	require.Contains(t, err.Error(), sdkerrors.ErrUnauthorized.Error())
	require.Nil(t, res)

	// Check mock was called as expected
	ctrl.Finish()
}

func TestUpdateClassEnableBlindBox(t *testing.T) {
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
	ownerAddressBytes := []byte{0, 1, 0, 1, 0, 1, 0, 1, 1, 1}
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
	maxSupply := uint64(5)

	// Mock keeper calls
	oldClassData := types.ClassData{
		Metadata: types.JsonInput(`{"aaaa": "bbbb"}`),
		Parent: types.ClassParent{
			Type:         types.ClassParentType_ISCN,
			IscnIdPrefix: iscnId.Prefix.String(),
		},
		Config: types.ClassConfig{
			Burnable:  false,
			MaxSupply: uint64(500),
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

	nftKeeper.
		EXPECT().
		GetTotalSupply(gomock.Any(), classId).
		Return(uint64(0))

	keeper.SetClassesByISCN(ctx, types.ClassesByISCN{
		IscnIdPrefix: iscnId.Prefix.String(),
		ClassIds:     []string{classId},
	})

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
		UpdateClass(gomock.Any(), gomock.Any()).
		Return(nil)

	// Ensure queue is empty
	revealQueue := keeper.GetClassRevealQueue(ctx)
	require.Equal(t, 0, len(revealQueue))

	// Run
	res, err := msgServer.UpdateClass(goCtx, &types.MsgUpdateClass{
		Creator: ownerAddress,
		ClassId: classId,
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

	// Check class is now enqueued
	revealQueue = keeper.GetClassRevealQueue(ctx)
	require.Contains(t, revealQueue, types.ClassRevealQueueEntry{
		ClassId:    classId,
		RevealTime: revealTime,
	})

	// Check mock was called as expected
	ctrl.Finish()
}

func TestUpdateClassDisableBlindBox(t *testing.T) {
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
	ownerAddressBytes := []byte{0, 1, 0, 1, 0, 1, 0, 1, 1, 1}
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
	maxSupply := uint64(5)
	// Mock keeper calls
	revealTime := *testutil.MustParseTime(time.RFC3339, "2022-04-28T00:00:00Z")
	oldClassData := types.ClassData{
		Metadata: types.JsonInput(`{"aaaa": "bbbb"}`),
		Parent: types.ClassParent{
			Type:         types.ClassParentType_ISCN,
			IscnIdPrefix: iscnId.Prefix.String(),
		},
		Config: types.ClassConfig{
			Burnable:  false,
			MaxSupply: uint64(500),
			BlindBoxConfig: &types.BlindBoxConfig{
				MintPeriods: []types.MintPeriod{
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
				},
				RevealTime: revealTime,
			},
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
		UpdateClass(gomock.Any(), gomock.Any()).
		Return(nil)

	// Assume entry is inserted
	keeper.SetClassRevealQueueEntry(ctx, types.ClassRevealQueueEntry{
		ClassId:    classId,
		RevealTime: revealTime,
	})
	// Run
	res, err := msgServer.UpdateClass(goCtx, &types.MsgUpdateClass{
		Creator: ownerAddress,
		ClassId: classId,
		Input: types.ClassInput{
			Name:        name,
			Symbol:      symbol,
			Description: description,
			Uri:         uri,
			UriHash:     uriHash,
			Metadata:    metadata,
			Config: types.ClassConfig{
				Burnable:       burnable,
				MaxSupply:      maxSupply,
				BlindBoxConfig: nil,
			},
		},
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
	require.Equal(t, iscnId.Prefix.String(), classData.Parent.IscnIdPrefix)
	require.Equal(t, iscnLatestVersion, classData.Parent.IscnVersionAtMint)
	require.Equal(t, burnable, classData.Config.Burnable)
	require.Equal(t, maxSupply, classData.Config.MaxSupply)
	require.Nil(t, classData.Config.BlindBoxConfig)

	// Check class is now enqueued
	revealQueue := keeper.GetClassRevealQueue(ctx)
	require.Equal(t, 0, len(revealQueue))

	// Check mock was called as expected
	ctrl.Finish()
}

func TestUpdateClassUpdateMintPeriod(t *testing.T) {
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
	ownerAddressBytes := []byte{0, 1, 0, 1, 0, 1, 0, 1, 1, 1}
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
	maxSupply := uint64(5)

	// Mock keeper calls
	oldMintPeriods := []types.MintPeriod{
		{
			StartTime:        *testutil.MustParseTime(time.RFC3339, "2022-04-19T00:00:00Z"),
			AllowedAddresses: []string{ownerAddress},
			MintPrice:        uint64(2048),
		},
	}
	oldRevealTime := *testutil.MustParseTime(time.RFC3339, "2022-05-01T00:00:00Z")
	oldClassData := types.ClassData{
		Metadata: types.JsonInput(`{"aaaa": "bbbb"}`),
		Parent: types.ClassParent{
			Type:         types.ClassParentType_ISCN,
			IscnIdPrefix: iscnId.Prefix.String(),
		},
		Config: types.ClassConfig{
			Burnable:  false,
			MaxSupply: uint64(500),
			BlindBoxConfig: &types.BlindBoxConfig{
				MintPeriods: oldMintPeriods,
				RevealTime:  oldRevealTime,
			},
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

	// Assume entry is inserted
	keeper.SetClassRevealQueueEntry(ctx, types.ClassRevealQueueEntry{
		ClassId:    classId,
		RevealTime: oldRevealTime,
	})

	updatedMintPeriods := []types.MintPeriod{
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
	updatedRevealTime := *testutil.MustParseTime(time.RFC3339, "2022-04-28T00:00:00Z")

	nftKeeper.
		EXPECT().
		GetTotalSupply(gomock.Any(), classId).
		Return(uint64(0))

	keeper.SetClassesByISCN(ctx, types.ClassesByISCN{
		IscnIdPrefix: iscnId.Prefix.String(),
		ClassIds:     []string{classId},
	})

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
		UpdateClass(gomock.Any(), gomock.Any()).
		Return(nil)

	// Run
	res, err := msgServer.UpdateClass(goCtx, &types.MsgUpdateClass{
		Creator: ownerAddress,
		ClassId: classId,
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
					MintPeriods: updatedMintPeriods,
					RevealTime:  updatedRevealTime,
				},
			},
		},
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
	require.Equal(t, iscnId.Prefix.String(), classData.Parent.IscnIdPrefix)
	require.Equal(t, iscnLatestVersion, classData.Parent.IscnVersionAtMint)
	require.Equal(t, burnable, classData.Config.Burnable)
	require.Equal(t, maxSupply, classData.Config.MaxSupply)
	require.Equal(t, updatedRevealTime, classData.Config.BlindBoxConfig.RevealTime)

	require.Equal(t, len(updatedMintPeriods), len(classData.Config.BlindBoxConfig.MintPeriods))
	for i, mintPeriod := range classData.Config.BlindBoxConfig.MintPeriods {
		require.Equal(t, mintPeriod.StartTime, updatedMintPeriods[i].StartTime)
		require.ElementsMatch(t, mintPeriod.AllowedAddresses, updatedMintPeriods[i].AllowedAddresses)
		require.Equal(t, mintPeriod.MintPrice, updatedMintPeriods[i].MintPrice)
	}

	// Check class is now enqueued
	revealQueue := keeper.GetClassRevealQueue(ctx)
	require.Contains(t, revealQueue, types.ClassRevealQueueEntry{
		ClassId:    classId,
		RevealTime: updatedRevealTime,
	})

	// Check mock was called as expected
	ctrl.Finish()
}

func TestUpdateClassInvalidMintPeriod(t *testing.T) {
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
	ownerAddressBytes := []byte{0, 1, 0, 1, 0, 1, 0, 1, 1, 1}
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
	maxSupply := uint64(5)

	mintPeriods := []types.MintPeriod{
		{
			StartTime:        *testutil.MustParseTime(time.RFC3339, "2022-04-21T00:00:00Z"),
			AllowedAddresses: make([]string, 0),
			MintPrice:        0,
		},
	}
	revealTime := *testutil.MustParseTime(time.RFC3339, "2322-04-20T00:00:00Z")

	// Mock keeper calls
	oldClassData := types.ClassData{
		Metadata: types.JsonInput(`{"aaaa": "bbbb"}`),
		Parent: types.ClassParent{
			Type:         types.ClassParentType_ISCN,
			IscnIdPrefix: iscnId.Prefix.String(),
		},
		Config: types.ClassConfig{
			Burnable:  false,
			MaxSupply: uint64(500),
			BlindBoxConfig: &types.BlindBoxConfig{
				MintPeriods: mintPeriods,
				RevealTime:  revealTime,
			},
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

	newMintPeriods := []types.MintPeriod{
		{
			StartTime:        *testutil.MustParseTime(time.RFC3339, "2922-04-21T00:00:00Z"),
			AllowedAddresses: make([]string, 0),
			MintPrice:        0,
		},
	}

	// Run
	res, err := msgServer.UpdateClass(goCtx, &types.MsgUpdateClass{
		Creator: ownerAddress,
		ClassId: classId,
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
					MintPeriods: newMintPeriods,
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

func TestUpdateClassMintPeriodInvalidAllowListAddress(t *testing.T) {
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
	ownerAddressBytes := []byte{0, 1, 0, 1, 0, 1, 0, 1, 1, 1}
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
	maxSupply := uint64(5)

	mintPeriods := []types.MintPeriod{
		{
			StartTime:        *testutil.MustParseTime(time.RFC3339, "2022-04-19T00:00:00Z"),
			AllowedAddresses: []string{"invalid address"},
			MintPrice:        0,
		},
	}
	revealTime := *testutil.MustParseTime(time.RFC3339, "2322-04-20T00:00:00Z")

	// Mock keeper calls
	oldClassData := types.ClassData{
		Metadata: types.JsonInput(`{"aaaa": "bbbb"}`),
		Parent: types.ClassParent{
			Type:         types.ClassParentType_ISCN,
			IscnIdPrefix: iscnId.Prefix.String(),
		},
		Config: types.ClassConfig{
			Burnable:  false,
			MaxSupply: uint64(500),
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
		Creator: ownerAddress,
		ClassId: classId,
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

func TestUpdateClassNoMintPeriod(t *testing.T) {
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
	ownerAddressBytes := []byte{0, 1, 0, 1, 0, 1, 0, 1, 1, 1}
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
	maxSupply := uint64(5)

	mintPeriods := []types.MintPeriod{
		{
			StartTime:        *testutil.MustParseTime(time.RFC3339, "2022-04-21T00:00:00Z"),
			AllowedAddresses: make([]string, 0),
			MintPrice:        0,
		},
	}
	revealTime := *testutil.MustParseTime(time.RFC3339, "2322-04-20T00:00:00Z")

	// Mock keeper calls
	oldClassData := types.ClassData{
		Metadata: types.JsonInput(`{"aaaa": "bbbb"}`),
		Parent: types.ClassParent{
			Type:         types.ClassParentType_ISCN,
			IscnIdPrefix: iscnId.Prefix.String(),
		},
		Config: types.ClassConfig{
			Burnable:  false,
			MaxSupply: uint64(500),
			BlindBoxConfig: &types.BlindBoxConfig{
				MintPeriods: mintPeriods,
				RevealTime:  revealTime,
			},
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
		Creator: ownerAddress,
		ClassId: classId,
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

func TestUpdateClassMaxSupplyNotLessThanMintableCount(t *testing.T) {
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
	ownerAddressBytes := []byte{0, 1, 0, 1, 0, 1, 0, 1, 1, 1}
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
	maxSupply := uint64(499)

	// Mock keeper calls
	oldClassData := types.ClassData{
		Metadata: types.JsonInput(`{"aaaa": "bbbb"}`),
		Parent: types.ClassParent{
			Type:         types.ClassParentType_ISCN,
			IscnIdPrefix: iscnId.Prefix.String(),
		},
		Config: types.ClassConfig{
			Burnable:  false,
			MaxSupply: uint64(500),
		},
		MintableCount: uint64(500),
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

	nftKeeper.
		EXPECT().
		GetTotalSupply(gomock.Any(), classId).
		Return(uint64(0))

	keeper.SetClassesByISCN(ctx, types.ClassesByISCN{
		IscnIdPrefix: iscnId.Prefix.String(),
		ClassIds:     []string{classId},
	})

	// Ensure queue is empty
	revealQueue := keeper.GetClassRevealQueue(ctx)
	require.Equal(t, 0, len(revealQueue))

	// Run
	res, err := msgServer.UpdateClass(goCtx, &types.MsgUpdateClass{
		Creator: ownerAddress,
		ClassId: classId,
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
	require.Contains(t, err.Error(), sdkerrors.ErrInvalidRequest.Error())
	require.Nil(t, res)

	// Check class is not enqueued
	revealQueue = keeper.GetClassRevealQueue(ctx)
	require.Equal(t, 0, len(revealQueue))

	// Check mock was called as expected
	ctrl.Finish()
}
