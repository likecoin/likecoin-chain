package keeper

func (k Keeper) MustMarshalIscnId(iscnId IscnId) []byte {
	return k.cdc.MustMarshalBinaryBare(&iscnId)
}

func (k Keeper) MustUnmarshalIscnId(iscnBytes []byte) (iscnId IscnId) {
	k.cdc.MustUnmarshalBinaryBare(iscnBytes, &iscnId)
	return iscnId
}

func (k Keeper) MustMarshalIscnPrefixId(iscnId IscnId) []byte {
	iscnId.Version = 0
	return k.MustMarshalIscnId(iscnId)
}

func (k Keeper) MustMarshalStoreRecord(record *StoreRecord) []byte {
	return k.cdc.MustMarshalBinaryBare(record)
}

func (k Keeper) MustUnmarshalStoreRecord(recordBytes []byte) (record StoreRecord) {
	k.cdc.MustUnmarshalBinaryBare(recordBytes, &record)
	return record
}

func (k Keeper) MustMarshalContentIdRecord(record *ContentIdRecord) []byte {
	return k.cdc.MustMarshalBinaryBare(record)
}

func (k Keeper) MustUnmarshalContentIdRecord(recordBytes []byte) (record ContentIdRecord) {
	k.cdc.MustUnmarshalBinaryBare(recordBytes, &record)
	return record
}
