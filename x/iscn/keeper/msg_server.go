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
	registryName := k.RegistryName(ctx)
	id := types.GenerateNewIscnIdWithSeed(registryName, ctx.TxBytes())
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
		return nil, err
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
	from, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address: %s", err.Error())
	}
	parentId, err := types.ParseIscnId(msg.IscnId)
	if err != nil {
		return nil, sdkerrors.Wrapf(types.ErrInvalidIscnId, "%s", err.Error())
	}
	parentSeq := k.GetIscnIdSequence(ctx, parentId)
	if parentSeq == 0 {
		return nil, sdkerrors.Wrapf(types.ErrRecordNotFound, "parent ISCN ID %s not found", parentId.String())
	}
	parentStoreRecord := k.GetStoreRecord(ctx, parentSeq)
	parentCid := parentStoreRecord.Cid()
	id := NewIscnId(parentId.Prefix.RegistryName, parentId.Prefix.ContentId, parentId.Version+1)
	recordJsonLd, err := msg.Record.ToJsonLd(&types.IscnRecordJsonLdInfo{
		Id:         id,
		Timestamp:  ctx.BlockTime(),
		ParentIpld: &parentCid,
	})
	if err != nil {
		return nil, sdkerrors.Wrapf(types.ErrEncodingJsonLd, "%s", err.Error())
	}
	cid, err := k.AddIscnRecord(ctx, id, from, recordJsonLd, msg.Record.ContentFingerprints)
	if err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, from.String()),
		),
	)
	return &types.MsgUpdateIscnRecordResponse{
		IscnId:     id.String(),
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
	contentIdRecord := k.GetContentIdRecord(ctx, id)
	if contentIdRecord == nil {
		return nil, sdkerrors.Wrapf(types.ErrRecordNotFound, "%s", id.String())
	}
	if id.Version != contentIdRecord.LatestVersion {
		return nil, sdkerrors.Wrapf(types.ErrInvalidIscnVersion, "expected version: %d", contentIdRecord.LatestVersion)
	}
	prevOwner := contentIdRecord.OwnerAddress()
	if !from.Equals(prevOwner) {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "sender not ISCN record owner, expect %s", prevOwner.String())
	}
	contentIdRecord.OwnerAddressBytes = newOwner.Bytes()
	k.SetContentIdRecord(ctx, id, contentIdRecord)
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeIscnRecord,
			sdk.NewAttribute(types.AttributeKeyIscnId, id.String()),
			sdk.NewAttribute(types.AttributeKeyIscnIdPrefix, id.Prefix.String()),
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
