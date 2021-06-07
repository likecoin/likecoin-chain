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
	IscnIdPrefix    = types.IscnIdPrefix
	CID             = types.CID
	IscnInput       = types.IscnInput
	StoreRecord     = types.StoreRecord
	ContentIdRecord = types.ContentIdRecord
)

var (
	ParamKeyRegistryName = types.ParamKeyRegistryName
	ParamKeyFeePerByte   = types.ParamKeyFeePerByte

	SequenceCountKey            = types.SequenceCountKey
	SequenceToStoreRecordPrefix = types.SequenceToStoreRecordPrefix
	CidToSequencePrefix         = types.CidToSequencePrefix
	IscnIdToSequencePrefix      = types.IscnIdToSequencePrefix
	ContentIdRecordPrefix       = types.ContentIdRecordPrefix
	FingerprintSequencePrefix   = types.FingerprintSequencePrefix
	OwnerSequencePrefix         = types.OwnerSequencePrefix

	NewIscnId = types.NewIscnId
)
