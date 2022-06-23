package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/likecoin/likechain/x/likenft/types"
)

func (k msgServer) CreateRoyaltyConfig(goCtx context.Context, msg *types.MsgCreateRoyaltyConfig) (*types.MsgCreateRoyaltyConfigResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Check user is class owner

	// Check if the value already exists
	_, isFound := k.GetRoyaltyConfig(
		ctx,
		msg.ClassId,
	)
	if isFound {
		// TODO: customize error
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "index already set")
	}

	var royaltyConfigByClass = types.RoyaltyConfigByClass{
		ClassId:       msg.ClassId,
		RoyaltyConfig: msg.RoyaltyConfig.ToConfig(),
	}

	k.SetRoyaltyConfig(
		ctx,
		royaltyConfigByClass,
	)

	// TODO emit event

	return &types.MsgCreateRoyaltyConfigResponse{
		RoyaltyConfig: royaltyConfigByClass.RoyaltyConfig,
	}, nil
}

func (k msgServer) UpdateRoyaltyConfig(goCtx context.Context, msg *types.MsgUpdateRoyaltyConfig) (*types.MsgUpdateRoyaltyConfigResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Check user is class owner

	// Check if the value exists
	_, isFound := k.GetRoyaltyConfig(
		ctx,
		msg.ClassId,
	)
	if !isFound {
		// TODO: customize error
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, "index not set")
	}

	var royaltyConfigByClass = types.RoyaltyConfigByClass{
		ClassId:       msg.ClassId,
		RoyaltyConfig: msg.RoyaltyConfig.ToConfig(),
	}

	k.SetRoyaltyConfig(ctx, royaltyConfigByClass)

	// TODO emit event

	return &types.MsgUpdateRoyaltyConfigResponse{
		RoyaltyConfig: royaltyConfigByClass.RoyaltyConfig,
	}, nil
}

func (k msgServer) DeleteRoyaltyConfig(goCtx context.Context, msg *types.MsgDeleteRoyaltyConfig) (*types.MsgDeleteRoyaltyConfigResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Check user is class owner

	// Check if the value exists
	_, isFound := k.GetRoyaltyConfig(
		ctx,
		msg.ClassId,
	)
	if !isFound {
		// TODO: customize error
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, "index not set")
	}

	k.RemoveRoyaltyConfig(
		ctx,
		msg.ClassId,
	)

	// TODO emit event

	return &types.MsgDeleteRoyaltyConfigResponse{}, nil
}
