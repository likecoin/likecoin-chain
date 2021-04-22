package types

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

func NewQueryParamsRequest() *QueryParamsRequest {
	return &QueryParamsRequest{}
}
