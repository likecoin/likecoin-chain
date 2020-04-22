package iscn

import (
	"encoding/base64"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgCreateIscn:
			return handleMsgCreateIscn(ctx, msg, keeper)
		case MsgAddAuthor:
			return handleMsgAddAuthor(ctx, msg, keeper)
		default:
			errMsg := fmt.Sprintf("unrecognized iscn message type: %T", msg)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleMsgCreateIscn(ctx sdk.Context, msg MsgCreateIscn, keeper Keeper) sdk.Result {
	id, err := keeper.AddIscnRecord(ctx, msg.From, &msg.IscnRecord)
	if err != nil {
		return sdk.Result{
			/* TODO: proper error*/
			Code:      err.Code(),
			Codespace: err.Codespace(),
			Log:       err.Error(),
		}
	}
	idStr := base64.URLEncoding.EncodeToString(id) // TODO: formatting iscn
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			EventTypeCreateIscn,
			sdk.NewAttribute(AttributeKeyIscnId, idStr),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.From.String()),
		),
	})

	return sdk.Result{Events: ctx.EventManager().Events()}
}

func handleMsgAddAuthor(ctx sdk.Context, msg MsgAddAuthor, keeper Keeper) sdk.Result {
	// TODO: extract so we can reuse logic in MsgCreateIscn
	authorCid := keeper.SetAuthor(ctx, &msg.AuthorInfo)
	cidStr := base64.URLEncoding.EncodeToString(authorCid) // TODO: formatting cid
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			EventTypeAddAuthor,
			sdk.NewAttribute(AttributeKeyAuthorCid, cidStr),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.From.String()),
		),
	})

	return sdk.Result{Events: ctx.EventManager().Events()}
}
