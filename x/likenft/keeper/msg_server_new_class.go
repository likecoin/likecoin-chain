package keeper

import (
	"context"
	"fmt"

	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/nft"
	"github.com/likecoin/likecoin-chain/v4/x/likenft/types"
)

func (k msgServer) NewClass(goCtx context.Context, msg *types.MsgNewClass) (*types.MsgNewClassResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	parent, err := k.NewClassParentFromInput(ctx, msg.Parent, msg.Creator)
	if err != nil {
		return nil, err
	}

	// check user is parent owner
	if err := k.assertBech32EqualsAccAddress(msg.Creator, parent.Owner); err != nil {
		return nil, err
	}

	// Sanitize class config
	cleanClassConfig, err := k.sanitizeClassConfig(ctx, msg.Input.Config, 0)
	if cleanClassConfig == nil || err != nil {
		return nil, err
	}
	msg.Input.Config = *cleanClassConfig

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
		Metadata: msg.Input.Metadata,
		Parent:   parent.ClassParent,
		Config:   msg.Input.Config,
		BlindBoxState: types.BlindBoxState{
			ToBeRevealed: msg.Input.Config.IsBlindBox(),
		},
	}
	classDataInAny, err := cdctypes.NewAnyWithValue(&classData)
	if err != nil {
		return nil, types.ErrFailedToMarshalData.Wrapf("%s", err.Error())
	}
	class := nft.Class{
		Id:          newClassId,
		Name:        msg.Input.Name,
		Symbol:      msg.Input.Symbol,
		Description: msg.Input.Description,
		Uri:         msg.Input.Uri,
		UriHash:     msg.Input.UriHash,
		Data:        classDataInAny,
	}
	// Deduct fee
	err = k.DeductFeePerByte(ctx, parent.Owner, class.Size(), nil)
	if err != nil {
		return nil, err
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

	// Enqueue class for reveal
	if classData.Config.IsBlindBox() {
		k.SetClassRevealQueueEntry(ctx, types.ClassRevealQueueEntry{
			ClassId:    newClassId,
			RevealTime: classData.Config.BlindBoxConfig.RevealTime,
		})
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
