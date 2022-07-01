package keeper

import (
	"context"

	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/likecoin/likecoin-chain/v3/backport/cosmos-sdk/v0.46.0-alpha2/x/nft"
	"github.com/likecoin/likecoin-chain/v3/x/likenft/types"
)

func (k msgServer) UpdateClass(goCtx context.Context, msg *types.MsgUpdateClass) (*types.MsgUpdateClassResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Verify class exists
	oldClass, classData, err := k.GetClass(ctx, msg.ClassId)
	if err != nil {
		return nil, err
	}

	// Verify no tokens minted under class
	totalSupply := k.nftKeeper.GetTotalSupply(ctx, oldClass.Id)
	if totalSupply > 0 {
		return nil, types.ErrCannotUpdateClassWithMintedTokens.Wrap("Cannot update class with minted tokens")
	}

	// Verify and Cleanup class config
	cleanClassConfig, err := k.sanitizeClassConfig(ctx, msg.Input.Config, classData.BlindBoxState.ContentCount)
	if cleanClassConfig == nil || err != nil {
		return nil, err
	}
	msg.Input.Config = *cleanClassConfig

	// Check class parent relation is valid and current user is owner
	// also refresh parent info (e.g. iscn latest version)
	parent, err := k.ValidateAndRefreshClassParent(ctx, oldClass.Id, classData.Parent)
	if err != nil {
		return nil, err
	}

	if err := k.assertBech32EqualsAccAddress(msg.Creator, parent.Owner); err != nil {
		return nil, err
	}

	originalConfig := classData.Config
	updatedConfig := msg.Input.Config

	// Update class
	classData.Metadata = msg.Input.Metadata
	classData.Parent = parent.ClassParent
	classData.Config = msg.Input.Config
	classData.BlindBoxState.ToBeRevealed = msg.Input.Config.IsBlindBox()
	classDataInAny, err := cdctypes.NewAnyWithValue(&classData)
	if err != nil {
		return nil, types.ErrFailedToMarshalData.Wrapf("%s", err.Error())
	}
	newClass := nft.Class{
		Id:          oldClass.Id,
		Name:        msg.Input.Name,
		Symbol:      msg.Input.Symbol,
		Description: msg.Input.Description,
		Uri:         msg.Input.Uri,
		UriHash:     msg.Input.UriHash,
		Data:        classDataInAny,
	}
	// Deduct minting fee if new content is longer
	lengthDiff := newClass.Size() - oldClass.Size()
	if lengthDiff > 0 {
		err = k.DeductFeePerByte(ctx, parent.Owner, lengthDiff)
		if err != nil {
			return nil, err
		}
	}
	if err := k.nftKeeper.UpdateClass(ctx, newClass); err != nil {
		return nil, types.ErrFailedToUpdateClass.Wrapf("%s", err.Error())
	}

	// Dequeue original reveal schedule
	if originalConfig.IsBlindBox() {
		k.RemoveClassRevealQueueEntry(ctx, originalConfig.BlindBoxConfig.RevealTime, newClass.Id)
	}

	// Enqueue new reveal schedule
	if updatedConfig.IsBlindBox() {
		k.SetClassRevealQueueEntry(ctx, types.ClassRevealQueueEntry{
			ClassId:    newClass.Id,
			RevealTime: updatedConfig.BlindBoxConfig.RevealTime,
		})
	}

	// Emit event
	ctx.EventManager().EmitTypedEvent(&types.EventUpdateClass{
		ClassId:            newClass.Id,
		ParentIscnIdPrefix: classData.Parent.IscnIdPrefix,
		ParentAccount:      classData.Parent.Account,
	})

	return &types.MsgUpdateClassResponse{
		Class: newClass,
	}, nil
}
