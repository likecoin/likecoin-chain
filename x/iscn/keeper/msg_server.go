package keeper

import (
	context "context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/likecoin/likechain/x/iscn/types"
)

type msgServer struct {
	Keeper
}

var _ types.MsgServer = msgServer{}

func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

// CreateIscnRecord defines a method to create ISCN record
func (k msgServer) CreateIscnRecord(goCtx context.Context, msg *MsgCreateIscnRecord) (*MsgCreateIscnRecordResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	from, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address: %s", err.Error())
	}
	registryId := k.RegistryId(ctx)
	id := types.GenerateNewIscnIdWithSeed(registryId, ctx.TxBytes())
	if k.GetIscnIdVersion(ctx, id) != 0 {
		return nil, sdkerrors.Wrapf(types.ErrReusingIscnId, "%s", id.String())
	}
	recordJsonLd, err := msg.Record.ToJsonLd(&types.IscnRecordJsonLdInfo{
		Id:         id,
		Timestamp:  ctx.BlockTime(),
		ParentIpld: nil,
	})
	if err != nil {
		return nil, sdkerrors.Wrapf(types.ErrEncodingJsonLd, "%s", err.Error())
	}
	cid, err := k.AddIscnRecord(ctx, id, from, recordJsonLd, msg.Record.ContentFingerprints)
	if err != nil {
		return nil, sdkerrors.Wrapf(types.ErrAddingIscnRecord, "%s", err.Error())
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, from.String()),
		),
	)
	return &types.MsgCreateIscnRecordResponse{
		IscnId:     id.String(),
		RecordIpld: cid.String(),
	}, nil
}

// UpdateIscnRecord defines a method to update existing ISCN record
func (k msgServer) UpdateIscnRecord(goCtx context.Context, msg *MsgUpdateIscnRecord) (*MsgUpdateIscnRecordResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	parentId, err := types.ParseIscnId(msg.IscnId)
	if err != nil {
		return nil, sdkerrors.Wrapf(types.ErrInvalidIscnId, "%s", err.Error())
	}
	currentVersion := k.GetIscnIdVersion(ctx, parentId)
	if parentId.Version != currentVersion {
		return nil, sdkerrors.Wrapf(types.ErrInvalidIscnVersion, "expected version: %d, got: %d", currentVersion, parentId.Version)
	}
	from, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address: %s", err.Error())
	}
	owner := k.GetIscnIdOwner(ctx, parentId)
	if !from.Equals(owner) {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "sender not ISCN record owner, expect %s, got %s", from.String(), owner.String())
	}
	parentCid := k.GetIscnIdCid(ctx, parentId)
	if parentCid == nil {
		return nil, sdkerrors.Wrapf(types.ErrCidNotFound, "%s", parentId.String())
	}
	id := IscnId{
		RegistryId: parentId.RegistryId,
		TracingId:  parentId.TracingId,
		Version:    parentId.Version + 1,
	}
	recordJsonLd, err := msg.Record.ToJsonLd(&types.IscnRecordJsonLdInfo{
		Id:         id,
		Timestamp:  ctx.BlockTime(),
		ParentIpld: parentCid,
	})
	if err != nil {
		return nil, sdkerrors.Wrapf(types.ErrEncodingJsonLd, "%s", err.Error())
	}
	cid, err := k.AddIscnRecord(ctx, id, from, recordJsonLd, msg.Record.ContentFingerprints)
	if err != nil {
		return nil, sdkerrors.Wrapf(types.ErrAddingIscnRecord, "%s", err.Error())
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, from.String()),
		),
	)
	return &types.MsgUpdateIscnRecordResponse{
		IscnId:     parentId.String(),
		RecordIpld: cid.String(),
	}, nil
}

// ChangeIscnRecordOwnership defines a method to update the ownership of existing ISCN record
func (k msgServer) ChangeIscnRecordOwnership(goCtx context.Context, msg *MsgChangeIscnRecordOwnership) (*MsgChangeIscnRecordOwnershipResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	from, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address: %s", err.Error())
	}
	newOwner, err := sdk.AccAddressFromBech32(msg.NewOwner)
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid new owner address: %s", err.Error())
	}
	id, err := types.ParseIscnId(msg.IscnId)
	if err != nil {
		return nil, sdkerrors.Wrapf(types.ErrInvalidIscnId, "%s", err.Error())
	}
	currentVersion := k.GetIscnIdVersion(ctx, id)
	if id.Version != currentVersion {
		return nil, sdkerrors.Wrapf(types.ErrInvalidIscnVersion, "expected version: %d, got: %d", currentVersion, id.Version)
	}
	owner := k.GetIscnIdOwner(ctx, id)
	if !from.Equals(owner) {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "sender not ISCN record owner, expect %s, got %s", from.String(), owner.String())
	}
	if k.GetIscnIdCid(ctx, id) == nil {
		return nil, sdkerrors.Wrapf(types.ErrCidNotFound, "%s", id.String())
	}
	k.SetIscnIdOwner(ctx, id, newOwner)
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeIscnRecord,
			sdk.NewAttribute(types.AttributeKeyIscnId, id.String()),
			sdk.NewAttribute(types.AttributeKeyIscnIdPrefix, id.Prefix()),
			sdk.NewAttribute(types.AttributeKeyIscnOwner, newOwner.String()),
		),
	)
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, from.String()),
		),
	)
	return &types.MsgChangeIscnRecordOwnershipResponse{}, nil
}
