package keeper

import (
	"context"
	"fmt"

	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/likecoin/likechain/backport/cosmos-sdk/v0.46.0-alpha2/x/nft"
	iscntypes "github.com/likecoin/likechain/x/iscn/types"
	"github.com/likecoin/likechain/x/likenft/types"
)

func (k msgServer) NewClass(goCtx context.Context, msg *types.MsgNewClass) (*types.MsgNewClassResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	var safeParent types.ClassParent
	var safeUserAddress sdk.AccAddress
	if msg.Parent.Type == types.ClassParentType_ISCN {
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
		safeUserAddress, err = sdk.AccAddressFromBech32(msg.Creator)
		if err != nil {
			return nil, sdkerrors.ErrInvalidAddress.Wrapf("%s", err.Error())
		}
		if !iscnRecord.OwnerAddress().Equals(safeUserAddress) {
			return nil, sdkerrors.ErrUnauthorized.Wrapf("%s is not the owner of the ISCN %s", msg.Creator, iscnId.Prefix.String())
		}
		safeParent = types.ClassParent{
			Type:              types.ClassParentType_ISCN,
			IscnIdPrefix:      iscnId.Prefix.String(),
			IscnVersionAtMint: iscnRecord.LatestVersion,
		}
	} else if msg.Parent.Type == types.ClassParentType_ACCOUNT {
		var err error
		safeUserAddress, err = sdk.AccAddressFromBech32(msg.Creator)
		if err != nil {
			return nil, sdkerrors.ErrInvalidAddress.Wrapf("%s", err.Error())
		}
		safeParent = types.ClassParent{
			Type:    types.ClassParentType_ACCOUNT,
			Account: safeUserAddress.String(),
		}
	} else {
		return nil, sdkerrors.ErrInvalidRequest.Wrapf("Unsupported parent type %s", msg.Parent.Type.String())
	}

	// Make class id
	var existingClassIds []string
	var newClassId string
	if safeParent.Type == types.ClassParentType_ISCN {
		value, found := k.GetClassesByISCN(ctx, safeParent.IscnIdPrefix)
		if found {
			existingClassIds = value.ClassIds
		}
		var err error
		newClassId, err = types.NewClassIdForISCN(safeParent.IscnIdPrefix, len(existingClassIds))
		if newClassId == "" || err != nil {
			return nil, types.ErrFailedToMarshalData.Wrapf("%s", err.Error())
		}
	} else if safeParent.Type == types.ClassParentType_ACCOUNT {
		value, found := k.GetClassesByAccount(ctx, safeUserAddress)
		if found {
			existingClassIds = value.ClassIds
		}
		var err error
		newClassId, err = types.NewClassIdForAccount(safeUserAddress, len(existingClassIds))
		if newClassId == "" || err != nil {
			return nil, types.ErrFailedToMarshalData.Wrapf("%s", err.Error())
		}
	} else {
		panic(fmt.Sprintf("Unsupported parent type %s after initial check", safeParent.Type.String()))
	}

	// Create Class
	classData := types.ClassData{
		Metadata: msg.Metadata,
		Parent:   safeParent,
		Config: types.ClassConfig{
			Burnable: msg.Burnable,
		},
	}
	classDataInAny, err := cdctypes.NewAnyWithValue(&classData)
	if err != nil {
		return nil, types.ErrFailedToMarshalData.Wrapf("%s", err.Error())
	}
	class := nft.Class{
		Id:          newClassId,
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
	classIds := append(existingClassIds, newClassId)
	if safeParent.Type == types.ClassParentType_ISCN {
		k.SetClassesByISCN(ctx, types.ClassesByISCN{
			IscnIdPrefix: safeParent.IscnIdPrefix,
			ClassIds:     classIds,
		})
	} else if safeParent.Type == types.ClassParentType_ACCOUNT {
		k.SetClassesByAccount(ctx, types.ClassesByAccount{
			Account:  safeParent.Account,
			ClassIds: classIds,
		})
	} else {
		panic(fmt.Sprintf("Unsupported parent type %s after initial check", safeParent.Type.String()))
	}

	// Emit event
	ctx.EventManager().EmitTypedEvent(&types.EventNewClass{
		IscnIdPrefix: safeParent.IscnIdPrefix,
		ClassId:      newClassId,
		Owner:        safeUserAddress.String(),
	})

	return &types.MsgNewClassResponse{
		Class: class,
	}, nil
}
