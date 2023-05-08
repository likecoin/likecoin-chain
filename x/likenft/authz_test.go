package likenft_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/authz"

	iscntypes "github.com/likecoin/likecoin-chain/v3/x/iscn/types"
	"github.com/likecoin/likecoin-chain/v3/x/likenft/types"

	testutil "github.com/likecoin/likecoin-chain/v3/testutil"
)

var (
	priv1 = secp256k1.GenPrivKey()
	addr1 = sdk.AccAddress(priv1.PubKey().Address())
	priv2 = secp256k1.GenPrivKey()
	addr2 = sdk.AccAddress(priv2.PubKey().Address())
	priv3 = secp256k1.GenPrivKey()
	addr3 = sdk.AccAddress(priv3.PubKey().Address())

	fingerprint1 = "hash://sha256/9564b85669d5e96ac969dd0161b8475bbced9e5999c6ec598da718a3045d6f2e"

	stakeholder1 = iscntypes.IscnInput(`
{
	"entity": {
		"@id": "did:cosmos:5sy29r37gfxvxz21rh4r0ktpuc46pzjrmz29g45",
		"name": "Chung Wu"
	},
	"rewardProportion": 95,
	"contributionType": "http://schema.org/author"
}`)

	stakeholder2 = iscntypes.IscnInput(`
{
	"rewardProportion": 5,
	"contributionType": "http://schema.org/citation",
	"footprint": "https://en.wikipedia.org/wiki/Fibonacci_number",
	"description": "The blog post referred the matrix form of computing Fibonacci numbers."
}`)

	contentMetadata1 = iscntypes.IscnInput(`
{
	"@context": "http://schema.org/",
	"@type": "CreativeWorks",
	"title": "使用矩陣計算遞歸關係式",
	"description": "An article on computing recursive function with matrix multiplication.",
	"datePublished": "2019-04-19",
	"version": 1,
	"url": "https://nnkken.github.io/post/recursive-relation/",
	"author": "https://github.com/nnkken",
	"usageInfo": "https://creativecommons.org/licenses/by/4.0",
	"keywords": "matrix,recursion"
}`)
)

type authorizationTestSetupClass struct {
	ClassId string
	NftIds  []string
}

type authorizationTestSetupIscn struct {
	IscnId  iscntypes.IscnId
	Record  iscntypes.IscnRecord
	Classes []authorizationTestSetupClass
}

type authorizationTestSetupOwner struct {
	PrivKey *secp256k1.PrivKey
	Addr    sdk.AccAddress
	Iscns   []authorizationTestSetupIscn
}

type authorizationTestSetup struct {
	App        *testutil.TestingApp
	Owners     []authorizationTestSetupOwner
	OtherAddrs []authorizationTestSetupOwner
}

func setupAppAndNfts(t *testing.T) authorizationTestSetup {
	var msg sdk.Msg

	app := testutil.SetupTestApp([]testutil.GenesisBalance{
		{addr1.String(), "1000000000000000000nanolike"},
		{addr2.String(), "1000000000000000000nanolike"},
		{addr3.String(), "1000000000000000000nanolike"},
	})

	baseRecord := iscntypes.IscnRecord{
		ContentFingerprints: []string{fingerprint1},
		Stakeholders:        []iscntypes.IscnInput{stakeholder1, stakeholder2},
		ContentMetadata:     contentMetadata1,
	}
	nftIds := []string{"test-nft-id-a", "test-nft-id-b"}

	app.NextHeader(1234567890)
	app.SetForTx()

	setup := authorizationTestSetup{
		App:        app,
		OtherAddrs: []authorizationTestSetupOwner{{PrivKey: priv3, Addr: addr3}},
	}
	for i, privKey := range []*secp256k1.PrivKey{priv1, priv2} {
		owner := authorizationTestSetupOwner{
			PrivKey: privKey,
			Addr:    sdk.AccAddress(privKey.PubKey().Address()),
		}
		for j := 0; j < 2; j++ {
			record := baseRecord
			record.RecordNotes = fmt.Sprintf("%d-%d", i, j)
			msg = iscntypes.NewMsgCreateIscnRecord(owner.Addr, &record, 0)
			res := app.DeliverMsgNoError(t, msg, owner.PrivKey)
			iscnId := testutil.GetIscnIdFromResult(t, res)
			iscn := authorizationTestSetupIscn{
				IscnId: iscnId,
				Record: record,
			}
			for k := 0; k < 2; k++ {
				msg := types.NewMsgNewClass(
					owner.Addr.String(),
					types.ClassParentInput{
						Type:         types.ClassParentType_ISCN,
						IscnIdPrefix: iscnId.Prefix.String(),
					},
					types.ClassInput{
						Name:     fmt.Sprintf("testclass%d-%d-%d", i, j, k),
						Symbol:   fmt.Sprintf("TEST%d-%d-%d", i, j, k),
						Metadata: types.JsonInput(`{}`),
					},
				)
				res = app.DeliverMsgNoError(t, msg, owner.PrivKey)
				classId := testutil.GetClassIdFromResult(t, res)
				class := authorizationTestSetupClass{
					ClassId: classId,
					NftIds:  nftIds,
				}
				for _, nftId := range nftIds {
					msg := types.NewMsgMintNFT(owner.Addr.String(), classId, nftId, &types.NFTInput{})
					app.DeliverMsgNoError(t, msg, owner.PrivKey)
				}
				iscn.Classes = append(iscn.Classes, class)
			}
			owner.Iscns = append(owner.Iscns, iscn)
		}
		setup.Owners = append(setup.Owners, owner)
	}
	return setup
}

func TestAuthorization(t *testing.T) {
	var msg sdk.Msg
	app := testutil.SetupTestApp([]testutil.GenesisBalance{
		{addr1.String(), "1000000000000000000nanolike"},
		{addr2.String(), "1000000000000000000nanolike"},
		{addr3.String(), "1000000000000000000nanolike"},
	})

	app.NextHeader(1234567890)
	app.SetForTx()
	record := iscntypes.IscnRecord{
		RecordNotes:         "addr1",
		ContentFingerprints: []string{fingerprint1},
		Stakeholders:        []iscntypes.IscnInput{stakeholder1, stakeholder2},
		ContentMetadata:     contentMetadata1,
	}
	msg = iscntypes.NewMsgCreateIscnRecord(addr1, &record, 0)
	res := app.DeliverMsgNoError(t, msg, priv1)
	iscnId1 := testutil.GetIscnIdFromResult(t, res)

	record = iscntypes.IscnRecord{
		RecordNotes:         "addr3",
		ContentFingerprints: []string{fingerprint1},
		Stakeholders:        []iscntypes.IscnInput{stakeholder1, stakeholder2},
		ContentMetadata:     contentMetadata1,
	}
	msg = iscntypes.NewMsgCreateIscnRecord(addr3, &record, 0)
	res = app.DeliverMsgNoError(t, msg, priv3)
	iscnId3 := testutil.GetIscnIdFromResult(t, res)

	msg = types.NewMsgNewClass(
		addr1.String(),
		types.ClassParentInput{
			Type:         types.ClassParentType_ISCN,
			IscnIdPrefix: iscnId1.Prefix.String(),
		},
		types.ClassInput{
			Name:        "testclass1-1",
			Symbol:      "TEST1-1",
			Description: "test class 1-1",
			Metadata:    types.JsonInput(`{}`),
		},
	)
	res = app.DeliverMsgNoError(t, msg, priv1)
	classId1 := testutil.GetClassIdFromResult(t, res)

	msg = types.NewMsgNewClass(
		addr3.String(),
		types.ClassParentInput{
			Type:         types.ClassParentType_ISCN,
			IscnIdPrefix: iscnId3.Prefix.String(),
		},
		types.ClassInput{
			Name:        "testclass3-1",
			Symbol:      "TEST3-1",
			Description: "test class 3-1",
			Metadata:    types.JsonInput(`{}`),
		},
	)
	res = app.DeliverMsgNoError(t, msg, priv3)
	classId3 := testutil.GetClassIdFromResult(t, res)

	expiration := time.Unix(2000000000, 0)
	msg, err := authz.NewMsgGrant(addr1, addr2, types.NewUpdateClassAuthorization(classId1), &expiration)
	require.NoError(t, err)
	app.DeliverMsgNoError(t, msg, priv1)

	newClassInput := types.ClassInput{
		Name:        "testclass1-2",
		Symbol:      "TEST1-2",
		Description: "test class 1-2",
		Metadata:    types.JsonInput(`{}`),
	}
	updateClassMsg := types.NewMsgUpdateClass(addr1.String(), classId1, newClassInput)
	msgExec := authz.NewMsgExec(addr2, []sdk.Msg{updateClassMsg})
	msg = &msgExec
	app.DeliverMsgNoError(t, msg, priv2)

	ctx := app.SetForQuery()
	class, ok := app.NftKeeper.GetClass(ctx, classId1)
	require.True(t, ok)
	require.Equal(t, newClassInput.Name, class.Name)
	require.Equal(t, newClassInput.Symbol, class.Symbol)
	require.Equal(t, newClassInput.Description, class.Description)

	app.NextHeader(1234567891)
	app.SetForTx()

	newClassInput = types.ClassInput{
		Name:        "testclass1-3",
		Symbol:      "TEST1-3",
		Description: "test class 1-3",
		Metadata:    types.JsonInput(`{}`),
	}
	updateClassMsg = types.NewMsgUpdateClass(addr1.String(), classId1, newClassInput)
	msgExec = authz.NewMsgExec(addr3, []sdk.Msg{updateClassMsg})
	msg = &msgExec
	_, _, simErr, _ := app.DeliverMsg(msg, priv3)
	require.ErrorContains(t, simErr, "authorization not found")

	newClassInput = types.ClassInput{
		Name:        "testclass3-2",
		Symbol:      "TEST3-2",
		Description: "test class 3-2",
		Metadata:    types.JsonInput(`{}`),
	}
	updateClassMsg = types.NewMsgUpdateClass(addr3.String(), classId3, newClassInput)
	msgExec = authz.NewMsgExec(addr2, []sdk.Msg{updateClassMsg})
	msg = &msgExec
	_, _, simErr, _ = app.DeliverMsg(msg, priv2)
	require.ErrorContains(t, simErr, "authorization not found")

	expiration = time.Unix(2000000000, 0)
	msg, err = authz.NewMsgGrant(addr1, addr2, types.NewMintNFTAuthorization(classId1), &expiration)
	require.NoError(t, err)
	app.DeliverMsgNoError(t, msg, priv1)

	msgMint := types.NewMsgMintNFT(addr1.String(), classId1, "token-1-by-2", &types.NFTInput{})
	msgExec = authz.NewMsgExec(addr2, []sdk.Msg{msgMint})
	msg = &msgExec
	app.DeliverMsgNoError(t, msg, priv2)

	msgMint = types.NewMsgMintNFT(addr3.String(), classId3, "token-3-by-3", &types.NFTInput{})
	msgExec = authz.NewMsgExec(addr3, []sdk.Msg{msgMint})
	msg = &msgExec
	app.DeliverMsgNoError(t, msg, priv3)

	msgMint = types.NewMsgMintNFT(addr3.String(), classId3, "token-3-by-2", &types.NFTInput{})
	msgExec = authz.NewMsgExec(addr2, []sdk.Msg{msgMint})
	msg = &msgExec
	_, _, simErr, _ = app.DeliverMsg(msg, priv2)
	require.ErrorContains(t, simErr, "authorization not found")

	msgMint = types.NewMsgMintNFT(addr1.String(), classId1, "token-1-by-3", &types.NFTInput{})
	msgExec = authz.NewMsgExec(addr3, []sdk.Msg{msgMint})
	msg = &msgExec
	_, _, simErr, _ = app.DeliverMsg(msg, priv3)
	require.ErrorContains(t, simErr, "authorization not found")

	msg = types.NewMsgCreateRoyaltyConfig(addr1.String(), classId1, types.RoyaltyConfigInput{
		RateBasisPoints: 100,
		Stakeholders: []types.RoyaltyStakeholderInput{
			{Account: addr1.String(), Weight: 70},
			{Account: addr2.String(), Weight: 30},
		},
	})
	app.DeliverMsgNoError(t, msg, priv1)

	msg = types.NewMsgCreateRoyaltyConfig(addr3.String(), classId3, types.RoyaltyConfigInput{
		RateBasisPoints: 200,
		Stakeholders: []types.RoyaltyStakeholderInput{
			{Account: addr3.String(), Weight: 120},
			{Account: addr2.String(), Weight: 80},
		},
	})
	app.DeliverMsgNoError(t, msg, priv3)

	expiration = time.Unix(2000000000, 0)
	msg, err = authz.NewMsgGrant(addr1, addr2, types.NewUpdateRoyaltyConfigAuthorization(classId1), &expiration)
	require.NoError(t, err)
	app.DeliverMsgNoError(t, msg, priv1)

	newRoyaltyConfig := types.RoyaltyConfigInput{
		RateBasisPoints: 300,
		Stakeholders: []types.RoyaltyStakeholderInput{
			{Account: addr2.String(), Weight: 100},
		},
	}
	msgUpdateRoyaltyConfig := types.NewMsgUpdateRoyaltyConfig(addr1.String(), classId1, newRoyaltyConfig)
	msgExec = authz.NewMsgExec(addr2, []sdk.Msg{msgUpdateRoyaltyConfig})
	msg = &msgExec
	app.DeliverMsgNoError(t, msg, priv2)

	ctx = app.SetForQuery()
	config, ok := app.LikeNftKeeper.GetRoyaltyConfig(ctx, classId1)
	require.True(t, ok)
	require.Equal(t, newRoyaltyConfig.RateBasisPoints, config.RateBasisPoints)
	require.Equal(t, len(newRoyaltyConfig.Stakeholders), len(config.Stakeholders))
	for i, stakeholder := range config.Stakeholders {
		require.Equal(t, stakeholder.Account.String(), newRoyaltyConfig.Stakeholders[i].Account)
		require.Equal(t, stakeholder.Weight, newRoyaltyConfig.Stakeholders[i].Weight)
	}

	app.NextHeader(1234567892)
	app.SetForTx()

	msgUpdateRoyaltyConfig = types.NewMsgUpdateRoyaltyConfig(addr1.String(), classId1, newRoyaltyConfig)
	msgExec = authz.NewMsgExec(addr3, []sdk.Msg{msgUpdateRoyaltyConfig})
	msg = &msgExec
	_, _, simErr, _ = app.DeliverMsg(msg, priv3)
	require.ErrorContains(t, simErr, "authorization not found")

	msgUpdateRoyaltyConfig = types.NewMsgUpdateRoyaltyConfig(addr3.String(), classId1, newRoyaltyConfig)
	msgExec = authz.NewMsgExec(addr2, []sdk.Msg{msgUpdateRoyaltyConfig})
	msg = &msgExec
	_, _, simErr, _ = app.DeliverMsg(msg, priv2)
	require.ErrorContains(t, simErr, "authorization not found")

	msg = types.NewMsgCreateOffer(addr1.String(), classId1, "token-1-by-2", 1, time.Unix(1240000000, 0))
	app.DeliverMsgNoError(t, msg, priv1)

	msg = types.NewMsgCreateOffer(addr3.String(), classId3, "token-3-by-3", 1, time.Unix(1240000000, 0))
	app.DeliverMsgNoError(t, msg, priv3)

	expiration = time.Unix(2000000000, 0)
	msg, err = authz.NewMsgGrant(addr1, addr2, types.NewUpdateOfferAuthorization(classId1, "token-1-by-2"), &expiration)
	require.NoError(t, err)
	app.DeliverMsgNoError(t, msg, priv1)

	msgUpdateOffer := types.NewMsgUpdateOffer(addr1.String(), classId1, "token-1-by-2", 2, time.Unix(1240000001, 0))
	msgExec = authz.NewMsgExec(addr2, []sdk.Msg{msgUpdateOffer})
	msg = &msgExec
	app.DeliverMsgNoError(t, msg, priv2)

	msgUpdateOffer = types.NewMsgUpdateOffer(addr1.String(), classId1, "token-1-by-2", 3, time.Unix(1240000002, 0))
	msgExec = authz.NewMsgExec(addr3, []sdk.Msg{msgUpdateOffer})
	msg = &msgExec
	_, _, simErr, _ = app.DeliverMsg(msg, priv3)
	require.ErrorContains(t, simErr, "authorization not found")

	msgUpdateOffer = types.NewMsgUpdateOffer(addr3.String(), classId1, "token-3-by-3", 3, time.Unix(1240000002, 0))
	msgExec = authz.NewMsgExec(addr2, []sdk.Msg{msgUpdateOffer})
	msg = &msgExec
	_, _, simErr, _ = app.DeliverMsg(msg, priv2)
	require.ErrorContains(t, simErr, "authorization not found")

	msg = types.NewMsgCreateListing(addr1.String(), classId1, "token-1-by-2", 1, time.Unix(1240000000, 0), false)
	app.DeliverMsgNoError(t, msg, priv1)

	msg = types.NewMsgCreateListing(addr3.String(), classId3, "token-3-by-3", 1, time.Unix(1240000000, 0), false)
	app.DeliverMsgNoError(t, msg, priv3)

	expiration = time.Unix(2000000000, 0)
	msg, err = authz.NewMsgGrant(addr1, addr2, types.NewUpdateListingAuthorization(classId1, "token-1-by-2"), &expiration)
	require.NoError(t, err)
	app.DeliverMsgNoError(t, msg, priv1)

	msgUpdateListing := types.NewMsgUpdateListing(addr1.String(), classId1, "token-1-by-2", 2, time.Unix(1240000001, 0), false)
	msgExec = authz.NewMsgExec(addr2, []sdk.Msg{msgUpdateListing})
	msg = &msgExec
	app.DeliverMsgNoError(t, msg, priv2)

	msgUpdateListing = types.NewMsgUpdateListing(addr1.String(), classId1, "token-1-by-2", 3, time.Unix(1240000002, 0), false)
	msgExec = authz.NewMsgExec(addr3, []sdk.Msg{msgUpdateListing})
	msg = &msgExec
	_, _, simErr, _ = app.DeliverMsg(msg, priv3)
	require.ErrorContains(t, simErr, "authorization not found")

	msgUpdateListing = types.NewMsgUpdateListing(addr3.String(), classId1, "token-3-by-3", 3, time.Unix(1240000002, 0), false)
	msgExec = authz.NewMsgExec(addr2, []sdk.Msg{msgUpdateListing})
	msg = &msgExec
	_, _, simErr, _ = app.DeliverMsg(msg, priv2)
	require.ErrorContains(t, simErr, "authorization not found")
}
