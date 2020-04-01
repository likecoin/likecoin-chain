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
		default:
			errMsg := fmt.Sprintf("unrecognized iscn message type: %T", msg)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleMsgCreateIscn(ctx sdk.Context, msg MsgCreateIscn, keeper Keeper) sdk.Result {
	id, err := keeper.AddIscnRecord(ctx, msg.From, &msg.IscnRecord)
	if err != nil {
		return sdk.Result{ /* TODO: error*/ }
	}
	idStr := base64.URLEncoding.EncodeToString(id) // TODO: formatting iscn
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			EventTypeCreateIscn,
			sdk.NewAttribute(AttributeKeyIscn, idStr),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.From.String()),
		),
	})

	return sdk.Result{Events: ctx.EventManager().Events()}
}
