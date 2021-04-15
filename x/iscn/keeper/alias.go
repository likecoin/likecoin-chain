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

	Params = types.Params
	IscnId = types.IscnId
	CID    = types.CID
)

const (
	TracingIdLength = types.TracingIdLength
	CidCodecType    = types.CidCodecType
)

var (
	ParamKeyRegistryId = types.ParamKeyRegistryId
	ParamKeyFeePerByte = types.ParamKeyFeePerByte
	IscnCountKey       = types.IscnCountKey
	CidBlockKey        = types.CidBlockKey
	IscnIdToCidKey     = types.IscnIdToCidKey

	GetCidBlockKey             = types.GetCidBlockKey
	GetCidToIscnIdKey          = types.GetCidToIscnIdKey
	GetIscnIdToCidKey          = types.GetIscnIdToCidKey
	GetIscnIdVersionKey        = types.GetIscnIdVersionKey
	GetIscnIdOwnerKey          = types.GetIscnIdOwnerKey
	GetFingerprintToCidKey     = types.GetFingerprintToCidKey
	GetFingerprintCidRecordKey = types.GetFingerprintCidRecordKey
)
