package types

const (
	ModuleName   = "iscn"
	StoreKey     = ModuleName
	QuerierRoute = ModuleName
	RouterKey    = ModuleName
)

var (
	IscnRecordKey = []byte{0x01}
	IscnCountKey  = []byte{0x02}
	AuthorKey     = []byte{0x03}
	RightTermsKey = []byte{0x04}
)

func GetIscnRecordKey(iscnId []byte) []byte {
	return append(IscnRecordKey, iscnId...)
}

func GetAuthorKey(authorCid []byte) []byte {
	return append(AuthorKey, authorCid...)
}

func GetRightTermsKey(rightTermsHash []byte) []byte {
	return append(RightTermsKey, rightTermsHash...)
}
