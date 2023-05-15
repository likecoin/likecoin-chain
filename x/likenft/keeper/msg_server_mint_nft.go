package keeper

import (
	"context"
	"fmt"

	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/nft"
	"github.com/likecoin/likecoin-chain/v4/x/likenft/types"
)

func (k msgServer) mintBlindBoxNFT(ctx sdk.Context, classId string, classData *types.ClassData, ownerAddress sdk.AccAddress, userAddress sdk.AccAddress, totalSupply uint64, msg *types.MsgMintNFT) (*nft.NFT, error) {
	params := k.GetParams(ctx)
	tokenId := fmt.Sprintf("nft%d", totalSupply+1)

	// Check if the class has already been revealed or not
	if !classData.BlindBoxState.ToBeRevealed {
		return nil, types.ErrFailedToMintNFT.Wrapf(fmt.Sprintf("The class %s has already been revealed", classId))
	}

	// Check blind box content supply
	if totalSupply >= classData.BlindBoxState.ContentCount {
		return nil, types.ErrNftNoSupply.Wrapf("All blind box content have been minted")
	}

	// Resolve the most applicable mint period
	mintPeriod, err := k.resolveValidMintPeriod(ctx, classId, *classData.Config.BlindBoxConfig, ownerAddress, userAddress)
	if err != nil {
		return nil, err
	}

	if mintPeriod == nil {
		return nil, sdkerrors.ErrUnauthorized.Wrapf(fmt.Sprintf("The user %s is not allowed to mint the class %s", userAddress, classId))
	}

	nftData := types.NFTData{
		ClassParent:  classData.Parent,
		ToBeRevealed: true,
	}

	nftDataInAny, err := cdctypes.NewAnyWithValue(&nftData)
	if err != nil {
		return nil, types.ErrFailedToMarshalData.Wrapf("%s", err.Error())
	}
	nft := nft.NFT{
		ClassId: classId,
		Id:      tokenId,
		Data:    nftDataInAny,
	}

	// Pay price to owner if mintPrice is not zero and the minter is not the owner
	if !ownerAddress.Equals(userAddress) && mintPeriod.MintPrice > 0 {
		spentableTokens := k.bankKeeper.GetBalance(ctx, userAddress, params.GetPriceDenom())
		if spentableTokens.Amount.Uint64() < mintPeriod.MintPrice {
			return nil, types.ErrInsufficientFunds.Wrapf("insufficient funds to mint tokenId %s", tokenId)
		}

		err = k.bankKeeper.SendCoins(ctx, userAddress, ownerAddress, sdk.NewCoins(sdk.NewCoin(params.GetPriceDenom(), sdk.NewInt(int64(mintPeriod.MintPrice)))))
		if err != nil {
			return nil, types.ErrFailedToMintNFT.Wrapf("%s", err.Error())
		}
	}

	err = k.nftKeeper.Mint(ctx, nft, userAddress)
	if err != nil {
		return nil, types.ErrFailedToMintNFT.Wrapf("%s", err.Error())
	}

	return &nft, nil
}

func (k msgServer) mintRegularNFT(ctx sdk.Context, classId string, classData *types.ClassData, userAddress sdk.AccAddress, msg *types.MsgMintNFT) (*nft.NFT, error) {
	// Validate NFTInput not null
	if msg.Input == nil {
		return nil, sdkerrors.ErrInvalidRequest.Wrapf("NFTInput cannot be null")
	}

	// Validate token id
	if err := nft.ValidateNFTID(msg.Id); err != nil {
		return nil, types.ErrInvalidTokenId.Wrapf("%s", err)
	}

	nftData := types.NFTData{
		Metadata:     msg.Input.Metadata,
		ClassParent:  classData.Parent,
		ToBeRevealed: false,
	}
	nftDataInAny, err := cdctypes.NewAnyWithValue(&nftData)
	if err != nil {
		return nil, types.ErrFailedToMarshalData.Wrapf("%s", err.Error())
	}
	nft := nft.NFT{
		ClassId: classId,
		Id:      msg.Id,
		Uri:     msg.Input.Uri,
		UriHash: msg.Input.UriHash,
		Data:    nftDataInAny,
	}
	// Deduct minting fee
	err = k.DeductFeePerByte(ctx, userAddress, nft.Size(), msg)
	if err != nil {
		return nil, err
	}
	err = k.nftKeeper.Mint(ctx, nft, userAddress)
	if err != nil {
		return nil, types.ErrFailedToMintNFT.Wrapf("%s", err.Error())
	}
	return &nft, nil
}

func (k msgServer) MintNFT(goCtx context.Context, msg *types.MsgMintNFT) (*types.MsgMintNFTResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Assert class exists
	class, classData, err := k.GetClass(ctx, msg.ClassId)
	if err != nil {
		return nil, err
	}

	// validate class parent relation & resolve owner
	// also refresh parent info (e.g. iscn latest version)
	parent, err := k.ValidateAndRefreshClassParent(ctx, class.Id, classData.Parent)
	if err != nil {
		return nil, err
	}

	userAddress, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("%s", err.Error())
	}

	totalSupply := k.nftKeeper.GetTotalSupply(ctx, class.Id)

	// Refresh recorded iscn version in class if needed and first mint
	if classData.Parent.IscnVersionAtMint != parent.IscnVersionAtMint &&
		totalSupply <= 0 {
		classData.Parent = parent.ClassParent
		classDataInAny, err := cdctypes.NewAnyWithValue(&classData)
		if err != nil {
			return nil, types.ErrFailedToMarshalData.Wrapf("%s", err.Error())
		}
		class.Data = classDataInAny
		err = k.nftKeeper.UpdateClass(ctx, class)
		if err != nil {
			return nil, types.ErrFailedToUpdateClass.Wrapf("%s", err.Error())
		}
	}

	// Assert supply is enough
	if classData.Config.MaxSupply > 0 &&
		totalSupply >= classData.Config.MaxSupply {
		return nil, types.ErrNftNoSupply.Wrapf("NFT Class has reached its maximum supply: %d", classData.Config.MaxSupply)
	}

	// Mint NFT
	var nft *nft.NFT
	if classData.Config.IsBlindBox() {
		nft, err = k.mintBlindBoxNFT(ctx, class.Id, &classData, parent.Owner, userAddress, totalSupply, msg)
		if err != nil {
			return nil, err
		}
	} else if parent.Owner.Equals(userAddress) {
		nft, err = k.mintRegularNFT(ctx, class.Id, &classData, userAddress, msg)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, sdkerrors.ErrUnauthorized.Wrapf("%s is not authorized", userAddress.String())
	}

	// Emit event
	ctx.EventManager().EmitTypedEvent(&types.EventMintNFT{
		ClassId:                 nft.ClassId,
		NftId:                   nft.Id,
		Owner:                   userAddress.String(),
		ClassParentIscnIdPrefix: classData.Parent.IscnIdPrefix,
		ClassParentAccount:      classData.Parent.Account,
	})

	return &types.MsgMintNFTResponse{
		Nft: *nft,
	}, nil
}
