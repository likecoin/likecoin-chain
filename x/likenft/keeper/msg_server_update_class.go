package keeper

import (
	"context"

	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/likecoin/likechain/backport/cosmos-sdk/v0.46.0-alpha2/x/nft"
	iscntypes "github.com/likecoin/likechain/x/iscn/types"
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

	// Verify iscn exists and is related
	var classData types.ClassData
	if err := k.cdc.Unmarshal(class.Data.Value, &classData); err != nil {
		return nil, types.ErrFailedToUnmarshalData.Wrapf(err.Error())
	}
	classesByISCN, found := k.GetClassesByISCN(ctx, classData.IscnIdPrefix)
	if !found {
		return nil, types.ErrNftClassNotRelatedToAnyIscn.Wrapf("NFT claims it is related to ISCN %s but no mapping is found", classData.IscnIdPrefix)
	}
	isRelated := false
	for _, validClassId := range classesByISCN.ClassIds {
		if validClassId == class.Id {
			// claimed relation is valid
			isRelated = true
			break
		}
	}
	if !isRelated {
		return nil, types.ErrNftClassNotRelatedToAnyIscn.Wrapf("NFT claims it is related to ISCN %s but no mapping is found", classData.IscnIdPrefix)
	}

	// Verify user is owner of iscn and thus the nft class
	iscnId, err := iscntypes.ParseIscnId(classData.IscnIdPrefix)
	if err != nil {
		return nil, types.ErrInvalidIscnId.Wrapf("%s", err.Error())
	}
	iscnRecord := k.iscnKeeper.GetContentIdRecord(ctx, iscnId.Prefix)
	if iscnRecord == nil {
		return nil, types.ErrIscnRecordNotFound.Wrapf("ISCN %s not found", iscnId.Prefix.String())
	}
	userAddress, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("%s", err.Error())
	}
	if !iscnRecord.OwnerAddress().Equals(userAddress) {
		return nil, sdkerrors.ErrUnauthorized.Wrapf("%s is not the owner of the ISCN %s", msg.Creator, iscnId.Prefix.String())
	}

	// Update class
	classData = types.ClassData{
		Metadata:          msg.Metadata,
		IscnIdPrefix:      iscnId.Prefix.String(),
		IscnVersionAtMint: iscnRecord.LatestVersion,
		Config: types.ClassConfig{
			Burnable: msg.Burnable,
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
		IscnIdPrefix: iscnId.Prefix.String(),
		ClassId:      class.Id,
		Owner:        iscnRecord.OwnerAddress().String(),
	})

	return &types.MsgUpdateClassResponse{
		Class: class,
	}, nil
}
