package whitelist

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
)

func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgSetWhitelist:
			return handleMsgSetWhitelist(ctx, msg, keeper)
		default:
			errMsg := fmt.Sprintf("unrecognized whitelist message type: %T", msg)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleMsgSetWhitelist(ctx sdk.Context, msg MsgSetWhitelist, keeper Keeper) sdk.Result {
	approver := keeper.Approver(ctx)
	if !approver.Equals(msg.Approver) {
		return ErrInvalidApprover(keeper.Codespace()).Result()
	}
	keeper.SetWhitelist(ctx, msg.Whitelist)
	bz, err := json.Marshal(msg.Whitelist)
	if err != nil {
		panic(err)
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			EventTypeSetWhitelist,
			sdk.NewAttribute(AttributeKeyWhitelist, string(bz)),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Approver.String()),
		),
	})

	return sdk.Result{Events: ctx.EventManager().Events()}
}

func WrapStakingHandler(keeper Keeper, stakingHandler sdk.Handler) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case staking.MsgCreateValidator:
			result := checkWhitelist(ctx, keeper, msg)
			if result.Code != 0 {
				return result
			}
		}
		return stakingHandler(ctx, msg)
	}
}

func checkWhitelist(ctx sdk.Context, keeper Keeper, msg staking.MsgCreateValidator) sdk.Result {
	whitelist := keeper.GetWhitelist(ctx)
	if len(whitelist) > 0 {
		for _, v := range whitelist {
			if msg.ValidatorAddress.Equals(v) {
				return sdk.Result{Code: 0}
			}
		}
		return ErrValidatorNotInWEhitelist(keeper.Codespace()).Result()
	}
	return sdk.Result{}
}
