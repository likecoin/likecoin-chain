package likenft_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/authz"

	"github.com/likecoin/likecoin-chain/v4/x/likenft/types"
)

func TestOfferAuthorizations(t *testing.T) {
	var msg sdk.Msg
	var msgGrant *authz.MsgGrant
	var msgExec authz.MsgExec
	var err error

	setup := setupAppAndNfts(t)
	app := setup.App

	granter := setup.Owners[0]
	grantee := setup.OtherAddrs[0]

	grantedClassId := setup.Owners[0].Iscns[0].Classes[0].ClassId
	ungrantedClassId := setup.Owners[0].Iscns[0].Classes[1].ClassId
	grantedNftId := setup.Owners[0].Iscns[0].Classes[0].NftIds[0]
	ungrantedNftId := setup.Owners[0].Iscns[0].Classes[0].NftIds[1]

	expiration := time.Unix(1300000000, 0)
	msgGrant, err = authz.NewMsgGrant(granter.Addr, grantee.Addr, &types.CreateOfferAuthorization{
		ClassId: grantedClassId,
		NftId:   grantedNftId,
	}, &expiration)
	require.NoError(t, err)
	app.DeliverMsgNoError(t, msgGrant, granter.PrivKey)

	msg = types.NewMsgCreateOffer(granter.Addr.String(), grantedClassId, grantedNftId, 1, time.Unix(1234569999, 0))
	msgExec = authz.NewMsgExec(grantee.Addr, []sdk.Msg{msg})
	app.DeliverMsgNoError(t, &msgExec, grantee.PrivKey)

	msg = types.NewMsgCreateOffer(granter.Addr.String(), grantedClassId, ungrantedNftId, 1, time.Unix(1234569999, 0))
	msgExec = authz.NewMsgExec(grantee.Addr, []sdk.Msg{msg})
	app.DeliverMsgSimError(t, &msgExec, grantee.PrivKey, "NFT ID mismatch")

	msg = types.NewMsgCreateOffer(granter.Addr.String(), ungrantedClassId, grantedNftId, 1, time.Unix(1234569999, 0))
	msgExec = authz.NewMsgExec(grantee.Addr, []sdk.Msg{msg})
	app.DeliverMsgSimError(t, &msgExec, grantee.PrivKey, "class ID mismatch")

	expiration = time.Unix(1300000000, 0)
	msgGrant, err = authz.NewMsgGrant(granter.Addr, grantee.Addr, &types.UpdateOfferAuthorization{
		ClassId: grantedClassId,
		NftId:   grantedNftId,
	}, &expiration)
	require.NoError(t, err)
	app.DeliverMsgNoError(t, msgGrant, granter.PrivKey)

	msg = types.NewMsgUpdateOffer(granter.Addr.String(), grantedClassId, grantedNftId, 2, time.Unix(1234569999, 0))
	msgExec = authz.NewMsgExec(grantee.Addr, []sdk.Msg{msg})
	app.DeliverMsgNoError(t, &msgExec, grantee.PrivKey)

	msg = types.NewMsgUpdateOffer(granter.Addr.String(), grantedClassId, ungrantedNftId, 2, time.Unix(1234569999, 0))
	msgExec = authz.NewMsgExec(grantee.Addr, []sdk.Msg{msg})
	app.DeliverMsgSimError(t, &msgExec, grantee.PrivKey, "NFT ID mismatch")

	msg = types.NewMsgUpdateOffer(granter.Addr.String(), ungrantedClassId, grantedNftId, 2, time.Unix(1234569999, 0))
	msgExec = authz.NewMsgExec(grantee.Addr, []sdk.Msg{msg})
	app.DeliverMsgSimError(t, &msgExec, grantee.PrivKey, "class ID mismatch")

	expiration = time.Unix(1300000000, 0)
	msgGrant, err = authz.NewMsgGrant(granter.Addr, grantee.Addr, &types.DeleteOfferAuthorization{
		ClassId: grantedClassId,
		NftId:   grantedNftId,
	}, &expiration)
	require.NoError(t, err)
	app.DeliverMsgNoError(t, msgGrant, granter.PrivKey)

	msg = types.NewMsgDeleteOffer(granter.Addr.String(), grantedClassId, ungrantedNftId)
	msgExec = authz.NewMsgExec(grantee.Addr, []sdk.Msg{msg})
	app.DeliverMsgSimError(t, &msgExec, grantee.PrivKey, "NFT ID mismatch")

	msg = types.NewMsgDeleteOffer(granter.Addr.String(), ungrantedClassId, grantedNftId)
	msgExec = authz.NewMsgExec(grantee.Addr, []sdk.Msg{msg})
	app.DeliverMsgSimError(t, &msgExec, grantee.PrivKey, "class ID mismatch")

	msg = types.NewMsgDeleteOffer(granter.Addr.String(), grantedClassId, grantedNftId)
	msgExec = authz.NewMsgExec(grantee.Addr, []sdk.Msg{msg})
	app.DeliverMsgNoError(t, &msgExec, grantee.PrivKey)
}

func TestOfferAuthorizationsForAllNftIds(t *testing.T) {
	var msg sdk.Msg
	var msgGrant *authz.MsgGrant
	var msgExec authz.MsgExec
	var err error

	setup := setupAppAndNfts(t)
	app := setup.App

	granter := setup.Owners[0]
	grantee := setup.OtherAddrs[0]

	grantedClassId := setup.Owners[0].Iscns[0].Classes[0].ClassId
	ungrantedClassId := setup.Owners[0].Iscns[0].Classes[1].ClassId
	nftId0 := setup.Owners[0].Iscns[0].Classes[0].NftIds[0]
	nftId1 := setup.Owners[0].Iscns[0].Classes[0].NftIds[1]

	expiration := time.Unix(1300000000, 0)
	msgGrant, err = authz.NewMsgGrant(granter.Addr, grantee.Addr, &types.CreateOfferAuthorization{
		ClassId: grantedClassId,
	}, &expiration)
	require.NoError(t, err)
	app.DeliverMsgNoError(t, msgGrant, granter.PrivKey)

	msg = types.NewMsgCreateOffer(granter.Addr.String(), grantedClassId, nftId0, 1, time.Unix(1234569999, 0))
	msgExec = authz.NewMsgExec(grantee.Addr, []sdk.Msg{msg})
	app.DeliverMsgNoError(t, &msgExec, grantee.PrivKey)

	msg = types.NewMsgCreateOffer(granter.Addr.String(), grantedClassId, nftId1, 1, time.Unix(1234569999, 0))
	msgExec = authz.NewMsgExec(grantee.Addr, []sdk.Msg{msg})
	app.DeliverMsgNoError(t, &msgExec, grantee.PrivKey)

	msg = types.NewMsgCreateOffer(granter.Addr.String(), ungrantedClassId, nftId0, 1, time.Unix(1234569999, 0))
	msgExec = authz.NewMsgExec(grantee.Addr, []sdk.Msg{msg})
	app.DeliverMsgSimError(t, &msgExec, grantee.PrivKey, "class ID mismatch")

	expiration = time.Unix(1300000000, 0)
	msgGrant, err = authz.NewMsgGrant(granter.Addr, grantee.Addr, &types.UpdateOfferAuthorization{
		ClassId: grantedClassId,
	}, &expiration)
	require.NoError(t, err)
	app.DeliverMsgNoError(t, msgGrant, granter.PrivKey)

	msg = types.NewMsgUpdateOffer(granter.Addr.String(), grantedClassId, nftId0, 2, time.Unix(1234569999, 0))
	msgExec = authz.NewMsgExec(grantee.Addr, []sdk.Msg{msg})
	app.DeliverMsgNoError(t, &msgExec, grantee.PrivKey)

	msg = types.NewMsgUpdateOffer(granter.Addr.String(), grantedClassId, nftId1, 2, time.Unix(1234569999, 0))
	msgExec = authz.NewMsgExec(grantee.Addr, []sdk.Msg{msg})
	app.DeliverMsgNoError(t, &msgExec, grantee.PrivKey)

	msg = types.NewMsgUpdateOffer(granter.Addr.String(), ungrantedClassId, nftId0, 2, time.Unix(1234569999, 0))
	msgExec = authz.NewMsgExec(grantee.Addr, []sdk.Msg{msg})
	app.DeliverMsgSimError(t, &msgExec, grantee.PrivKey, "class ID mismatch")

	expiration = time.Unix(1300000000, 0)
	msgGrant, err = authz.NewMsgGrant(granter.Addr, grantee.Addr, &types.DeleteOfferAuthorization{
		ClassId: grantedClassId,
	}, &expiration)
	require.NoError(t, err)
	app.DeliverMsgNoError(t, msgGrant, granter.PrivKey)

	msg = types.NewMsgDeleteOffer(granter.Addr.String(), ungrantedClassId, nftId0)
	msgExec = authz.NewMsgExec(grantee.Addr, []sdk.Msg{msg})
	app.DeliverMsgSimError(t, &msgExec, grantee.PrivKey, "class ID mismatch")

	msg = types.NewMsgDeleteOffer(granter.Addr.String(), grantedClassId, nftId0)
	msgExec = authz.NewMsgExec(grantee.Addr, []sdk.Msg{msg})
	app.DeliverMsgNoError(t, &msgExec, grantee.PrivKey)

	msg = types.NewMsgDeleteOffer(granter.Addr.String(), grantedClassId, nftId1)
	msgExec = authz.NewMsgExec(grantee.Addr, []sdk.Msg{msg})
	app.DeliverMsgNoError(t, &msgExec, grantee.PrivKey)
}
