package iscn

import (
	"github.com/likecoin/likechain/x/iscn/types"
)

const (
	ModuleName   = types.ModuleName
	StoreKey     = types.StoreKey
	QuerierRoute = types.QuerierRoute
	RouterKey    = types.RouterKey
	QueryRecord  = types.QueryRecord
	QueryParams  = types.QueryParams
)

var (
	ModuleCdc                   = types.ModuleCdc
	NewMsgCreateIscn            = types.NewMsgCreateIscn
	ErrInvalidApprover          = types.ErrInvalidApprover
	ErrValidatorNotInWEhitelist = types.ErrValidatorNotInWEhitelist
	KeyFeePerByte               = types.KeyFeePerByte
	DefaultParams               = types.DefaultParams
	DefaultGenesisState         = types.DefaultGenesisState
	DefaultCodespace            = types.DefaultCodespace
	ValidateGenesis             = types.ValidateGenesis
	IscnRecordKey               = types.IscnRecordKey
	GetIscnRecordKey            = types.GetIscnRecordKey
	EventTypeCreateIscn         = types.EventTypeCreateIscn
	AttributeKeyIscn            = types.AttributeKeyIscn
	AttributeValueCategory      = types.AttributeValueCategory
	RegisterCodec               = types.RegisterCodec
)

type (
	MsgCreateIscn = types.MsgCreateIscn
	IscnRecord    = types.IscnRecord
	Params        = types.Params
	GenesisState  = types.GenesisState
)
