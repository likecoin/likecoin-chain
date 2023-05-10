package likenft_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/authz"

	"github.com/cosmos/cosmos-sdk/x/nft"
	"github.com/likecoin/likecoin-chain/v4/x/likenft/types"
)

func TestNFTAuthorizations(t *testing.T) {
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

	grantedNftId := "this-is-granted"
	ungrantedNftId := setup.Owners[0].Iscns[0].Classes[0].NftIds[0]

	expiration := time.Unix(1300000000, 0)
	msgGrant, err = authz.NewMsgGrant(granter.Addr, grantee.Addr, &types.MintNFTAuthorization{
		ClassId: grantedClassId,
	}, &expiration)
	require.NoError(t, err)
	app.DeliverMsgNoError(t, msgGrant, granter.PrivKey)

	msg = types.NewMsgMintNFT(granter.Addr.String(), grantedClassId, grantedNftId, &types.NFTInput{})
	msgExec = authz.NewMsgExec(grantee.Addr, []sdk.Msg{msg})
	app.DeliverMsgNoError(t, &msgExec, grantee.PrivKey)

	msg = types.NewMsgMintNFT(granter.Addr.String(), ungrantedClassId, grantedNftId, &types.NFTInput{})
	msgExec = authz.NewMsgExec(grantee.Addr, []sdk.Msg{msg})
	app.DeliverMsgSimError(t, &msgExec, grantee.PrivKey, "class ID mismatch")

	expiration = time.Unix(1300000000, 0)
	msgGrant, err = authz.NewMsgGrant(granter.Addr, grantee.Addr, &types.SendNFTAuthorization{
		ClassId: grantedClassId,
		Id:      grantedNftId,
	}, &expiration)
	require.NoError(t, err)
	app.DeliverMsgNoError(t, msgGrant, granter.PrivKey)

	msg = &nft.MsgSend{
		ClassId:  grantedClassId,
		Id:       grantedNftId,
		Sender:   granter.Addr.String(),
		Receiver: grantee.Addr.String(),
	}
	msgExec = authz.NewMsgExec(grantee.Addr, []sdk.Msg{msg})
	app.DeliverMsgNoError(t, &msgExec, grantee.PrivKey)

	msg = &nft.MsgSend{
		ClassId:  grantedClassId,
		Id:       ungrantedNftId,
		Sender:   granter.Addr.String(),
		Receiver: grantee.Addr.String(),
	}
	msgExec = authz.NewMsgExec(grantee.Addr, []sdk.Msg{msg})
	app.DeliverMsgSimError(t, &msgExec, grantee.PrivKey, "NFT ID mismatch")

	expiration = time.Unix(1300000000, 0)
	msgGrant, err = authz.NewMsgGrant(granter.Addr, grantee.Addr, &types.SendNFTAuthorization{
		ClassId: grantedClassId,
	}, &expiration)
	require.NoError(t, err)
	app.DeliverMsgNoError(t, msgGrant, granter.PrivKey)

	msg = &nft.MsgSend{
		ClassId:  grantedClassId,
		Id:       ungrantedNftId,
		Sender:   granter.Addr.String(),
		Receiver: grantee.Addr.String(),
	}
	msgExec = authz.NewMsgExec(grantee.Addr, []sdk.Msg{msg})
	app.DeliverMsgNoError(t, &msgExec, grantee.PrivKey)

	msg = &nft.MsgSend{
		ClassId:  ungrantedClassId,
		Id:       ungrantedNftId,
		Sender:   granter.Addr.String(),
		Receiver: grantee.Addr.String(),
	}
	msgExec = authz.NewMsgExec(grantee.Addr, []sdk.Msg{msg})
	app.DeliverMsgSimError(t, &msgExec, grantee.PrivKey, "class ID mismatch")
}
