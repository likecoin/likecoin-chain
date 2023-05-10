package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/likecoin/likecoin-chain/v4/x/likenft/types"
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
	if err := k.assertBech32EqualsAccAddress(msg.Creator, owner); err != nil {
		return nil, err
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
