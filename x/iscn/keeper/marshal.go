package keeper

func (k Keeper) MustMarshalIscnId(iscnId IscnId) []byte {
	return k.cdc.MustMarshal(&iscnId)
}

func (k Keeper) MustUnmarshalIscnId(iscnBytes []byte) (iscnId IscnId) {
	k.cdc.MustUnmarshal(iscnBytes, &iscnId)
	return iscnId
}

func (k Keeper) MustMarshalIscnIdPrefix(iscnIdPrefix IscnIdPrefix) []byte {
	return k.cdc.MustMarshal(&iscnIdPrefix)
}

func (k Keeper) MustUnmarshalIscnIdPrefix(iscnIdPrefixBytes []byte) (iscnIdPrefix IscnIdPrefix) {
	k.cdc.MustUnmarshal(iscnIdPrefixBytes, &iscnIdPrefix)
	return iscnIdPrefix
}

func (k Keeper) MustMarshalStoreRecord(record *StoreRecord) []byte {
	return k.cdc.MustMarshal(record)
}

func (k Keeper) MustUnmarshalStoreRecord(recordBytes []byte) (record StoreRecord) {
	k.cdc.MustUnmarshal(recordBytes, &record)
	return record
}

func (k Keeper) MustMarshalContentIdRecord(record *ContentIdRecord) []byte {
	return k.cdc.MustMarshal(record)
}

func (k Keeper) MustUnmarshalContentIdRecord(recordBytes []byte) (record ContentIdRecord) {
	k.cdc.MustUnmarshal(recordBytes, &record)
	return record
}
