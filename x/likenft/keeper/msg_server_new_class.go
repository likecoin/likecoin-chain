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

func (k msgServer) NewClass(goCtx context.Context, msg *types.MsgNewClass) (*types.MsgNewClassResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Assert iscn id is valid
	iscnId, err := iscntypes.ParseIscnId(msg.Parent.IscnIdPrefix)
	if err != nil {
		return nil, types.ErrInvalidIscnId.Wrapf("%s", err.Error())
	}
	// Assert iscn exists
	iscnRecord := k.iscnKeeper.GetContentIdRecord(ctx, iscnId.Prefix)
	if iscnRecord == nil {
		return nil, types.ErrIscnRecordNotFound.Wrapf("ISCN %s not found", iscnId.Prefix.String())
	}
	// Assert current user is owner
	userAddress, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("%s", err.Error())
	}
	if !iscnRecord.OwnerAddress().Equals(userAddress) {
		return nil, sdkerrors.ErrUnauthorized.Wrapf("%s is not the owner of the ISCN %s", msg.Creator, iscnId.Prefix.String())
	}

	// Make class id
	var existingClassIds []string
	value, found := k.GetClassesByISCN(ctx, iscnId.Prefix.String())
	if found {
		existingClassIds = value.ClassIds
	}
	newClassId, err := types.NewClassId(iscnId.Prefix.String(), len(existingClassIds))
	if newClassId == nil || err != nil {
		return nil, types.ErrFailedToMarshalData.Wrapf("%s", err.Error())
	}

	// Create Class
	classData := types.ClassData{
		Metadata: msg.Metadata,
		Parent: types.ClassParent{
			Type:              types.ClassParentType_ISCN,
			IscnIdPrefix:      iscnId.Prefix.String(),
			IscnVersionAtMint: iscnRecord.LatestVersion,
		},
		Config: types.ClassConfig{
			Burnable: msg.Burnable,
		},
	}
	classDataInAny, err := cdctypes.NewAnyWithValue(&classData)
	if err != nil {
		return nil, types.ErrFailedToMarshalData.Wrapf("%s", err.Error())
	}
	class := nft.Class{
		Id:          *newClassId,
		Name:        msg.Name,
		Symbol:      msg.Symbol,
		Description: msg.Description,
		Uri:         msg.Uri,
		UriHash:     msg.UriHash,
		Data:        classDataInAny,
	}
	err = k.nftKeeper.SaveClass(ctx, class)
	if err != nil {
		return nil, types.ErrFailedToSaveClass.Wrapf("%s", err.Error())
	}

	// Append iscn to class mapping
	classIds := append(existingClassIds, *newClassId)
	k.SetClassesByISCN(ctx, types.ClassesByISCN{
		IscnIdPrefix: iscnId.Prefix.String(),
		ClassIds:     classIds,
	})

	// Emit event
	ctx.EventManager().EmitTypedEvent(&types.EventNewClass{
		IscnIdPrefix: iscnId.Prefix.String(),
		ClassId:      *newClassId,
		Owner:        iscnRecord.OwnerAddress().String(),
	})

	return &types.MsgNewClassResponse{
		Class: class,
	}, nil
}
