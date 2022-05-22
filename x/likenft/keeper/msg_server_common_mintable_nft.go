package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/likecoin/likechain/x/likenft/types"
)

func (k msgServer) getParentOwnerAndValidateReqToMutateMintableNFT(ctx sdk.Context, creator string, classId string, willCreate bool) (*types.ClassParentAndOwner, error) {

	// Verify class exists
	class, classData, err := k.GetClass(ctx, classId)
	if err != nil {
		return nil, err
	}

	// Verify no tokens minted under class
	totalSupply := k.nftKeeper.GetTotalSupply(ctx, class.Id)
	if totalSupply > 0 {
		return nil, types.ErrCannotUpdateClassWithMintedTokens.Wrap("Cannot update class with minted tokens")
	}

	// Check max supply vs existing mintable count
	if willCreate && classData.Config.MaxSupply > 0 && classData.MintableCount >= classData.Config.MaxSupply {
		return nil, types.ErrNftNoSupply.Wrapf("NFT Class has reached its maximum supply: %d", classData.Config.MaxSupply)
	}

	// Check class parent relation is valid and current user is owner
	parentAndOwner, err := k.validateAndGetClassParentAndOwner(ctx, class.Id, &classData)
	if err != nil {
		return nil, err
	}

	userAddress, err := sdk.AccAddressFromBech32(creator)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("%s", err.Error())
	}
	if !parentAndOwner.Owner.Equals(userAddress) {
		return nil, sdkerrors.ErrUnauthorized.Wrapf("%s is not authorized", userAddress.String())
	}

	return parentAndOwner, nil
}
