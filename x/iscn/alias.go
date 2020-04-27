package iscn

import (
	"github.com/likecoin/likechain/x/iscn/types"
)

const (
	ModuleName           = types.ModuleName
	StoreKey             = types.StoreKey
	QuerierRoute         = types.QuerierRoute
	RouterKey            = types.RouterKey
	QueryIscnKernel      = types.QueryIscnKernel
	QueryParams          = types.QueryParams
	QueryCID             = types.QueryCID
	QueryCidBlockGet     = types.QueryCidBlockGet
	QueryCidBlockGetSize = types.QueryCidBlockGetSize
	QueryCidBlockHas     = types.QueryCidBlockHas
)

var (
	ModuleCdc                  = types.ModuleCdc
	NewMsgCreateIscn           = types.NewMsgCreateIscn
	ErrInvalidSender           = types.ErrInvalidSender
	KeyFeePerByte              = types.KeyFeePerByte
	DefaultParams              = types.DefaultParams
	DefaultGenesisState        = types.DefaultGenesisState
	DefaultCodespace           = types.DefaultCodespace
	ValidateGenesis            = types.ValidateGenesis
	CidBlockKey                = types.CidBlockKey
	IscnKernelKey              = types.IscnKernelKey
	IscnCountKey               = types.IscnCountKey
	CidToIscnIDKey             = types.CidToIscnIDKey
	GetCidBlockKey             = types.GetCidBlockKey
	GetIscnKernelKey           = types.GetIscnKernelKey
	GetCidToIscnIDKey          = types.GetCidToIscnIDKey
	IscnKernelCodecType        = types.IscnKernelCodecType
	IscnContentCodecType       = types.IscnContentCodecType
	RightTermsCodecType        = types.RightTermsCodecType
	EntityCodecType            = types.EntityCodecType
	StakeholdersCodecType      = types.StakeholdersCodecType
	RightsCodecType            = types.RightsCodecType
	EventTypeCreateIscn        = types.EventTypeCreateIscn
	EventTypeAddEntity         = types.EventTypeAddEntity
	EventTypeAddRightTerms     = types.EventTypeAddRightTerms
	EventTypeAddIscnContent    = types.EventTypeAddIscnContent
	EventTypeAddIscnKernel     = types.EventTypeAddIscnKernel
	AttributeKeyIscnID         = types.AttributeKeyIscnID
	AttributeKeyIscnKernelCid  = types.AttributeKeyIscnKernelCid
	AttributeKeyIscnContentCid = types.AttributeKeyIscnContentCid
	AttributeKeyEntityCid      = types.AttributeKeyEntityCid
	AttributeKeyRightTermsCid  = types.AttributeKeyRightTermsCid
	AttributeValueCategory     = types.AttributeValueCategory
	RegisterCodec              = types.RegisterCodec
	CidMbaseEncoder            = types.CidMbaseEncoder
	CheckIscnType              = types.CheckIscnType
	None                       = types.None
	Number                     = types.Number
	String                     = types.String
	NestedCID                  = types.NestedCID
	NestedIscnData             = types.NestedIscnData
	Array                      = types.Array
	Unknown                    = types.Unknown
	KernelSchema               = types.KernelSchema
	EntitySchema               = types.EntitySchema
	RegistryID                 = types.RegistryID
)

type (
	MsgCreateIscn    = types.MsgCreateIscn
	MsgAddEntity     = types.MsgAddEntity
	MsgAddRightTerms = types.MsgAddRightTerms
	CID              = types.CID
	Params           = types.Params
	GenesisState     = types.GenesisState
	RawIscnMap       = types.RawIscnMap
	IscnData         = types.IscnData
	IscnDataField    = types.IscnDataField
	IscnDataArray    = types.IscnDataArray
	IscnID           = types.IscnID
)
