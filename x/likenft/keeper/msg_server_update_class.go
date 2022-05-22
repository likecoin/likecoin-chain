package keeper

import (
	"context"

	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/likecoin/likechain/backport/cosmos-sdk/v0.46.0-alpha2/x/nft"
	"github.com/likecoin/likechain/x/likenft/types"
)

func (k msgServer) UpdateClass(goCtx context.Context, msg *types.MsgUpdateClass) (*types.MsgUpdateClassResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Verify class exists
	class, classData, err := k.GetClass(ctx, msg.ClassId)
	if err != nil {
		return nil, err
	}

	// Verify no tokens minted under class
	totalSupply := k.nftKeeper.GetTotalSupply(ctx, class.Id)
	if totalSupply > 0 {
		return nil, types.ErrCannotUpdateClassWithMintedTokens.Wrap("Cannot update class with minted tokens")
	}

	// Verify and Cleanup class config
	cleanClassConfig, err := k.sanitizeClassConfig(msg.Input.Config, classData.MintableCount)
	if cleanClassConfig == nil || err != nil {
		return nil, err
	}
	msg.Input.Config = *cleanClassConfig

	// Check class parent relation is valid and current user is owner
	// also refresh parent info (e.g. iscn latest version)
	parent, err := k.ValidateAndRefreshClassParent(ctx, class.Id, classData.Parent)
	if err != nil {
		return nil, err
	}

	userAddress, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("%s", err.Error())
	}
	if !parent.Owner.Equals(userAddress) {
		return nil, sdkerrors.ErrUnauthorized.Wrapf("%s is not authorized", userAddress.String())
	}

	originalConfig := classData.Config
	updatedConfig := msg.Input.Config

	// Update class
	classData.Metadata = msg.Input.Metadata
	classData.Parent = parent.ClassParent
	classData.Config = msg.Input.Config
	classData.ToBeRevealed = msg.Input.Config.IsBlindBox()
	classDataInAny, err := cdctypes.NewAnyWithValue(&classData)
	if err != nil {
		return nil, types.ErrFailedToMarshalData.Wrapf("%s", err.Error())
	}
	class = nft.Class{
		Id:          class.Id,
		Name:        msg.Input.Name,
		Symbol:      msg.Input.Symbol,
		Description: msg.Input.Description,
		Uri:         msg.Input.Uri,
		UriHash:     msg.Input.UriHash,
		Data:        classDataInAny,
	}
	if err := k.nftKeeper.UpdateClass(ctx, class); err != nil {
		return nil, types.ErrFailedToUpdateClass.Wrapf("%s", err.Error())
	}

	// Dequeue original reveal schedule
	if originalConfig.IsBlindBox() {
		k.RemoveClassRevealQueueEntry(ctx, originalConfig.BlindBoxConfig.RevealTime, class.Id)
	}

	// Enqueue new reveal schedule
	if updatedConfig.IsBlindBox() {
		k.SetClassRevealQueueEntry(ctx, types.ClassRevealQueueEntry{
			ClassId:    class.Id,
			RevealTime: updatedConfig.BlindBoxConfig.RevealTime,
		})
	}

	// Emit event
	ctx.EventManager().EmitTypedEvent(&types.EventUpdateClass{
		ClassId:            class.Id,
		ParentIscnIdPrefix: classData.Parent.IscnIdPrefix,
		ParentAccount:      classData.Parent.Account,
	})

	return &types.MsgUpdateClassResponse{
		Class: class,
	}, nil
}
