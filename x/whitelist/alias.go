package whitelist

import (
	"github.com/likecoin/likechain/x/whitelist/types"
)

const (
	ModuleName     = types.ModuleName
	StoreKey       = types.StoreKey
	QuerierRoute   = types.QuerierRoute
	RouterKey      = types.RouterKey
	QueryApprover  = types.QueryApprover
	QueryWhitelist = types.QueryWhitelist
)

var (
	ModuleCdc                   = types.ModuleCdc
	NewMsgSetWhitelist          = types.NewMsgSetWhitelist
	ErrInvalidApprover          = types.ErrInvalidApprover
	ErrValidatorNotInWEhitelist = types.ErrValidatorNotInWEhitelist
	KeyApprover                 = types.KeyApprover
	DefaultParams               = types.DefaultParams
	DefaultGenesisState         = types.DefaultGenesisState
	DefaultCodespace            = types.DefaultCodespace
	ValidateGenesis             = types.ValidateGenesis
	WhitelistKey                = types.WhitelistKey
	EventTypeSetWhitelist       = types.EventTypeSetWhitelist
	AttributeKeyWhitelist       = types.AttributeKeyWhitelist
	AttributeValueCategory      = types.AttributeValueCategory
	RegisterCodec               = types.RegisterCodec
)

type (
	MsgSetWhitelist = types.MsgSetWhitelist
	Whitelist       = types.Whitelist
	Params          = types.Params
	GenesisState    = types.GenesisState
)
