package types

const (
	ModuleName   = "iscn"
	StoreKey     = ModuleName
	QuerierRoute = ModuleName
	RouterKey    = ModuleName
)

var (
	IscnRecordKey = []byte{0x01}
)

func GetIscnRecordKey(iscnId []byte) []byte {
	return append(IscnRecordKey, iscnId...)
}
