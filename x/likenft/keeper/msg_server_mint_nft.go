package keeper

import (
	"context"

	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/likecoin/likechain/backport/cosmos-sdk/v0.46.0-alpha2/x/nft"
	"github.com/likecoin/likechain/x/likenft/types"
)

func (k msgServer) MintNFT(goCtx context.Context, msg *types.MsgMintNFT) (*types.MsgMintNFTResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Validate token id
	if err := nft.ValidateNFTID(msg.Id); err != nil {
		return nil, types.ErrInvalidTokenId.Wrapf("%s", err)
	}

	// Assert class exists
	class, found := k.nftKeeper.GetClass(ctx, msg.ClassId)
	if !found {
		return nil, types.ErrNftClassNotFound.Wrapf("Class id %s not found", msg.ClassId)
	}

	// Assert class has related iscn
	var classData types.ClassData
	if err := k.cdc.Unmarshal(class.Data.Value, &classData); err != nil {
		return nil, types.ErrFailedToUnmarshalData.Wrapf(err.Error())
	}
	if err := k.validateClassParentRelation(ctx, class.Id, classData.Parent); err != nil {
		return nil, err
	}

	// refresh parent info (e.g. iscn latest version) & check ownership
	parent, err := k.resolveClassParentAndOwner(ctx, classData.Parent.ToInput(), classData.Parent.Account)
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

	// Refresh recorded iscn version in class if needed and first mint
	if classData.Parent.IscnVersionAtMint != parent.IscnVersionAtMint &&
		k.nftKeeper.GetTotalSupply(ctx, class.Id) <= 0 {
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

	// Mint NFT
	nftData := types.NFTData{
		Metadata:    msg.Metadata,
		ClassParent: classData.Parent,
	}
	nftDataInAny, err := cdctypes.NewAnyWithValue(&nftData)
	if err != nil {
		return nil, types.ErrFailedToMarshalData.Wrapf("%s", err.Error())
	}
	nft := nft.NFT{
		ClassId: class.Id,
		Id:      msg.Id,
		Uri:     msg.Uri,
		UriHash: msg.UriHash,
		Data:    nftDataInAny,
	}
	err = k.nftKeeper.Mint(ctx, nft, userAddress)
	if err != nil {
		return nil, types.ErrFailedToMintNFT.Wrapf("%s", err.Error())
	}

	// Emit event
	ctx.EventManager().EmitTypedEvent(&types.EventMintNFT{
		ClassId:                 class.Id,
		NftId:                   msg.Id,
		Owner:                   userAddress.String(),
		ClassParentIscnIdPrefix: classData.Parent.IscnIdPrefix,
		ClassParentAccount:      classData.Parent.Account,
	})

	return &types.MsgMintNFTResponse{
		Nft: nft,
	}, nil
}
