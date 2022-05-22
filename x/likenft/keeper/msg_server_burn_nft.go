package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/likecoin/likechain/x/likenft/types"
)

func (k msgServer) BurnNFT(goCtx context.Context, msg *types.MsgBurnNFT) (*types.MsgBurnNFTResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check nft exists
	exists := k.nftKeeper.HasNFT(ctx, msg.ClassId, msg.NftId)
	if !exists {
		return nil, types.ErrNftNotFound.Wrapf("Class %s NFT %s does not exist", msg.ClassId, msg.NftId)
	}

	// Check user is owner
	owner := k.nftKeeper.GetOwner(ctx, msg.ClassId, msg.NftId)
	user, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("%s", err.Error())
	}
	if !owner.Equals(user) {
		return nil, sdkerrors.ErrUnauthorized.Wrapf("User %s is not owner of the NFT", msg.Creator)
	}

	// Check class is set to burnable
	class, classData, err := k.GetClass(ctx, msg.ClassId)
	if err != nil {
		return nil, err
	}
	if !classData.Config.Burnable {
		return nil, types.ErrNftNotBurnable.Wrapf("NFT of class %s is not burnable", class.Id)
	}

	// Burn NFT
	err = k.nftKeeper.Burn(ctx, msg.ClassId, msg.NftId)
	if err != nil {
		return nil, types.ErrFailedToBurnNFT.Wrapf("%s", err.Error())
	}

	// Emit event
	ctx.EventManager().EmitTypedEvent(&types.EventBurnNFT{
		ClassId:                 class.Id,
		NftId:                   msg.NftId,
		Owner:                   owner.String(),
		ClassParentIscnIdPrefix: classData.Parent.IscnIdPrefix,
		ClassParentAccount:      classData.Parent.Account,
	})

	return &types.MsgBurnNFTResponse{}, nil
}
