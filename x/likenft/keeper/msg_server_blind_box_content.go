package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/likecoin/likechain/backport/cosmos-sdk/v0.46.0-alpha2/x/nft"
	"github.com/likecoin/likechain/x/likenft/types"
)

func (k msgServer) validateReqToMutateBlindBoxContent(ctx sdk.Context, creator string, class nft.Class, classData types.ClassData, parent types.ClassParentWithOwner, willCreate bool) error {

	// Verify no tokens minted under class
	totalSupply := k.nftKeeper.GetTotalSupply(ctx, class.Id)
	if totalSupply > 0 {
		return types.ErrCannotUpdateClassWithMintedTokens.Wrap("Cannot update class with minted tokens")
	}

	// Check max supply vs existing blind box content count
	if willCreate && classData.Config.MaxSupply > 0 && classData.BlindBoxState.ContentCount >= classData.Config.MaxSupply {
		return types.ErrNftNoSupply.Wrapf("NFT Class has reached its maximum supply: %d", classData.Config.MaxSupply)
	}

	// Check class parent relation is valid and current user is owner
	if err := k.assertBech32EqualsAccAddress(creator, parent.Owner); err != nil {
		return err
	}

	return nil
}

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
		return nil, types.ErrBlindBoxContentNotFound
	}

	content := types.BlindBoxContent{
		ClassId: msg.ClassId,
		Id:      msg.Id,
		Input:   msg.Input,
	}

	// Deduct minting fee if new content is longer
	lengthDiff := content.Size() - oldBlindBoxContent.Size()
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
	k.SetBlindBoxContent(ctx, content)

	// Emit event
	ctx.EventManager().EmitTypedEvent(&types.EventUpdateBlindBoxContent{
		ClassId:                 msg.ClassId,
		ContentId:               msg.Id,
		ClassParentIscnIdPrefix: parent.IscnIdPrefix,
		ClassParentAccount:      parent.Account,
	})

	return &types.MsgUpdateBlindBoxContentResponse{
		BlindBoxContent: content,
	}, nil
}

func (k msgServer) DeleteBlindBoxContent(goCtx context.Context, msg *types.MsgDeleteBlindBoxContent) (*types.MsgDeleteBlindBoxContentResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	class, classData, err := k.GetClass(ctx, msg.ClassId)
	if err != nil {
		return nil, err
	}
	parent, err := k.ValidateAndRefreshClassParent(ctx, msg.ClassId, classData.Parent)
	if err != nil {
		return nil, err
	}
	if err := k.validateReqToMutateBlindBoxContent(ctx, msg.Creator, class, classData, parent, false); err != nil {
		return nil, err
	}

	// check id already exists
	if _, exists := k.GetBlindBoxContent(ctx, msg.ClassId, msg.Id); !exists {
		return nil, types.ErrBlindBoxContentNotFound
	}

	// remove record
	k.RemoveBlindBoxContent(ctx, msg.ClassId, msg.Id)

	// Emit event
	ctx.EventManager().EmitTypedEvent(&types.EventDeleteBlindBoxContent{
		ClassId:                 msg.ClassId,
		ContentId:               msg.Id,
		ClassParentIscnIdPrefix: parent.IscnIdPrefix,
		ClassParentAccount:      parent.Account,
	})

	return &types.MsgDeleteBlindBoxContentResponse{}, nil
}
