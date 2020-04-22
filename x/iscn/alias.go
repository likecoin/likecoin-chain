package iscn

import (
	"github.com/likecoin/likechain/x/iscn/types"
)

const (
	ModuleName      = types.ModuleName
	StoreKey        = types.StoreKey
	QuerierRoute    = types.QuerierRoute
	RouterKey       = types.RouterKey
	QueryAuthor     = types.QueryAuthor
	QueryIscnRecord = types.QueryIscnRecord
	QueryParams     = types.QueryParams
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
	IscnCountKey                = types.IscnCountKey
	AuthorKey                   = types.AuthorKey
	RightTermsKey               = types.RightTermsKey
	GetIscnRecordKey            = types.GetIscnRecordKey
	GetAuthorKey                = types.GetAuthorKey
	GetRightTermsKey            = types.GetRightTermsKey
	EventTypeCreateIscn         = types.EventTypeCreateIscn
	EventTypeAddAuthor          = types.EventTypeAddAuthor
	AttributeKeyIscnId          = types.AttributeKeyIscnId
	AttributeKeyAuthorCid       = types.AttributeKeyAuthorCid
	AttributeValueCategory      = types.AttributeValueCategory
	RegisterCodec               = types.RegisterCodec
)

type (
	MsgCreateIscn = types.MsgCreateIscn
	MsgAddAuthor  = types.MsgAddAuthor
	IscnRecord    = types.IscnRecord
	Author        = types.Author
	RightTerms    = types.RightTerms
	Params        = types.Params
	GenesisState  = types.GenesisState
)
