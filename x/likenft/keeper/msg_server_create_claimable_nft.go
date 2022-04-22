package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/likecoin/likechain/x/likenft/types"
)

func (k msgServer) CreateClaimableNFT(goCtx context.Context, msg *types.MsgCreateClaimableNFT) (*types.MsgCreateClaimableNFTResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Verify class exists
	class, found := k.nftKeeper.GetClass(ctx, msg.ClassId)
	if !found {
		return nil, types.ErrNftClassNotFound.Wrapf("Class id %s not found", msg.ClassId)
	}

	// Verify no tokens minted under class
	totalSupply := k.nftKeeper.GetTotalSupply(ctx, class.Id)
	if totalSupply > 0 {
		return nil, types.ErrCannotUpdateClassWithMintedTokens.Wrap("Cannot update class with minted tokens")
	}

	// Check class parent relation is valid and current user is owner
	var classData types.ClassData
	if err := k.cdc.Unmarshal(class.Data.Value, &classData); err != nil {
		return nil, types.ErrFailedToUnmarshalData.Wrapf(err.Error())
	}

	parentAndOwner, err := k.validateAndGetClassParentAndOwner(ctx, class.Id, &classData)
	if err != nil {
		return nil, err
	}

	userAddress, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("%s", err.Error())
	}
	if !parentAndOwner.Owner.Equals(userAddress) {
		return nil, sdkerrors.ErrUnauthorized.Wrapf("%s is not authorized", userAddress.String())
	}

	// check id not already exist
	if _, exists := k.GetClaimableNFT(ctx, class.Id, msg.Id); exists {
		return nil, types.ErrClaimableNftAlreadyExists
	}

	// set record
	k.SetClaimableNFT(ctx, types.ClaimableNFT{
		ClassId: msg.ClassId,
		Id:      msg.Id,
		Input:   msg.Input,
	})

	// TODO emit event
	return &types.MsgCreateClaimableNFTResponse{}, nil
}
