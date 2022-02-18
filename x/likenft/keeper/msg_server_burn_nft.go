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
	exists := k.nftKeeper.HasNFT(ctx, msg.ClassID, msg.NftID)
	if !exists {
		return nil, types.ErrNftNotFound.Wrapf("Class %s NFT %s does not exist", msg.ClassID, msg.NftID)
	}

	// Check user is owner
	owner := k.nftKeeper.GetOwner(ctx, msg.ClassID, msg.NftID)
	user, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("%s", err.Error())
	}
	if !owner.Equals(user) {
		return nil, sdkerrors.ErrUnauthorized.Wrapf("User %s is not owner of the NFT", msg.Creator)
	}

	// Check class is set to burnable
	class, found := k.nftKeeper.GetClass(ctx, msg.ClassID)
	if !found {
		return nil, types.ErrNftClassNotFound.Wrapf("NFT Class %s not found", msg.ClassID)
	}
	var classData types.ClassData
	if err := k.cdc.Unmarshal(class.Data.Value, &classData); err != nil {
		return nil, types.ErrFailedToUnmarshalData.Wrapf(err.Error())
	}
	if !classData.Config.Burnable {
		return nil, types.ErrNftNotBurnable.Wrapf("NFT of class %s is not burnable", class.Id)
	}

	// Burn NFT
	err = k.nftKeeper.Burn(ctx, msg.ClassID, msg.NftID)
	if err != nil {
		return nil, types.ErrFailedToBurnNFT.Wrapf("%s", err.Error())
	}

	return &types.MsgBurnNFTResponse{}, nil
}
