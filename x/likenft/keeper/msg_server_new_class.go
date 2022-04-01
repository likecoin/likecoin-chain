package keeper

import (
	"context"
	"fmt"

	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/likecoin/likechain/backport/cosmos-sdk/v0.46.0-alpha2/x/nft"
	"github.com/likecoin/likechain/x/likenft/types"
)

func (k msgServer) NewClass(goCtx context.Context, msg *types.MsgNewClass) (*types.MsgNewClassResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	parent, err := k.resolveClassParentAndOwner(ctx, msg.Parent, msg.Creator)
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

	// Make class id
	var existingClassIds []string
	var newClassId string
	if parent.Type == types.ClassParentType_ISCN {
		value, found := k.GetClassesByISCN(ctx, parent.IscnIdPrefix)
		if found {
			existingClassIds = value.ClassIds
		}
		var err error
		newClassId, err = types.NewClassIdForISCN(parent.IscnIdPrefix, len(existingClassIds))
		if newClassId == "" || err != nil {
			return nil, types.ErrFailedToMarshalData.Wrapf("%s", err.Error())
		}
	} else if parent.Type == types.ClassParentType_ACCOUNT {
		value, found := k.GetClassesByAccount(ctx, parent.Owner)
		if found {
			existingClassIds = value.ClassIds
		}
		var err error
		newClassId, err = types.NewClassIdForAccount(parent.Owner, len(existingClassIds))
		if newClassId == "" || err != nil {
			return nil, types.ErrFailedToMarshalData.Wrapf("%s", err.Error())
		}
	} else {
		panic(fmt.Sprintf("Unsupported parent type %s after initial check", parent.Type.String()))
	}

	// Create Class
	classData := types.ClassData{
		Metadata: msg.Metadata,
		Parent:   parent.ClassParent,
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
	if parent.Type == types.ClassParentType_ISCN {
		k.SetClassesByISCN(ctx, types.ClassesByISCN{
			IscnIdPrefix: parent.IscnIdPrefix,
			ClassIds:     classIds,
		})
	} else if parent.Type == types.ClassParentType_ACCOUNT {
		k.SetClassesByAccount(ctx, types.ClassesByAccount{
			Account:  parent.Account,
			ClassIds: classIds,
		})
	} else {
		panic(fmt.Sprintf("Unsupported parent type %s after initial check", parent.Type.String()))
	}

	// Emit event
	ctx.EventManager().EmitTypedEvent(&types.EventNewClass{
		ClassId:            newClassId,
		ParentIscnIdPrefix: classData.Parent.IscnIdPrefix,
		ParentAccount:      classData.Parent.Account,
	})

	return &types.MsgNewClassResponse{
		Class: class,
	}, nil
}
