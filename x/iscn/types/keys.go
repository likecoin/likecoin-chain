package types

const (
	ModuleName   = "iscn"
	StoreKey     = ModuleName
	QuerierRoute = ModuleName
	RouterKey    = ModuleName
)

var (
	CidBlockKey    = []byte{0x01}
	IscnKernelKey  = []byte{0x02}
	IscnCountKey   = []byte{0x03}
	CidToIscnIDKey = []byte{0x04}
)

func GetIscnKernelKey(iscnID []byte) []byte {
	return append(IscnKernelKey, iscnID...)
}

func GetCidBlockKey(cid []byte) []byte {
	return append(CidBlockKey, cid...)
}

func GetCidToIscnIDKey(cid []byte) []byte {
	return append(CidToIscnIDKey, cid...)
}
