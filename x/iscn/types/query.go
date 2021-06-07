package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewQueryRecordsByIdRequest(iscnId IscnId, fromVersion, toVersion uint64) *QueryRecordsByIdRequest {
	return &QueryRecordsByIdRequest{
		IscnId:      iscnId.String(),
		FromVersion: fromVersion,
		ToVersion:   toVersion,
	}
}

func NewQueryRecordsByFingerprintRequest(fingerprint string, fromSeq uint64) *QueryRecordsByFingerprintRequest {
	return &QueryRecordsByFingerprintRequest{
		Fingerprint:  fingerprint,
		FromSequence: fromSeq,
	}
}

func NewQueryRecordsByOwnerRequest(owner sdk.AccAddress, fromSeq uint64) *QueryRecordsByOwnerRequest {
	return &QueryRecordsByOwnerRequest{
		Owner:        owner.String(),
		FromSequence: fromSeq,
	}
}

func NewQueryParamsRequest() *QueryParamsRequest {
	return &QueryParamsRequest{}
}
