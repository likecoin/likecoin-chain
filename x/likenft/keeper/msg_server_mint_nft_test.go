package keeper_test

import (
	"fmt"
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
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

func TestMintOwnerNFTNormal(t *testing.T) {
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
	classId := "likenft1aabbccddeeff"
	tokenId := "token1"
	uri := "ipfs://a1b2c3"
	uriHash := "a1b2c3"
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

	// Mock keeper calls
	classIscnVersionAtMint := uint64(1)
	classData := types.ClassData{
		Metadata: types.JsonInput(`{"aaaa": "bbbb"}`),
		Parent: types.ClassParent{
			Type:              types.ClassParentType_ISCN,
			IscnIdPrefix:      iscnId.Prefix.String(),
			IscnVersionAtMint: classIscnVersionAtMint,
		},
		Config: types.ClassConfig{
			Burnable:  false,
			MaxSupply: uint64(5),
		},
	}
	classDataInAny, _ := cdctypes.NewAnyWithValue(&classData)
	nftKeeper.
		EXPECT().
		GetClass(gomock.Any(), gomock.Eq(classId)).
		Return(nft.Class{
			Id:          classId,
			Name:        "Class Name",
			Symbol:      "ABC",
			Description: "Testing Class 123",
			Uri:         "ipfs://abcdef",
			UriHash:     "abcdef",
			Data:        classDataInAny,
		}, true)

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

	// Test for subsequent nft mint at this case
	// No class update
	nftKeeper.
		EXPECT().
		GetTotalSupply(gomock.Any(), gomock.Eq(classId)).
		Return(uint64(1))

	wrappedOwnerAddress, _ := sdk.AccAddressFromBech32(ownerAddress)
	nftKeeper.
		EXPECT().
		Mint(gomock.Any(), gomock.Any(), wrappedOwnerAddress)

	// Run
	res, err := msgServer.MintNFT(goCtx, &types.MsgMintNFT{
		Creator: ownerAddress,
		ClassId: classId,
		Id:      tokenId,
		Input: types.NFTInput{
			Uri:      uri,
			UriHash:  uriHash,
			Metadata: metadata,
		},
	})

	// Check output
	require.NoError(t, err)
	require.Equal(t, classId, res.Nft.ClassId)
	require.Equal(t, uri, res.Nft.Uri)
	require.Equal(t, uriHash, res.Nft.UriHash)

	var nftData types.NFTData
	err = nftData.Unmarshal(res.Nft.Data.Value)
	require.NoErrorf(t, err, "Error unmarshal class data")
	require.Equal(t, metadata, nftData.Metadata)
	require.Equal(t, iscnId.Prefix.String(), nftData.ClassParent.IscnIdPrefix)
	require.Equal(t, classIscnVersionAtMint, nftData.ClassParent.IscnVersionAtMint)

	// Check mock was called as expected
	ctrl.Finish()
}

type IFirstTokenUpdateClassMatcher struct {
	iscnVersionAtMint uint64
}

func firstTokenUpdateClassMatcher(iscnVersionAtMint uint64) gomock.Matcher {
	return IFirstTokenUpdateClassMatcher{
		iscnVersionAtMint,
	}
}

func (m IFirstTokenUpdateClassMatcher) Matches(x interface{}) bool {
	class := x.(nft.Class)
	var classData types.ClassData
	if err := classData.Unmarshal(class.Data.Value); err != nil {
		return false
	}
	return classData.Parent.IscnVersionAtMint == m.iscnVersionAtMint
}

func (m IFirstTokenUpdateClassMatcher) String() string {
	return fmt.Sprintf("data item iscnVersionAtMint is equal to %d", m.iscnVersionAtMint)
}

func TestMintNFTFirstToken(t *testing.T) {
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
	classId := "likenft1aabbccddeeff"
	tokenId := "token1"
	uri := "ipfs://a1b2c3"
	uriHash := "a1b2c3"
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

	// Mock keeper calls
	classIscnVersionAtMint := uint64(1)
	classData := types.ClassData{
		Metadata: types.JsonInput(`{"aaaa": "bbbb"}`),
		Parent: types.ClassParent{
			Type:              types.ClassParentType_ISCN,
			IscnIdPrefix:      iscnId.Prefix.String(),
			IscnVersionAtMint: classIscnVersionAtMint,
		},
		Config: types.ClassConfig{
			Burnable:  false,
			MaxSupply: uint64(5),
		},
	}
	classDataInAny, _ := cdctypes.NewAnyWithValue(&classData)
	nftKeeper.
		EXPECT().
		GetClass(gomock.Any(), gomock.Eq(classId)).
		Return(nft.Class{
			Id:          classId,
			Name:        "Class Name",
			Symbol:      "ABC",
			Description: "Testing Class 123",
			Uri:         "ipfs://abcdef",
			UriHash:     "abcdef",
			Data:        classDataInAny,
		}, true)

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

	// Test for first nft mint at this case
	// Should refresh iscn version in class data
	nftKeeper.
		EXPECT().
		GetTotalSupply(gomock.Any(), gomock.Eq(classId)).
		Return(uint64(0))

	updateClassMatcher := firstTokenUpdateClassMatcher(iscnLatestVersion)
	nftKeeper.
		EXPECT().
		UpdateClass(gomock.Any(), updateClassMatcher).
		Return(nil)

	wrappedOwnerAddress, _ := sdk.AccAddressFromBech32(ownerAddress)
	nftKeeper.
		EXPECT().
		Mint(gomock.Any(), gomock.Any(), wrappedOwnerAddress)

	// Run
	res, err := msgServer.MintNFT(goCtx, &types.MsgMintNFT{
		Creator: ownerAddress,
		ClassId: classId,
		Id:      tokenId,
		Input: types.NFTInput{
			Uri:      uri,
			UriHash:  uriHash,
			Metadata: metadata,
		},
	})

	// Check output
	require.NoError(t, err)
	require.Equal(t, classId, res.Nft.ClassId)
	require.Equal(t, uri, res.Nft.Uri)
	require.Equal(t, uriHash, res.Nft.UriHash)

	var nftData types.NFTData
	err = nftData.Unmarshal(res.Nft.Data.Value)
	require.NoErrorf(t, err, "Error unmarshal class data")
	require.Equal(t, metadata, nftData.Metadata)
	require.Equal(t, iscnId.Prefix.String(), nftData.ClassParent.IscnIdPrefix)
	require.Equal(t, iscnLatestVersion, nftData.ClassParent.IscnVersionAtMint)

	// Check mock was called as expected
	ctrl.Finish()
}

func TestMintNFTOwnerInvalidTokenID(t *testing.T) {
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
	classId := "likenft1aabbccddeeff"
	uri := "ipfs://a1b2c3"
	uriHash := "a1b2c3"
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

	// Mock keeper calls
	classIscnVersionAtMint := uint64(1)
	classData := types.ClassData{
		Metadata: types.JsonInput(`{"aaaa": "bbbb"}`),
		Parent: types.ClassParent{
			Type:              types.ClassParentType_ISCN,
			IscnIdPrefix:      iscnId.Prefix.String(),
			IscnVersionAtMint: classIscnVersionAtMint,
		},
		Config: types.ClassConfig{
			Burnable:       false,
			MaxSupply:      uint64(5),
			EnableBlindBox: false,
		},
	}
	classDataInAny, _ := cdctypes.NewAnyWithValue(&classData)
	nftKeeper.
		EXPECT().
		GetClass(gomock.Any(), gomock.Eq(classId)).
		Return(nft.Class{
			Id:          classId,
			Name:        "Class Name",
			Symbol:      "ABC",
			Description: "Testing Class 123",
			Uri:         "ipfs://abcdef",
			UriHash:     "abcdef",
			Data:        classDataInAny,
		}, true)

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

	// Test for first nft mint at this case
	// Should refresh iscn version in class data
	nftKeeper.
		EXPECT().
		GetTotalSupply(gomock.Any(), gomock.Eq(classId)).
		Return(uint64(0))

	updateClassMatcher := firstTokenUpdateClassMatcher(iscnLatestVersion)
	nftKeeper.
		EXPECT().
		UpdateClass(gomock.Any(), updateClassMatcher).
		Return(nil)

	// Run
	res, err := msgServer.MintNFT(goCtx, &types.MsgMintNFT{
		Creator: ownerAddress,
		ClassId: classId,
		Id:      "123456", // x/nft requires token id to start with letter
		Input: types.NFTInput{
			Uri:      uri,
			UriHash:  uriHash,
			Metadata: metadata,
		},
	})

	// Check output
	require.Error(t, err)
	require.Contains(t, err.Error(), types.ErrInvalidTokenId.Error())
	require.Nil(t, res)

	// Check mock was called as expected
	ctrl.Finish()
}

func TestMintNFTClassNotFound(t *testing.T) {
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
	tokenId := "token1"
	uri := "ipfs://a1b2c3"
	uriHash := "a1b2c3"
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

	// Mock keeper calls
	nftKeeper.
		EXPECT().
		GetClass(gomock.Any(), gomock.Eq(classId)).
		Return(nft.Class{}, false)

	// Run
	res, err := msgServer.MintNFT(goCtx, &types.MsgMintNFT{
		Creator: ownerAddress,
		ClassId: classId,
		Id:      tokenId,
		Input: types.NFTInput{
			Uri:      uri,
			UriHash:  uriHash,
			Metadata: metadata,
		},
	})

	// Check output
	require.Error(t, err)
	require.Contains(t, err.Error(), types.ErrNftClassNotFound.Error())
	require.Nil(t, res)

	// Check mock was called as expected
	ctrl.Finish()
}

func TestMintNFTMissingIscnRelation(t *testing.T) {
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
	classId := "likenft1aabbccddeeff"
	tokenId := "token1"
	uri := "ipfs://a1b2c3"
	uriHash := "a1b2c3"
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

	// Mock keeper calls
	classData := types.ClassData{
		Metadata: types.JsonInput(`{"aaaa": "bbbb"}`),
		Parent: types.ClassParent{
			Type:              types.ClassParentType_ISCN,
			IscnIdPrefix:      iscnId.Prefix.String(),
			IscnVersionAtMint: uint64(1),
		},
		Config: types.ClassConfig{
			Burnable:  false,
			MaxSupply: uint64(0),
		},
	}
	classDataInAny, _ := cdctypes.NewAnyWithValue(&classData)
	nftKeeper.
		EXPECT().
		GetClass(gomock.Any(), gomock.Eq(classId)).
		Return(nft.Class{
			Id:          classId,
			Name:        "Class Name",
			Symbol:      "ABC",
			Description: "Testing Class 123",
			Uri:         "ipfs://abcdef",
			UriHash:     "abcdef",
			Data:        classDataInAny,
		}, true)

	// Run
	res, err := msgServer.MintNFT(goCtx, &types.MsgMintNFT{
		Creator: ownerAddress,
		ClassId: classId,
		Id:      tokenId,
		Input: types.NFTInput{
			Uri:      uri,
			UriHash:  uriHash,
			Metadata: metadata,
		},
	})

	// Check output
	require.Error(t, err)
	require.Contains(t, err.Error(), types.ErrNftClassNotRelatedToAnyIscn.Error())
	require.Nil(t, res)

	// Check mock was called as expected
	ctrl.Finish()
}

func TestMintNFTNotRelatedToIscn(t *testing.T) {
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
	classId := "likenft1aabbccddeeff"
	tokenId := "token1"
	uri := "ipfs://a1b2c3"
	uriHash := "a1b2c3"
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

	// Mock keeper calls
	classData := types.ClassData{
		Metadata: types.JsonInput(`{"aaaa": "bbbb"}`),
		Parent: types.ClassParent{
			Type:              types.ClassParentType_ISCN,
			IscnIdPrefix:      iscnId.Prefix.String(),
			IscnVersionAtMint: uint64(1),
		},
		Config: types.ClassConfig{
			Burnable:  false,
			MaxSupply: uint64(0),
		},
	}
	classDataInAny, _ := cdctypes.NewAnyWithValue(&classData)
	nftKeeper.
		EXPECT().
		GetClass(gomock.Any(), gomock.Eq(classId)).
		Return(nft.Class{
			Id:          classId,
			Name:        "Class Name",
			Symbol:      "ABC",
			Description: "Testing Class 123",
			Uri:         "ipfs://abcdef",
			UriHash:     "abcdef",
			Data:        classDataInAny,
		}, true)

	keeper.SetClassesByISCN(ctx, types.ClassesByISCN{
		IscnIdPrefix: iscnId.Prefix.String(),
		ClassIds:     []string{"likenft1112233445566"},
	})

	// Run
	res, err := msgServer.MintNFT(goCtx, &types.MsgMintNFT{
		Creator: ownerAddress,
		ClassId: classId,
		Id:      tokenId,
		Input: types.NFTInput{
			Uri:      uri,
			UriHash:  uriHash,
			Metadata: metadata,
		},
	})

	// Check output
	require.Error(t, err)
	require.Contains(t, err.Error(), types.ErrNftClassNotRelatedToAnyIscn.Error())
	require.Nil(t, res)

	// Check mock was called as expected
	ctrl.Finish()
}

func TestMintNFTIscnNotFound(t *testing.T) {
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
	classId := "likenft1aabbccddeeff"
	tokenId := "token1"
	uri := "ipfs://a1b2c3"
	uriHash := "a1b2c3"
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

	// Mock keeper calls
	classData := types.ClassData{
		Metadata: types.JsonInput(`{"aaaa": "bbbb"}`),
		Parent: types.ClassParent{
			Type:              types.ClassParentType_ISCN,
			IscnIdPrefix:      iscnId.Prefix.String(),
			IscnVersionAtMint: uint64(1),
		},
		Config: types.ClassConfig{
			Burnable:  false,
			MaxSupply: uint64(0),
		},
	}
	classDataInAny, _ := cdctypes.NewAnyWithValue(&classData)
	nftKeeper.
		EXPECT().
		GetClass(gomock.Any(), gomock.Eq(classId)).
		Return(nft.Class{
			Id:          classId,
			Name:        "Class Name",
			Symbol:      "ABC",
			Description: "Testing Class 123",
			Uri:         "ipfs://abcdef",
			UriHash:     "abcdef",
			Data:        classDataInAny,
		}, true)

	keeper.SetClassesByISCN(ctx, types.ClassesByISCN{
		IscnIdPrefix: iscnId.Prefix.String(),
		ClassIds:     []string{classId},
	})

	iscnKeeper.
		EXPECT().
		GetContentIdRecord(gomock.Any(), gomock.Eq(iscnId.Prefix)).
		Return(nil)

	// Run
	res, err := msgServer.MintNFT(goCtx, &types.MsgMintNFT{
		Creator: ownerAddress,
		ClassId: classId,
		Id:      tokenId,
		Input: types.NFTInput{
			Uri:      uri,
			UriHash:  uriHash,
			Metadata: metadata,
		},
	})

	// Check output
	require.Error(t, err)
	require.Contains(t, err.Error(), types.ErrIscnRecordNotFound.Error())
	require.Nil(t, res)

	// Check mock was called as expected
	ctrl.Finish()
}

func TestMintNFTInvalidUserAddress(t *testing.T) {
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
	iscnId := iscntypes.NewIscnId("likecoin-chain", "abcdef", 1)
	classId := "likenft1aabbccddeeff"
	tokenId := "token1"
	uri := "ipfs://a1b2c3"
	uriHash := "a1b2c3"
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

	// Mock keeper calls
	classData := types.ClassData{
		Metadata: types.JsonInput(`{"aaaa": "bbbb"}`),
		Parent: types.ClassParent{
			Type:              types.ClassParentType_ISCN,
			IscnIdPrefix:      iscnId.Prefix.String(),
			IscnVersionAtMint: uint64(1),
		},
		Config: types.ClassConfig{
			Burnable:  false,
			MaxSupply: uint64(0),
		},
	}
	classDataInAny, _ := cdctypes.NewAnyWithValue(&classData)
	nftKeeper.
		EXPECT().
		GetClass(gomock.Any(), gomock.Eq(classId)).
		Return(nft.Class{
			Id:          classId,
			Name:        "Class Name",
			Symbol:      "ABC",
			Description: "Testing Class 123",
			Uri:         "ipfs://abcdef",
			UriHash:     "abcdef",
			Data:        classDataInAny,
		}, true)

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
	res, err := msgServer.MintNFT(goCtx, &types.MsgMintNFT{
		Creator: "not a valid address",
		ClassId: classId,
		Id:      tokenId,
		Input: types.NFTInput{
			Uri:      uri,
			UriHash:  uriHash,
			Metadata: metadata,
		},
	})

	// Check output
	require.Error(t, err)
	require.Contains(t, err.Error(), sdkerrors.ErrInvalidAddress.Error())
	require.Nil(t, res)

	// Check mock was called as expected
	ctrl.Finish()
}

func TestMintBlindBoxNFTNormal(t *testing.T) {
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
	ctx = ctx.WithBlockHeader(tmproto.Header{Time: *testutil.MustParseTime(time.RFC3339, "2022-04-20T00:00:00Z")})
	updatedGoCtx := sdk.WrapSDKContext(ctx)

	// Test Input
	ownerAddressBytes := []byte{0, 1, 0, 1, 0, 1, 0, 1}
	ownerAddress, _ := sdk.Bech32ifyAddressBytes("cosmos", ownerAddressBytes)
	minterAddressBytes := []byte{0, 1, 0, 1, 0, 1, 1, 0, 0, 1, 1}
	minterAddress, _ := sdk.Bech32ifyAddressBytes("cosmos", minterAddressBytes)
	iscnId := iscntypes.NewIscnId("likecoin-chain", "abcdef", 1)
	classId := "likenft1aabbccddeeff"
	tokenId := "token1"
	uri := "ipfs://a1b2c3"
	uriHash := "a1b2c3"
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
	mintPrice := uint64(5000000000)
	revealTime := testutil.MustParseTime(time.RFC3339, "2322-04-20T00:00:00Z")

	// Mock keeper calls
	classIscnVersionAtMint := uint64(1)
	classData := types.ClassData{
		Metadata: types.JsonInput(`{"aaaa": "bbbb"}`),
		Parent: types.ClassParent{
			Type:              types.ClassParentType_ISCN,
			IscnIdPrefix:      iscnId.Prefix.String(),
			IscnVersionAtMint: classIscnVersionAtMint,
		},
		Config: types.ClassConfig{
			Burnable:       false,
			MaxSupply:      uint64(5000000000),
			EnableBlindBox: true,
			ClaimPeriods: []*types.ClaimPeriod{
				{
					StartTime:        testutil.MustParseTime(time.RFC3339, "2020-01-01T00:00:00Z"),
					AllowedAddresses: nil,
					MintPrice:        mintPrice,
				},
			},
			RevealTime: revealTime,
		},
	}
	classDataInAny, _ := cdctypes.NewAnyWithValue(&classData)
	nftKeeper.
		EXPECT().
		GetClass(gomock.Any(), gomock.Eq(classId)).
		Return(nft.Class{
			Id:          classId,
			Name:        "Class Name",
			Symbol:      "ABC",
			Description: "Testing Class 123",
			Uri:         "ipfs://abcdef",
			UriHash:     "abcdef",
			Data:        classDataInAny,
		}, true)

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

	// Test for subsequent nft mint at this case
	// No class update
	nftKeeper.
		EXPECT().
		GetTotalSupply(gomock.Any(), gomock.Eq(classId)).
		Return(uint64(1))

	wrappedMinterAddress, _ := sdk.AccAddressFromBech32(minterAddress)
	nftKeeper.
		EXPECT().
		Mint(gomock.Any(), gomock.Any(), wrappedMinterAddress)

	bankKeeper.EXPECT().GetBalance(
		gomock.Any(),
		gomock.Eq(wrappedMinterAddress),
		gomock.Any(),
	).
		Return(sdk.Coin{
			Denom:  "nanolike",
			Amount: sdk.NewInt(100000000000),
		})

	wrappedOwnerAddress, _ := sdk.AccAddressFromBech32(ownerAddress)
	bankKeeper.EXPECT().SendCoins(
		gomock.Any(),
		gomock.Eq(wrappedMinterAddress),
		gomock.Eq(wrappedOwnerAddress),
		gomock.Eq(sdk.Coins{sdk.Coin{
			Denom:  "nanolike",
			Amount: sdk.NewInt(int64(mintPrice)),
		}}),
	)

	// Run
	res, err := msgServer.MintNFT(updatedGoCtx, &types.MsgMintNFT{
		Creator: minterAddress,
		ClassId: classId,
		Id:      tokenId,
		Input: types.NFTInput{
			Uri:      uri,
			UriHash:  uriHash,
			Metadata: metadata,
		},
	})

	// Check output
	require.NoError(t, err)
	require.Equal(t, classId, res.Nft.ClassId)
	require.Equal(t, uri, res.Nft.Uri)
	require.Equal(t, uriHash, res.Nft.UriHash)

	var nftData types.NFTData
	err = nftData.Unmarshal(res.Nft.Data.Value)
	require.NoErrorf(t, err, "Error unmarshal class data")
	require.Equal(t, types.JsonInput(nil), nftData.Metadata) // TODO: update metadata to support templates
	require.Equal(t, iscnId.Prefix.String(), nftData.ClassParent.IscnIdPrefix)
	require.Equal(t, classIscnVersionAtMint, nftData.ClassParent.IscnVersionAtMint)

	// Check mock was called as expected
	ctrl.Finish()
}

func TestMintBlindBoxNFTOwnerNoPay(t *testing.T) {
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
	ctx = ctx.WithBlockHeader(tmproto.Header{Time: *testutil.MustParseTime(time.RFC3339, "2022-04-20T00:00:00Z")})
	updatedGoCtx := sdk.WrapSDKContext(ctx)

	// Test Input
	ownerAddressBytes := []byte{0, 1, 0, 1, 0, 1, 0, 1}
	ownerAddress, _ := sdk.Bech32ifyAddressBytes("cosmos", ownerAddressBytes)
	iscnId := iscntypes.NewIscnId("likecoin-chain", "abcdef", 1)
	classId := "likenft1aabbccddeeff"
	tokenId := "token1"
	uri := "ipfs://a1b2c3"
	uriHash := "a1b2c3"
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
	mintPrice := uint64(5000000000)
	revealTime := testutil.MustParseTime(time.RFC3339, "2322-04-20T00:00:00Z")

	// Mock keeper calls
	classIscnVersionAtMint := uint64(1)
	classData := types.ClassData{
		Metadata: types.JsonInput(`{"aaaa": "bbbb"}`),
		Parent: types.ClassParent{
			Type:              types.ClassParentType_ISCN,
			IscnIdPrefix:      iscnId.Prefix.String(),
			IscnVersionAtMint: classIscnVersionAtMint,
		},
		Config: types.ClassConfig{
			Burnable:       false,
			MaxSupply:      uint64(5000000000),
			EnableBlindBox: true,
			ClaimPeriods: []*types.ClaimPeriod{
				{
					StartTime:        testutil.MustParseTime(time.RFC3339, "2020-01-01T00:00:00Z"),
					AllowedAddresses: nil,
					MintPrice:        mintPrice,
				},
			},
			RevealTime: revealTime,
		},
	}
	classDataInAny, _ := cdctypes.NewAnyWithValue(&classData)
	nftKeeper.
		EXPECT().
		GetClass(gomock.Any(), gomock.Eq(classId)).
		Return(nft.Class{
			Id:          classId,
			Name:        "Class Name",
			Symbol:      "ABC",
			Description: "Testing Class 123",
			Uri:         "ipfs://abcdef",
			UriHash:     "abcdef",
			Data:        classDataInAny,
		}, true)

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

	// Test for subsequent nft mint at this case
	// No class update
	nftKeeper.
		EXPECT().
		GetTotalSupply(gomock.Any(), gomock.Eq(classId)).
		Return(uint64(1))

	wrappedMinterAddress, _ := sdk.AccAddressFromBech32(ownerAddress)
	nftKeeper.
		EXPECT().
		Mint(gomock.Any(), gomock.Any(), wrappedMinterAddress)

	// Run
	res, err := msgServer.MintNFT(updatedGoCtx, &types.MsgMintNFT{
		Creator: ownerAddress,
		ClassId: classId,
		Id:      tokenId,
		Input: types.NFTInput{
			Uri:      uri,
			UriHash:  uriHash,
			Metadata: metadata,
		},
	})

	// Check output
	require.NoError(t, err)
	require.Equal(t, classId, res.Nft.ClassId)
	require.Equal(t, uri, res.Nft.Uri)
	require.Equal(t, uriHash, res.Nft.UriHash)

	var nftData types.NFTData
	err = nftData.Unmarshal(res.Nft.Data.Value)
	require.NoErrorf(t, err, "Error unmarshal class data")
	require.Equal(t, types.JsonInput(nil), nftData.Metadata) // TODO: update metadata to support templates
	require.Equal(t, iscnId.Prefix.String(), nftData.ClassParent.IscnIdPrefix)
	require.Equal(t, classIscnVersionAtMint, nftData.ClassParent.IscnVersionAtMint)

	// Check mock was called as expected
	ctrl.Finish()
}

func TestMintBlindBoxNFTOwnerIgnoreAllowList(t *testing.T) {
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
	ctx = ctx.WithBlockHeader(tmproto.Header{Time: *testutil.MustParseTime(time.RFC3339, "2022-04-20T00:00:00Z")})
	updatedGoCtx := sdk.WrapSDKContext(ctx)

	// Test Input
	ownerAddressBytes := []byte{0, 1, 0, 1, 0, 1, 0, 1}
	ownerAddress, _ := sdk.Bech32ifyAddressBytes("cosmos", ownerAddressBytes)
	iscnId := iscntypes.NewIscnId("likecoin-chain", "abcdef", 1)
	classId := "likenft1aabbccddeeff"
	tokenId := "token1"
	uri := "ipfs://a1b2c3"
	uriHash := "a1b2c3"
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
	mintPrice := uint64(5000000000)
	revealTime := testutil.MustParseTime(time.RFC3339, "2322-04-20T00:00:00Z")

	// Mock keeper calls
	classIscnVersionAtMint := uint64(1)
	classData := types.ClassData{
		Metadata: types.JsonInput(`{"aaaa": "bbbb"}`),
		Parent: types.ClassParent{
			Type:              types.ClassParentType_ISCN,
			IscnIdPrefix:      iscnId.Prefix.String(),
			IscnVersionAtMint: classIscnVersionAtMint,
		},
		Config: types.ClassConfig{
			Burnable:       false,
			MaxSupply:      uint64(5000000000),
			EnableBlindBox: true,
			ClaimPeriods: []*types.ClaimPeriod{
				{
					StartTime:        testutil.MustParseTime(time.RFC3339, "2020-01-01T00:00:00Z"),
					AllowedAddresses: []string{},
					MintPrice:        mintPrice,
				},
			},
			RevealTime: revealTime,
		},
	}
	classDataInAny, _ := cdctypes.NewAnyWithValue(&classData)
	nftKeeper.
		EXPECT().
		GetClass(gomock.Any(), gomock.Eq(classId)).
		Return(nft.Class{
			Id:          classId,
			Name:        "Class Name",
			Symbol:      "ABC",
			Description: "Testing Class 123",
			Uri:         "ipfs://abcdef",
			UriHash:     "abcdef",
			Data:        classDataInAny,
		}, true)

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

	// Test for subsequent nft mint at this case
	// No class update
	nftKeeper.
		EXPECT().
		GetTotalSupply(gomock.Any(), gomock.Eq(classId)).
		Return(uint64(1))

	wrappedMinterAddress, _ := sdk.AccAddressFromBech32(ownerAddress)
	nftKeeper.
		EXPECT().
		Mint(gomock.Any(), gomock.Any(), wrappedMinterAddress)

	// Run
	res, err := msgServer.MintNFT(updatedGoCtx, &types.MsgMintNFT{
		Creator: ownerAddress,
		ClassId: classId,
		Id:      tokenId,
		Input: types.NFTInput{
			Uri:      uri,
			UriHash:  uriHash,
			Metadata: metadata,
		},
	})

	// Check output
	require.NoError(t, err)
	require.Equal(t, classId, res.Nft.ClassId)
	require.Equal(t, uri, res.Nft.Uri)
	require.Equal(t, uriHash, res.Nft.UriHash)

	var nftData types.NFTData
	err = nftData.Unmarshal(res.Nft.Data.Value)
	require.NoErrorf(t, err, "Error unmarshal class data")
	require.Equal(t, types.JsonInput(nil), nftData.Metadata) // TODO: update metadata to support templates
	require.Equal(t, iscnId.Prefix.String(), nftData.ClassParent.IscnIdPrefix)
	require.Equal(t, classIscnVersionAtMint, nftData.ClassParent.IscnVersionAtMint)

	// Check mock was called as expected
	ctrl.Finish()
}

func TestMintBlindBoxNFTChangingClaimPeriodPrice(t *testing.T) {
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
	ctx = ctx.WithBlockHeader(tmproto.Header{Time: *testutil.MustParseTime(time.RFC3339, "2022-04-20T00:00:00Z")})
	updatedGoCtx := sdk.WrapSDKContext(ctx)

	// Test Input
	ownerAddressBytes := []byte{0, 1, 0, 1, 0, 1, 0, 1}
	ownerAddress, _ := sdk.Bech32ifyAddressBytes("cosmos", ownerAddressBytes)
	minterAddressBytes := []byte{0, 1, 0, 1, 0, 1, 1, 0, 0, 1, 1}
	minterAddress, _ := sdk.Bech32ifyAddressBytes("cosmos", minterAddressBytes)
	iscnId := iscntypes.NewIscnId("likecoin-chain", "abcdef", 1)
	classId := "likenft1aabbccddeeff"
	tokenId := "token1"
	uri := "ipfs://a1b2c3"
	uriHash := "a1b2c3"
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
	firstClaimPeriodPrice := uint64(10000000000)
	secondClaimPeriodPrice := uint64(20000000000)
	thirdClaimPeriodPrice := uint64(30000000000)

	revealTime := testutil.MustParseTime(time.RFC3339, "2322-04-20T00:00:00Z")

	// Mock keeper calls
	classIscnVersionAtMint := uint64(1)
	classData := types.ClassData{
		Metadata: types.JsonInput(`{"aaaa": "bbbb"}`),
		Parent: types.ClassParent{
			Type:              types.ClassParentType_ISCN,
			IscnIdPrefix:      iscnId.Prefix.String(),
			IscnVersionAtMint: classIscnVersionAtMint,
		},
		Config: types.ClassConfig{
			Burnable:       false,
			MaxSupply:      uint64(5000000000),
			EnableBlindBox: true,
			// Assume sorted descending when creating/updating class
			ClaimPeriods: []*types.ClaimPeriod{
				{
					StartTime:        testutil.MustParseTime(time.RFC3339, "2048-04-01T00:00:00Z"),
					AllowedAddresses: nil,
					MintPrice:        thirdClaimPeriodPrice,
				},
				{
					StartTime:        testutil.MustParseTime(time.RFC3339, "2022-04-01T00:00:00Z"),
					AllowedAddresses: nil,
					MintPrice:        secondClaimPeriodPrice,
				},
				{
					StartTime:        testutil.MustParseTime(time.RFC3339, "2022-01-01T00:00:00Z"),
					AllowedAddresses: nil,
					MintPrice:        firstClaimPeriodPrice,
				},
			},
			RevealTime: revealTime,
		},
	}
	classDataInAny, _ := cdctypes.NewAnyWithValue(&classData)
	nftKeeper.
		EXPECT().
		GetClass(gomock.Any(), gomock.Eq(classId)).
		Return(nft.Class{
			Id:          classId,
			Name:        "Class Name",
			Symbol:      "ABC",
			Description: "Testing Class 123",
			Uri:         "ipfs://abcdef",
			UriHash:     "abcdef",
			Data:        classDataInAny,
		}, true)

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

	// Test for subsequent nft mint at this case
	// No class update
	nftKeeper.
		EXPECT().
		GetTotalSupply(gomock.Any(), gomock.Eq(classId)).
		Return(uint64(1))

	wrappedMinterAddress, _ := sdk.AccAddressFromBech32(minterAddress)
	nftKeeper.
		EXPECT().
		Mint(gomock.Any(), gomock.Any(), wrappedMinterAddress)

	bankKeeper.EXPECT().GetBalance(
		gomock.Any(),
		gomock.Eq(wrappedMinterAddress),
		gomock.Any(),
	).
		Return(sdk.Coin{
			Denom:  "nanolike",
			Amount: sdk.NewInt(100000000000),
		})

	wrappedOwnerAddress, _ := sdk.AccAddressFromBech32(ownerAddress)
	bankKeeper.EXPECT().SendCoins(
		gomock.Any(),
		gomock.Eq(wrappedMinterAddress),
		gomock.Eq(wrappedOwnerAddress),
		gomock.Eq(sdk.Coins{sdk.Coin{
			Denom:  "nanolike",
			Amount: sdk.NewInt(int64(secondClaimPeriodPrice)),
		}}),
	)

	// Run
	res, err := msgServer.MintNFT(updatedGoCtx, &types.MsgMintNFT{
		Creator: minterAddress,
		ClassId: classId,
		Id:      tokenId,
		Input: types.NFTInput{
			Uri:      uri,
			UriHash:  uriHash,
			Metadata: metadata,
		},
	})

	// Check output
	require.NoError(t, err)
	require.Equal(t, classId, res.Nft.ClassId)
	require.Equal(t, uri, res.Nft.Uri)
	require.Equal(t, uriHash, res.Nft.UriHash)

	var nftData types.NFTData
	err = nftData.Unmarshal(res.Nft.Data.Value)
	require.NoErrorf(t, err, "Error unmarshal class data")
	require.Equal(t, types.JsonInput(nil), nftData.Metadata) // TODO: update metadata to support templates
	require.Equal(t, iscnId.Prefix.String(), nftData.ClassParent.IscnIdPrefix)
	require.Equal(t, classIscnVersionAtMint, nftData.ClassParent.IscnVersionAtMint)

	// Check mock was called as expected
	ctrl.Finish()
}

func TestMintBlindBoxNFTFree(t *testing.T) {
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
	ctx = ctx.WithBlockHeader(tmproto.Header{Time: *testutil.MustParseTime(time.RFC3339, "2022-04-20T00:00:00Z")})
	updatedGoCtx := sdk.WrapSDKContext(ctx)

	// Test Input
	ownerAddressBytes := []byte{0, 1, 0, 1, 0, 1, 0, 1}
	minterAddressBytes := []byte{0, 1, 0, 1, 0, 1, 1, 0, 0, 1, 1}
	minterAddress, _ := sdk.Bech32ifyAddressBytes("cosmos", minterAddressBytes)
	iscnId := iscntypes.NewIscnId("likecoin-chain", "abcdef", 1)
	classId := "likenft1aabbccddeeff"
	tokenId := "token1"
	uri := "ipfs://a1b2c3"
	uriHash := "a1b2c3"
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
	mintPrice := uint64(0)
	revealTime := testutil.MustParseTime(time.RFC3339, "2322-04-20T00:00:00Z")

	// Mock keeper calls
	classIscnVersionAtMint := uint64(1)
	classData := types.ClassData{
		Metadata: types.JsonInput(`{"aaaa": "bbbb"}`),
		Parent: types.ClassParent{
			Type:              types.ClassParentType_ISCN,
			IscnIdPrefix:      iscnId.Prefix.String(),
			IscnVersionAtMint: classIscnVersionAtMint,
		},
		Config: types.ClassConfig{
			Burnable:       false,
			MaxSupply:      uint64(5000000000),
			EnableBlindBox: true,
			ClaimPeriods: []*types.ClaimPeriod{
				{
					StartTime:        testutil.MustParseTime(time.RFC3339, "2020-01-01T00:00:00Z"),
					AllowedAddresses: nil,
					MintPrice:        mintPrice,
				},
			},
			RevealTime: revealTime,
		},
	}
	classDataInAny, _ := cdctypes.NewAnyWithValue(&classData)
	nftKeeper.
		EXPECT().
		GetClass(gomock.Any(), gomock.Eq(classId)).
		Return(nft.Class{
			Id:          classId,
			Name:        "Class Name",
			Symbol:      "ABC",
			Description: "Testing Class 123",
			Uri:         "ipfs://abcdef",
			UriHash:     "abcdef",
			Data:        classDataInAny,
		}, true)

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

	// Test for subsequent nft mint at this case
	// No class update
	nftKeeper.
		EXPECT().
		GetTotalSupply(gomock.Any(), gomock.Eq(classId)).
		Return(uint64(1))

	wrappedMinterAddress, _ := sdk.AccAddressFromBech32(minterAddress)
	nftKeeper.
		EXPECT().
		Mint(gomock.Any(), gomock.Any(), wrappedMinterAddress)

	// Run
	res, err := msgServer.MintNFT(updatedGoCtx, &types.MsgMintNFT{
		Creator: minterAddress,
		ClassId: classId,
		Id:      tokenId,
		Input: types.NFTInput{
			Uri:      uri,
			UriHash:  uriHash,
			Metadata: metadata,
		},
	})

	// Check output
	require.NoError(t, err)
	require.Equal(t, classId, res.Nft.ClassId)
	require.Equal(t, uri, res.Nft.Uri)
	require.Equal(t, uriHash, res.Nft.UriHash)

	var nftData types.NFTData
	err = nftData.Unmarshal(res.Nft.Data.Value)
	require.NoErrorf(t, err, "Error unmarshal class data")
	require.Equal(t, types.JsonInput(nil), nftData.Metadata) // TODO: update metadata to support templates
	require.Equal(t, iscnId.Prefix.String(), nftData.ClassParent.IscnIdPrefix)
	require.Equal(t, classIscnVersionAtMint, nftData.ClassParent.IscnVersionAtMint)

	// Check mock was called as expected
	ctrl.Finish()
}

func TestMintBlindBoxNFTInsufficientFunds(t *testing.T) {
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
	ctx = ctx.WithBlockHeader(tmproto.Header{Time: *testutil.MustParseTime(time.RFC3339, "2022-04-20T00:00:00Z")})
	updatedGoCtx := sdk.WrapSDKContext(ctx)

	// Test Input
	ownerAddressBytes := []byte{0, 1, 0, 1, 0, 1, 0, 1}
	minterAddressBytes := []byte{0, 1, 0, 1, 0, 1, 1, 0, 0, 1, 1}
	minterAddress, _ := sdk.Bech32ifyAddressBytes("cosmos", minterAddressBytes)
	iscnId := iscntypes.NewIscnId("likecoin-chain", "abcdef", 1)
	classId := "likenft1aabbccddeeff"
	tokenId := "token1"
	uri := "ipfs://a1b2c3"
	uriHash := "a1b2c3"
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
	mintPrice := uint64(5000000000)
	revealTime := testutil.MustParseTime(time.RFC3339, "2322-04-20T00:00:00Z")

	// Mock keeper calls
	classIscnVersionAtMint := uint64(1)
	classData := types.ClassData{
		Metadata: types.JsonInput(`{"aaaa": "bbbb"}`),
		Parent: types.ClassParent{
			Type:              types.ClassParentType_ISCN,
			IscnIdPrefix:      iscnId.Prefix.String(),
			IscnVersionAtMint: classIscnVersionAtMint,
		},
		Config: types.ClassConfig{
			Burnable:       false,
			MaxSupply:      uint64(5000000000),
			EnableBlindBox: true,
			ClaimPeriods: []*types.ClaimPeriod{
				{
					StartTime:        testutil.MustParseTime(time.RFC3339, "2020-01-01T00:00:00Z"),
					AllowedAddresses: nil,
					MintPrice:        mintPrice,
				},
			},
			RevealTime: revealTime,
		},
	}
	classDataInAny, _ := cdctypes.NewAnyWithValue(&classData)
	nftKeeper.
		EXPECT().
		GetClass(gomock.Any(), gomock.Eq(classId)).
		Return(nft.Class{
			Id:          classId,
			Name:        "Class Name",
			Symbol:      "ABC",
			Description: "Testing Class 123",
			Uri:         "ipfs://abcdef",
			UriHash:     "abcdef",
			Data:        classDataInAny,
		}, true)

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

	// Test for subsequent nft mint at this case
	// No class update
	nftKeeper.
		EXPECT().
		GetTotalSupply(gomock.Any(), gomock.Eq(classId)).
		Return(uint64(1))

	wrappedMinterAddress, _ := sdk.AccAddressFromBech32(minterAddress)

	bankKeeper.EXPECT().GetBalance(
		gomock.Any(),
		gomock.Eq(wrappedMinterAddress),
		gomock.Any(),
	).
		Return(sdk.Coin{
			Denom:  "nanolike",
			Amount: sdk.NewInt(100),
		})

	// Run
	res, err := msgServer.MintNFT(updatedGoCtx, &types.MsgMintNFT{
		Creator: minterAddress,
		ClassId: classId,
		Id:      tokenId,
		Input: types.NFTInput{
			Uri:      uri,
			UriHash:  uriHash,
			Metadata: metadata,
		},
	})

	// Check output
	require.Error(t, err)
	require.Contains(t, err.Error(), types.ErrInsufficientFunds.Error())
	require.Nil(t, res)

	// Check mock was called as expected
	ctrl.Finish()
}

func TestMintNFTNotOwnerNoBlindBox(t *testing.T) {
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
	ctx = ctx.WithBlockHeader(tmproto.Header{Time: *testutil.MustParseTime(time.RFC3339, "2022-04-20T00:00:00Z")})
	updatedGoCtx := sdk.WrapSDKContext(ctx)

	// Test Input
	ownerAddressBytes := []byte{0, 1, 0, 1, 0, 1, 0, 1}
	minterAddressBytes := []byte{0, 1, 0, 1, 0, 1, 1, 0, 0, 1, 1}
	minterAddress, _ := sdk.Bech32ifyAddressBytes("cosmos", minterAddressBytes)
	iscnId := iscntypes.NewIscnId("likecoin-chain", "abcdef", 1)
	classId := "likenft1aabbccddeeff"
	tokenId := "token1"
	uri := "ipfs://a1b2c3"
	uriHash := "a1b2c3"
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

	// Mock keeper calls
	classIscnVersionAtMint := uint64(1)
	classData := types.ClassData{
		Metadata: types.JsonInput(`{"aaaa": "bbbb"}`),
		Parent: types.ClassParent{
			Type:              types.ClassParentType_ISCN,
			IscnIdPrefix:      iscnId.Prefix.String(),
			IscnVersionAtMint: classIscnVersionAtMint,
		},
		Config: types.ClassConfig{
			Burnable:       false,
			MaxSupply:      uint64(5000000000),
			EnableBlindBox: false,
		},
	}
	classDataInAny, _ := cdctypes.NewAnyWithValue(&classData)
	nftKeeper.
		EXPECT().
		GetClass(gomock.Any(), gomock.Eq(classId)).
		Return(nft.Class{
			Id:          classId,
			Name:        "Class Name",
			Symbol:      "ABC",
			Description: "Testing Class 123",
			Uri:         "ipfs://abcdef",
			UriHash:     "abcdef",
			Data:        classDataInAny,
		}, true)

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

	// Test for subsequent nft mint at this case
	// No class update
	nftKeeper.
		EXPECT().
		GetTotalSupply(gomock.Any(), gomock.Eq(classId)).
		Return(uint64(1))

	// Run
	res, err := msgServer.MintNFT(updatedGoCtx, &types.MsgMintNFT{
		Creator: minterAddress,
		ClassId: classId,
		Id:      tokenId,
		Input: types.NFTInput{
			Uri:      uri,
			UriHash:  uriHash,
			Metadata: metadata,
		},
	})

	// Check output
	require.Error(t, err)
	require.Contains(t, err.Error(), sdkerrors.ErrUnauthorized.Error())
	require.Nil(t, res)

	// Check mock was called as expected
	ctrl.Finish()
}

func TestMintNFTNoSupply(t *testing.T) {
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
	ctx = ctx.WithBlockHeader(tmproto.Header{Time: *testutil.MustParseTime(time.RFC3339, "2022-04-20T00:00:00Z")})
	updatedGoCtx := sdk.WrapSDKContext(ctx)

	// Test Input 1
	ownerAddressBytes := []byte{0, 1, 0, 1, 0, 1, 0, 1}
	ownerAddress, _ := sdk.Bech32ifyAddressBytes("cosmos", ownerAddressBytes)
	iscnId := iscntypes.NewIscnId("likecoin-chain", "abcdef", 1)
	classId := "likenft1aabbccddeeff"
	tokenId := "token1"
	uri := "ipfs://a1b2c3"
	uriHash := "a1b2c3"
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

	// Mock keeper calls
	classIscnVersionAtMint := uint64(1)
	classData := types.ClassData{
		Metadata: types.JsonInput(`{"aaaa": "bbbb"}`),
		Parent: types.ClassParent{
			Type:              types.ClassParentType_ISCN,
			IscnIdPrefix:      iscnId.Prefix.String(),
			IscnVersionAtMint: classIscnVersionAtMint,
		},
		Config: types.ClassConfig{
			Burnable:  false,
			MaxSupply: uint64(500),
		},
	}
	classDataInAny, _ := cdctypes.NewAnyWithValue(&classData)
	nftKeeper.
		EXPECT().
		GetClass(gomock.Any(), gomock.Eq(classId)).
		Return(nft.Class{
			Id:          classId,
			Name:        "Class Name",
			Symbol:      "ABC",
			Description: "Testing Class 123",
			Uri:         "ipfs://abcdef",
			UriHash:     "abcdef",
			Data:        classDataInAny,
		}, true)

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

	// Test for subsequent nft mint at this case
	// No class update
	nftKeeper.
		EXPECT().
		GetTotalSupply(gomock.Any(), gomock.Eq(classId)).
		Return(uint64(500))

	// Run
	res, err := msgServer.MintNFT(updatedGoCtx, &types.MsgMintNFT{
		Creator: ownerAddress,
		ClassId: classId,
		Id:      tokenId,
		Input: types.NFTInput{
			Uri:      uri,
			UriHash:  uriHash,
			Metadata: metadata,
		},
	})

	// Check output
	require.Error(t, err)
	require.Contains(t, err.Error(), types.ErrNftNoSupply.Error())
	require.Nil(t, res)

	// Check mock was called as expected
	ctrl.Finish()
}

func TestMintBlindBoxNFTAfterRevealTime(t *testing.T) {
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
	ctx = ctx.WithBlockHeader(tmproto.Header{Time: *testutil.MustParseTime(time.RFC3339, "2022-04-20T00:00:00Z")})
	updatedGoCtx := sdk.WrapSDKContext(ctx)

	// Test Input 1
	ownerAddressBytes := []byte{0, 1, 0, 1, 0, 1, 0, 1}
	minterAddressBytes := []byte{0, 1, 0, 1, 0, 1, 1, 0, 0, 1, 1}
	minterAddress, _ := sdk.Bech32ifyAddressBytes("cosmos", minterAddressBytes)
	iscnId := iscntypes.NewIscnId("likecoin-chain", "abcdef", 1)
	classId := "likenft1aabbccddeeff"
	tokenId := "token1"
	uri := "ipfs://a1b2c3"
	uriHash := "a1b2c3"
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
	revealTime := testutil.MustParseTime(time.RFC3339, "2020-01-01T00:00:00Z")

	// Mock keeper calls
	classIscnVersionAtMint := uint64(1)
	classData := types.ClassData{
		Metadata: types.JsonInput(`{"aaaa": "bbbb"}`),
		Parent: types.ClassParent{
			Type:              types.ClassParentType_ISCN,
			IscnIdPrefix:      iscnId.Prefix.String(),
			IscnVersionAtMint: classIscnVersionAtMint,
		},
		Config: types.ClassConfig{
			Burnable:       false,
			MaxSupply:      uint64(500),
			EnableBlindBox: true,
			ClaimPeriods: []*types.ClaimPeriod{
				{
					StartTime:        testutil.MustParseTime(time.RFC3339, "2020-01-01T00:00:00Z"),
					AllowedAddresses: nil,
					MintPrice:        uint64(0),
				},
			},
			RevealTime: revealTime,
		},
	}
	classDataInAny, _ := cdctypes.NewAnyWithValue(&classData)
	nftKeeper.
		EXPECT().
		GetClass(gomock.Any(), gomock.Eq(classId)).
		Return(nft.Class{
			Id:          classId,
			Name:        "Class Name",
			Symbol:      "ABC",
			Description: "Testing Class 123",
			Uri:         "ipfs://abcdef",
			UriHash:     "abcdef",
			Data:        classDataInAny,
		}, true)

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

	// Test for subsequent nft mint at this case
	// No class update
	nftKeeper.
		EXPECT().
		GetTotalSupply(gomock.Any(), gomock.Eq(classId)).
		Return(uint64(2))

	// Run
	res, err := msgServer.MintNFT(updatedGoCtx, &types.MsgMintNFT{
		Creator: minterAddress,
		ClassId: classId,
		Id:      tokenId,
		Input: types.NFTInput{
			Uri:      uri,
			UriHash:  uriHash,
			Metadata: metadata,
		},
	})

	// Check output
	require.Error(t, err)
	require.Contains(t, err.Error(), types.ErrFailedToMintNFT.Error())
	require.Nil(t, res)

	// Check mock was called as expected
	ctrl.Finish()
}

func TestMintNFTUnlimitedSupply(t *testing.T) {
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
	ctx = ctx.WithBlockHeader(tmproto.Header{Time: *testutil.MustParseTime(time.RFC3339, "2022-04-20T00:00:00Z")})
	updatedGoCtx := sdk.WrapSDKContext(ctx)

	// Test Input
	ownerAddressBytes := []byte{0, 1, 0, 1, 0, 1, 0, 1}
	ownerAddress, _ := sdk.Bech32ifyAddressBytes("cosmos", ownerAddressBytes)
	iscnId := iscntypes.NewIscnId("likecoin-chain", "abcdef", 1)
	classId := "likenft1aabbccddeeff"
	tokenId := "token1"
	uri := "ipfs://a1b2c3"
	uriHash := "a1b2c3"
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

	// Mock keeper calls
	classIscnVersionAtMint := uint64(1)
	classData := types.ClassData{
		Metadata: types.JsonInput(`{"aaaa": "bbbb"}`),
		Parent: types.ClassParent{
			Type:              types.ClassParentType_ISCN,
			IscnIdPrefix:      iscnId.Prefix.String(),
			IscnVersionAtMint: classIscnVersionAtMint,
		},
		Config: types.ClassConfig{
			Burnable:  false,
			MaxSupply: uint64(0),
		},
	}
	classDataInAny, _ := cdctypes.NewAnyWithValue(&classData)
	nftKeeper.
		EXPECT().
		GetClass(gomock.Any(), gomock.Eq(classId)).
		Return(nft.Class{
			Id:          classId,
			Name:        "Class Name",
			Symbol:      "ABC",
			Description: "Testing Class 123",
			Uri:         "ipfs://abcdef",
			UriHash:     "abcdef",
			Data:        classDataInAny,
		}, true)

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

	// Test for subsequent nft mint at this case
	// No class update
	nftKeeper.
		EXPECT().
		GetTotalSupply(gomock.Any(), gomock.Eq(classId)).
		Return(uint64(50))

	wrappedOwnerAddress, _ := sdk.AccAddressFromBech32(ownerAddress)
	nftKeeper.
		EXPECT().
		Mint(gomock.Any(), gomock.Any(), wrappedOwnerAddress)

	// Run
	res, err := msgServer.MintNFT(updatedGoCtx, &types.MsgMintNFT{
		Creator: ownerAddress,
		ClassId: classId,
		Id:      tokenId,
		Input: types.NFTInput{
			Uri:      uri,
			UriHash:  uriHash,
			Metadata: metadata,
		},
	})

	// Check output
	require.NoError(t, err)
	require.Equal(t, classId, res.Nft.ClassId)
	require.Equal(t, uri, res.Nft.Uri)
	require.Equal(t, uriHash, res.Nft.UriHash)

	var nftData types.NFTData
	err = nftData.Unmarshal(res.Nft.Data.Value)
	require.NoErrorf(t, err, "Error unmarshal class data")
	require.Equal(t, metadata, nftData.Metadata)
	require.Equal(t, iscnId.Prefix.String(), nftData.ClassParent.IscnIdPrefix)
	require.Equal(t, classIscnVersionAtMint, nftData.ClassParent.IscnVersionAtMint)

	// Check mock was called as expected
	ctrl.Finish()
}
