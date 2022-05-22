package keeper

import (
	"fmt"
	"math/rand"

	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/likecoin/likechain/backport/cosmos-sdk/v0.46.0-alpha2/x/nft"
	"github.com/likecoin/likechain/x/likenft/types"
	"github.com/likecoin/likechain/x/likenft/utils"
)

func (k Keeper) RevealMintableNFTs(ctx sdk.Context, classId string) error {
	// check if class is using blindbox
	class, classData, err := k.GetClass(ctx, classId)
	if err != nil {
		return err
	}
	if !classData.Config.IsBlindBox() {
		return types.ErrClassIsNotBlindBox
	}
	// validate class parent relation and resolve owner
	parentAndOwner, err := k.ValidateAndRefreshClassParent(ctx, classId, classData.Parent)
	if err != nil {
		return err
	}
	// mint all remaining supply to owner
	totalSupply := k.nftKeeper.GetTotalSupply(ctx, classId)
	remainingSupply := classData.MintableCount - totalSupply
	for i := 0; i < int(remainingSupply); i++ {
		tokenId := fmt.Sprintf("nft%d", int(totalSupply)+i+1)

		nftData := types.NFTData{
			ClassParent:  classData.Parent,
			ToBeRevealed: true,
		}

		nftDataInAny, err := cdctypes.NewAnyWithValue(&nftData)
		if err != nil {
			return types.ErrFailedToMarshalData.Wrapf("%s", err.Error())
		}
		nft := nft.NFT{
			ClassId: classId,
			Id:      tokenId,
			Data:    nftDataInAny,
		}
		err = k.nftKeeper.Mint(ctx, nft, parentAndOwner.Owner)
		if err != nil {
			return types.ErrFailedToMintNFT.Wrapf("%s", err.Error())
		}
	}

	// get list of mintable ids and shuffle
	var mintableIDs []string
	k.IterateMintableNFTs(ctx, classId, func(val types.MintableNFT) {
		mintableIDs = append(mintableIDs, val.Id)
	})

	// shuffle with last block hash as seed
	rand.Seed(utils.RandSeedFromLastBlock(ctx))
	rand.Shuffle(len(mintableIDs), func(i, j int) {
		mintableIDs[i], mintableIDs[j] = mintableIDs[j], mintableIDs[i]
	})

	// reveal tokens
	tokens := k.nftKeeper.GetNFTsOfClass(ctx, classId)
	if len(tokens) != len(mintableIDs) {
		// should not happen
		return fmt.Errorf("mintable length %d and minted tokens %d length mismatch", len(mintableIDs), len(tokens))
	}
	for i, token := range tokens {
		// get assigned data
		assigned, found := k.GetMintableNFT(ctx, classId, mintableIDs[i])
		if !found {
			return types.ErrMintableNftNotFound
		}
		// write data to token
		var nftData types.NFTData
		if err = nftData.Unmarshal(token.Data.Value); err != nil {
			return types.ErrFailedToUnmarshalData.Wrapf("%s", err.Error())
		}
		nftData.Metadata = assigned.Input.Metadata
		nftData.ToBeRevealed = false
		nftDataInAny, err := cdctypes.NewAnyWithValue(&nftData)
		if err != nil {
			return types.ErrFailedToMarshalData.Wrapf("%s", err.Error())
		}
		token.Uri = assigned.Input.Uri
		token.UriHash = assigned.Input.UriHash
		token.Data = nftDataInAny
		if err := k.nftKeeper.Update(ctx, token); err != nil {
			return types.ErrFailedToUpdateNFT.Wrapf("%s", err.Error())
		}
	}

	// Update revealed flag on class
	classData.ToBeRevealed = false
	classDataInAny, err := cdctypes.NewAnyWithValue(&classData)
	if err != nil {
		return types.ErrFailedToMarshalData.Wrapf("%s", err.Error())
	}
	class.Data = classDataInAny
	if err := k.nftKeeper.UpdateClass(ctx, class); err != nil {
		return types.ErrFailedToUpdateClass.Wrapf("%s", err.Error())
	}

	// Delete all mintables
	k.RemoveMintableNFTs(ctx, classId)

	return nil
}
