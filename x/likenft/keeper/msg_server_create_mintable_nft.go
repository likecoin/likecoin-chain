package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/likecoin/likechain/x/likenft/types"
)

func (k msgServer) CreateBlindBoxContent(goCtx context.Context, msg *types.MsgCreateBlindBoxContent) (*types.MsgCreateBlindBoxContentResponse, error) {
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

	// check id not already exist
	if _, exists := k.GetBlindBoxContent(ctx, msg.ClassId, msg.Id); exists {
		return nil, types.ErrBlindBoxContentAlreadyExists
	}

	content := types.BlindBoxContent{
		ClassId: msg.ClassId,
		Id:      msg.Id,
		Input:   msg.Input,
	}

	// Deduct minting fee
	userAddress, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf(err.Error())
	}
	err = k.DeductFeePerByte(ctx, userAddress, content.Size())
	if err != nil {
		return nil, err
	}

	// set record
	k.SetBlindBoxContent(ctx, content)

	// Emit event
	ctx.EventManager().EmitTypedEvent(&types.EventCreateBlindBoxContent{
		ClassId:                 msg.ClassId,
		ContentId:               msg.Id,
		ClassParentIscnIdPrefix: parent.IscnIdPrefix,
		ClassParentAccount:      parent.Account,
	})

	return &types.MsgCreateBlindBoxContentResponse{
		BlindBoxContent: content,
	}, nil
}
