package likenft_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/authz"

	"github.com/likecoin/likecoin-chain/v4/x/likenft/types"
)

func TestRoyaltyConfigAuthorizations(t *testing.T) {
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

	baseRoyaltyConfigInput := types.RoyaltyConfigInput{
		RateBasisPoints: 100,
		Stakeholders: []types.RoyaltyStakeholderInput{
			{Account: addr1.String(), Weight: 80},
			{Account: addr2.String(), Weight: 20},
		},
	}

	expiration := time.Unix(1300000000, 0)
	msgGrant, err = authz.NewMsgGrant(granter.Addr, grantee.Addr, &types.CreateRoyaltyConfigAuthorization{
		ClassId: grantedClassId,
	}, &expiration)
	require.NoError(t, err)
	app.DeliverMsgNoError(t, msgGrant, granter.PrivKey)

	msg = types.NewMsgCreateRoyaltyConfig(granter.Addr.String(), grantedClassId, baseRoyaltyConfigInput)
	msgExec = authz.NewMsgExec(grantee.Addr, []sdk.Msg{msg})
	app.DeliverMsgNoError(t, &msgExec, grantee.PrivKey)

	msg = types.NewMsgCreateRoyaltyConfig(granter.Addr.String(), ungrantedClassId, baseRoyaltyConfigInput)
	msgExec = authz.NewMsgExec(grantee.Addr, []sdk.Msg{msg})
	app.DeliverMsgSimError(t, &msgExec, grantee.PrivKey, "class ID mismatch")

	expiration = time.Unix(1300000000, 0)
	msgGrant, err = authz.NewMsgGrant(granter.Addr, grantee.Addr, &types.UpdateRoyaltyConfigAuthorization{
		ClassId: grantedClassId,
	}, &expiration)
	require.NoError(t, err)
	app.DeliverMsgNoError(t, msgGrant, granter.PrivKey)

	updatedRoyaltyConfigInput := baseRoyaltyConfigInput
	updatedRoyaltyConfigInput.RateBasisPoints = 200
	msg = types.NewMsgUpdateRoyaltyConfig(granter.Addr.String(), grantedClassId, updatedRoyaltyConfigInput)
	msgExec = authz.NewMsgExec(grantee.Addr, []sdk.Msg{msg})
	app.DeliverMsgNoError(t, &msgExec, grantee.PrivKey)

	msg = types.NewMsgUpdateRoyaltyConfig(granter.Addr.String(), ungrantedClassId, updatedRoyaltyConfigInput)
	msgExec = authz.NewMsgExec(grantee.Addr, []sdk.Msg{msg})
	app.DeliverMsgSimError(t, &msgExec, grantee.PrivKey, "class ID mismatch")

	expiration = time.Unix(1300000000, 0)
	msgGrant, err = authz.NewMsgGrant(granter.Addr, grantee.Addr, &types.DeleteRoyaltyConfigAuthorization{
		ClassId: grantedClassId,
	}, &expiration)
	require.NoError(t, err)
	app.DeliverMsgNoError(t, msgGrant, granter.PrivKey)

	msg = types.NewMsgDeleteRoyaltyConfig(granter.Addr.String(), ungrantedClassId)
	msgExec = authz.NewMsgExec(grantee.Addr, []sdk.Msg{msg})
	app.DeliverMsgSimError(t, &msgExec, grantee.PrivKey, "class ID mismatch")

	msg = types.NewMsgDeleteRoyaltyConfig(granter.Addr.String(), grantedClassId)
	msgExec = authz.NewMsgExec(grantee.Addr, []sdk.Msg{msg})
	app.DeliverMsgNoError(t, &msgExec, grantee.PrivKey)
}
