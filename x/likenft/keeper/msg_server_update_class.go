package keeper

import (
	"context"

	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/likecoin/likechain/backport/cosmos-sdk/v0.46.0-alpha2/x/nft"
	"github.com/likecoin/likechain/x/likenft/types"
)

func (k msgServer) UpdateClass(goCtx context.Context, msg *types.MsgUpdateClass) (*types.MsgUpdateClassResponse, error) {
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

	// Update class
	classData = types.ClassData{
		Metadata: msg.Metadata,
		Parent:   parent.ClassParent,
		Config: types.ClassConfig{
			Burnable:  msg.Burnable,
			MaxSupply: msg.MaxSupply,
		},
	}
	classDataInAny, err := cdctypes.NewAnyWithValue(&classData)
	if err != nil {
		return nil, types.ErrFailedToMarshalData.Wrapf("%s", err.Error())
	}
	class = nft.Class{
		Id:          class.Id,
		Name:        msg.Name,
		Symbol:      msg.Symbol,
		Description: msg.Description,
		Uri:         msg.Uri,
		UriHash:     msg.UriHash,
		Data:        classDataInAny,
	}
	if err := k.nftKeeper.UpdateClass(ctx, class); err != nil {
		return nil, types.ErrFailedToUpdateClass.Wrapf("%s", err.Error())
	}

	// Emit event
	ctx.EventManager().EmitTypedEvent(&types.EventUpdateClass{
		ClassId:            class.Id,
		ParentIscnIdPrefix: classData.Parent.IscnIdPrefix,
		ParentAccount:      classData.Parent.Account,
	})

	return &types.MsgUpdateClassResponse{
		Class: class,
	}, nil
}
