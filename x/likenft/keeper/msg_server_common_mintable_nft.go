package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/likecoin/likechain/backport/cosmos-sdk/v0.46.0-alpha2/x/nft"
	"github.com/likecoin/likechain/x/likenft/types"
)

func (k msgServer) validateReqToMutateBlindBoxContent(ctx sdk.Context, creator string, class nft.Class, classData types.ClassData, parent types.ClassParentWithOwner, willCreate bool) error {

	// Verify no tokens minted under class
	totalSupply := k.nftKeeper.GetTotalSupply(ctx, class.Id)
	if totalSupply > 0 {
		return types.ErrCannotUpdateClassWithMintedTokens.Wrap("Cannot update class with minted tokens")
	}

	// Check max supply vs existing mintable count
	if willCreate && classData.Config.MaxSupply > 0 && classData.BlindBoxState.ContentCount >= classData.Config.MaxSupply {
		return types.ErrNftNoSupply.Wrapf("NFT Class has reached its maximum supply: %d", classData.Config.MaxSupply)
	}

	// Check class parent relation is valid and current user is owner
	if err := k.assertBech32EqualsAccAddress(creator, parent.Owner); err != nil {
		return err
	}

	return nil
}
