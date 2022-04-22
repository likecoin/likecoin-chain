package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/likecoin/likechain/x/likenft/types"
)

func (k msgServer) validateRequestToMutateClaimableNFT(ctx sdk.Context, creator string, classId string) error {

	// Verify class exists
	class, found := k.nftKeeper.GetClass(ctx, classId)
	if !found {
		return types.ErrNftClassNotFound.Wrapf("Class id %s not found", classId)
	}

	// Verify no tokens minted under class
	totalSupply := k.nftKeeper.GetTotalSupply(ctx, class.Id)
	if totalSupply > 0 {
		return types.ErrCannotUpdateClassWithMintedTokens.Wrap("Cannot update class with minted tokens")
	}

	// Check class parent relation is valid and current user is owner
	var classData types.ClassData
	if err := k.cdc.Unmarshal(class.Data.Value, &classData); err != nil {
		return types.ErrFailedToUnmarshalData.Wrapf(err.Error())
	}

	parentAndOwner, err := k.validateAndGetClassParentAndOwner(ctx, class.Id, &classData)
	if err != nil {
		return err
	}

	userAddress, err := sdk.AccAddressFromBech32(creator)
	if err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("%s", err.Error())
	}
	if !parentAndOwner.Owner.Equals(userAddress) {
		return sdkerrors.ErrUnauthorized.Wrapf("%s is not authorized", userAddress.String())
	}

	return nil
}
