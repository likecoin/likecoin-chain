package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/likecoin/likechain/x/likenft/types"
)

func (k msgServer) UpdateBlindBoxContent(goCtx context.Context, msg *types.MsgUpdateBlindBoxContent) (*types.MsgUpdateBlindBoxContentResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	class, classData, err := k.GetClass(ctx, msg.ClassId)
	if err != nil {
		return nil, err
	}
	parent, err := k.ValidateAndRefreshClassParent(ctx, msg.ClassId, classData.Parent)
	if err != nil {
		return nil, err
	}
	if err := k.validateReqToMutateBlindBoxContent(ctx, msg.Creator, class, classData, parent, true); err != nil {
		return nil, err
	}

	// check id already exists
	oldBlindBoxContent, exists := k.GetBlindBoxContent(ctx, msg.ClassId, msg.Id)
	if !exists {
		return nil, types.ErrMintableNftNotFound
	}

	mintableNFT := types.BlindBoxContent{
		ClassId: msg.ClassId,
		Id:      msg.Id,
		Input:   msg.Input,
	}

	// Deduct minting fee if new content is longer
	lengthDiff := mintableNFT.Size() - oldBlindBoxContent.Size()
	if lengthDiff > 0 {
		userAddress, err := sdk.AccAddressFromBech32(msg.Creator)
		if err != nil {
			return nil, sdkerrors.ErrInvalidAddress.Wrapf(err.Error())
		}
		err = k.DeductFeePerByte(ctx, userAddress, lengthDiff)
		if err != nil {
			return nil, err
		}
	}

	// set record
	k.SetBlindBoxContent(ctx, mintableNFT)

	// Emit event
	ctx.EventManager().EmitTypedEvent(&types.EventUpdateBlindBoxContent{
		ClassId:                 msg.ClassId,
		ContentId:               msg.Id,
		ClassParentIscnIdPrefix: parent.IscnIdPrefix,
		ClassParentAccount:      parent.Account,
	})

	return &types.MsgUpdateBlindBoxContentResponse{
		BlindBoxContent: mintableNFT,
	}, nil
}
