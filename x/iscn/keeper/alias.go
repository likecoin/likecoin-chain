package keeper

import (
	"github.com/likecoin/likechain/x/iscn/types"
)

type (
	MsgCreateIscnRecord                  = types.MsgCreateIscnRecord
	MsgCreateIscnRecordResponse          = types.MsgCreateIscnRecordResponse
	MsgUpdateIscnRecord                  = types.MsgUpdateIscnRecord
	MsgUpdateIscnRecordResponse          = types.MsgUpdateIscnRecordResponse
	MsgChangeIscnRecordOwnership         = types.MsgChangeIscnRecordOwnership
	MsgChangeIscnRecordOwnershipResponse = types.MsgChangeIscnRecordOwnershipResponse

	Params          = types.Params
	IscnId          = types.IscnId
	CID             = types.CID
	IscnInput       = types.IscnInput
	StoreRecord     = types.StoreRecord
	TracingIdRecord = types.TracingIdRecord
)

const (
	TracingIdLength = types.TracingIdLength
	CidCodecType    = types.CidCodecType
)

var (
	ParamKeyRegistryId = types.ParamKeyRegistryId
	ParamKeyFeePerByte = types.ParamKeyFeePerByte

	SequenceCountKey            = types.SequenceCountKey
	SequenceToStoreRecordPrefix = types.SequenceToStoreRecordPrefix
	CidToSequencePrefix         = types.CidToSequencePrefix
	IscnIdToSequencePrefix      = types.IscnIdToSequencePrefix
	TracingIdRecordPrefix       = types.TracingIdRecordPrefix
	FingerprintSequencePrefix   = types.FingerprintSequencePrefix
)
