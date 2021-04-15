package types

func NewQueryIscnRecordsRequestByIscnId(iscnId IscnId) *QueryIscnRecordsRequest {
	return &QueryIscnRecordsRequest{
		IscnId: iscnId.String(),
	}
}

func NewQueryIscnRecordsRequestByFingerprint(fingerprint string) *QueryIscnRecordsRequest {
	return &QueryIscnRecordsRequest{
		Fingerprint: fingerprint,
	}
}

func NewQueryParamsRequest() *QueryParamsRequest {
	return &QueryParamsRequest{}
}
