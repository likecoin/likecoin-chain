package types

func (record StoreRecord) Cid() CID {
	return MustCidFromBytes(record.CidBytes)
}
