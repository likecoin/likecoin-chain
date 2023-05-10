package likenft_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/authz"

	"github.com/likecoin/likecoin-chain/v4/testutil"
	"github.com/likecoin/likecoin-chain/v4/x/likenft/types"
)

func TestClassAuthorizations(t *testing.T) {
	var msg sdk.Msg
	var msgGrant *authz.MsgGrant
	var msgExec authz.MsgExec
	var res *sdk.Result
	var err error

	setup := setupAppAndNfts(t)
	app := setup.App

	granter := setup.Owners[0]
	grantee := setup.OtherAddrs[0]

	grantedIscnIdPrefix := setup.Owners[0].Iscns[0].IscnId.Prefix.String()
	ungrantedIscnIdPrefix := setup.Owners[0].Iscns[1].IscnId.Prefix.String()

	grantedClassParentInput := types.ClassParentInput{
		Type:         types.ClassParentType_ISCN,
		IscnIdPrefix: grantedIscnIdPrefix,
	}
	ungrantedClassParentInput := types.ClassParentInput{
		Type:         types.ClassParentType_ISCN,
		IscnIdPrefix: ungrantedIscnIdPrefix,
	}
	baseClassInput := types.ClassInput{
		Name:     "test",
		Symbol:   "TEST",
		Metadata: types.JsonInput(`{}`),
	}

	expiration := time.Unix(1300000000, 0)
	msgGrant, err = authz.NewMsgGrant(granter.Addr, grantee.Addr, &types.NewClassAuthorization{
		IscnIdPrefix: grantedIscnIdPrefix,
	}, &expiration)
	require.NoError(t, err)
	app.DeliverMsgNoError(t, msgGrant, granter.PrivKey)

	msg = types.NewMsgNewClass(granter.Addr.String(), grantedClassParentInput, baseClassInput)
	msgExec = authz.NewMsgExec(grantee.Addr, []sdk.Msg{msg})
	res = app.DeliverMsgNoError(t, &msgExec, grantee.PrivKey)
	grantedClassId := testutil.GetClassIdFromResult(t, res)

	msg = types.NewMsgNewClass(granter.Addr.String(), ungrantedClassParentInput, baseClassInput)
	msgExec = authz.NewMsgExec(grantee.Addr, []sdk.Msg{msg})
	app.DeliverMsgSimError(t, &msgExec, grantee.PrivKey, "ISCN ID prefix mismatch")

	ungrantedClassId := setup.Owners[0].Iscns[0].Classes[0].ClassId

	expiration = time.Unix(1300000000, 0)
	msgGrant, err = authz.NewMsgGrant(granter.Addr, grantee.Addr, &types.UpdateClassAuthorization{
		ClassId: grantedClassId,
	}, &expiration)
	require.NoError(t, err)
	app.DeliverMsgNoError(t, msgGrant, granter.PrivKey)

	updatedClassInput := baseClassInput
	updatedClassInput.Name = "Updated"
	updatedClassInput.Symbol = "UPDATED"
	msg = types.NewMsgUpdateClass(granter.Addr.String(), grantedClassId, updatedClassInput)
	msgExec = authz.NewMsgExec(grantee.Addr, []sdk.Msg{msg})
	app.DeliverMsgNoError(t, &msgExec, grantee.PrivKey)

	msg = types.NewMsgUpdateClass(granter.Addr.String(), ungrantedClassId, updatedClassInput)
	msgExec = authz.NewMsgExec(grantee.Addr, []sdk.Msg{msg})
	app.DeliverMsgSimError(t, &msgExec, grantee.PrivKey, "class ID mismatch")
}
