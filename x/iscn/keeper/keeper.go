package keeper

import (
	"encoding/binary"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	authTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	paramTypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/likecoin/likechain/x/iscn/types"
)

type AccountKeeper interface {
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) authTypes.AccountI
}

type Keeper struct {
	storeKey      sdk.StoreKey
	cdc           codec.BinaryMarshaler
	paramstore    paramTypes.Subspace
	accountKeeper AccountKeeper
	bankKeeper    authTypes.BankKeeper
}

func NewKeeper(
	cdc codec.BinaryMarshaler, key sdk.StoreKey, accountKeeper AccountKeeper,
	bankKeeper authTypes.BankKeeper, paramstore paramTypes.Subspace,
) Keeper {
	return Keeper{
		storeKey:      key,
		cdc:           cdc,
		paramstore:    paramstore.WithKeyTable(ParamKeyTable()),
		accountKeeper: accountKeeper,
		bankKeeper:    bankKeeper,
	}
}

func ParamKeyTable() paramTypes.KeyTable {
	return paramTypes.NewKeyTable().RegisterParamSet(&Params{})
}

func (k Keeper) RegistryId(ctx sdk.Context) (res string) {
	k.paramstore.Get(ctx, ParamKeyRegistryId, &res)
	return
}

func (k Keeper) FeePerByte(ctx sdk.Context) (res sdk.DecCoin) {
	k.paramstore.Get(ctx, ParamKeyFeePerByte, &res)
	return
}

func (k Keeper) GetParams(ctx sdk.Context) Params {
	return Params{
		RegistryId: k.RegistryId(ctx),
		FeePerByte: k.FeePerByte(ctx),
	}
}

func (k Keeper) SetParams(ctx sdk.Context, params Params) {
	k.paramstore.SetParamSet(ctx, &params)
}

func (k Keeper) GetCidBlock(ctx sdk.Context, cid CID) []byte {
	key := GetCidBlockKey(cid)
	return ctx.KVStore(k.storeKey).Get(key)
}

func (k Keeper) HasCidBlock(ctx sdk.Context, cid CID) bool {
	key := GetCidBlockKey(cid)
	return ctx.KVStore(k.storeKey).Has(key)
}

func (k Keeper) SetCidBlock(ctx sdk.Context, cid CID, bz []byte) {
	key := GetCidBlockKey(cid)
	ctx.KVStore(k.storeKey).Set(key, bz)
}

func (k Keeper) IterateCidBlocks(ctx sdk.Context, f func(cid CID, bz []byte) bool) {
	it := sdk.KVStorePrefixIterator(ctx.KVStore(k.storeKey), CidBlockKey)
	defer it.Close()
	for ; it.Valid(); it.Next() {
		cidBytes := it.Key()[len(CidBlockKey):]
		cid := types.MustCidFromBytes(cidBytes)
		if f(cid, it.Value()) {
			break
		}
	}
}

func (k Keeper) GetCidIscnId(ctx sdk.Context, cid CID) *IscnId {
	key := GetCidToIscnIdKey(cid)
	bz := ctx.KVStore(k.storeKey).Get(key)
	if bz == nil {
		return nil
	}
	iscnId := IscnId{}
	k.cdc.MustUnmarshalBinaryBare(bz, &iscnId)
	return &iscnId
}

func (k Keeper) SetCidIscnId(ctx sdk.Context, cid CID, iscnId IscnId) {
	key := GetCidToIscnIdKey(cid)
	bz := k.cdc.MustMarshalBinaryBare(&iscnId)
	ctx.KVStore(k.storeKey).Set(key, bz)
}

func (k Keeper) GetIscnIdCid(ctx sdk.Context, iscnId IscnId) *CID {
	key := GetIscnIdToCidKey(k.cdc, iscnId)
	bz := ctx.KVStore(k.storeKey).Get(key)
	if bz == nil {
		return nil
	}
	cid := types.MustCidFromBytes(bz)
	return &cid
}

func (k Keeper) SetIscnIdCid(ctx sdk.Context, iscnId IscnId, cid CID) {
	key := GetIscnIdToCidKey(k.cdc, iscnId)
	ctx.KVStore(k.storeKey).Set(key, cid.Bytes())
}

func (k Keeper) IterateIscnIds(ctx sdk.Context, f func(iscnId IscnId, cid CID) bool) {
	it := sdk.KVStorePrefixIterator(ctx.KVStore(k.storeKey), IscnIdToCidKey)
	defer it.Close()
	for ; it.Valid(); it.Next() {
		iscnId := IscnId{}
		iscnIdBytes := it.Key()[len(IscnIdToCidKey):]
		k.cdc.MustUnmarshalBinaryBare(iscnIdBytes, &iscnId)
		cid := types.MustCidFromBytes(it.Value())
		if f(iscnId, cid) {
			break
		}
	}
}

func (k Keeper) AddFingerprintCid(ctx sdk.Context, fingerprint string, cid CID) {
	key := GetFingerprintCidRecordKey(fingerprint, cid)
	ctx.KVStore(k.storeKey).Set(key, []byte{0x01})
}

func (k Keeper) HasFingerprintCid(ctx sdk.Context, fingerprint string, cid CID) bool {
	key := GetFingerprintCidRecordKey(fingerprint, cid)
	return ctx.KVStore(k.storeKey).Has(key)
}

func (k Keeper) IterateFingerprints(ctx sdk.Context, f func(fingerprint string, cid CID) bool) {
	prefix := types.FingerprintToCidKey
	it := sdk.KVStorePrefixIterator(ctx.KVStore(k.storeKey), prefix)
	defer it.Close()
	for ; it.Valid(); it.Next() {
		key := it.Key()
		fingerprintLenBytes := key[len(prefix) : len(prefix)+4]
		fingerprintLen := binary.BigEndian.Uint32(fingerprintLenBytes)
		fingerprint := string(key[len(prefix)+4 : len(prefix)+4+int(fingerprintLen)])
		cidBytes := key[len(prefix)+4+int(fingerprintLen):]
		cid := types.MustCidFromBytes(cidBytes)
		if f(fingerprint, cid) {
			break
		}
	}
}

func (k Keeper) IterateFingerprintCids(ctx sdk.Context, fingerprint string, f func(cid CID) bool) {
	prefix := GetFingerprintToCidKey(fingerprint)
	it := sdk.KVStorePrefixIterator(ctx.KVStore(k.storeKey), prefix)
	defer it.Close()
	for ; it.Valid(); it.Next() {
		cidBytes := it.Key()[len(prefix):]
		cid := types.MustCidFromBytes(cidBytes)
		if f(cid) {
			break
		}
	}
}

func (k Keeper) GetIscnIdVersion(ctx sdk.Context, iscnId IscnId) uint64 {
	key := GetIscnIdVersionKey(k.cdc, iscnId)
	bz := ctx.KVStore(k.storeKey).Get(key)
	if bz == nil {
		return 0
	}
	return binary.BigEndian.Uint64(bz)
}

func (k Keeper) SetIscnIdVersion(ctx sdk.Context, iscnId IscnId, count uint64) {
	key := GetIscnIdVersionKey(k.cdc, iscnId)
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, count)
	ctx.KVStore(k.storeKey).Set(key, bz)
}

func (k Keeper) SetIscnIdOwner(ctx sdk.Context, iscnId IscnId, owner sdk.AccAddress) {
	key := GetIscnIdOwnerKey(k.cdc, iscnId)
	ctx.KVStore(k.storeKey).Set(key, []byte(owner))
}

func (k Keeper) GetIscnIdOwner(ctx sdk.Context, iscnId IscnId) sdk.AccAddress {
	key := GetIscnIdOwnerKey(k.cdc, iscnId)
	bz := ctx.KVStore(k.storeKey).Get(key)
	if bz == nil {
		return nil
	}
	return sdk.AccAddress(bz)
}

func (k Keeper) DeductFeeForIscn(ctx sdk.Context, feePayer sdk.AccAddress, data []byte) error {
	acc := k.accountKeeper.GetAccount(ctx, feePayer)
	if acc == nil {
		return fmt.Errorf("account %s for deducting fee", feePayer.String())
	}
	feePerByte := k.GetParams(ctx).FeePerByte
	feeAmount := feePerByte.Amount.MulInt64(int64(len(data)))
	fees := sdk.NewCoins(sdk.NewCoin(feePerByte.Denom, feeAmount.Ceil().RoundInt()))
	if !fees.IsZero() {
		err := ante.DeductFees(k.bankKeeper, ctx, acc, fees)
		if err != nil {
			return err
		}
	}
	return nil
}

func (k Keeper) AddIscnRecord(
	ctx sdk.Context, id IscnId, owner sdk.AccAddress, record []byte, fingerprints []string,
) (*CID, error) {
	cid := types.ComputeRecordCid(record)
	if k.GetCidBlock(ctx, cid) != nil {
		return nil, fmt.Errorf("CID %s already exists", cid.String())
	}
	err := k.DeductFeeForIscn(ctx, owner, record)
	if err != nil {
		return nil, err
	}
	k.SetCidBlock(ctx, cid, record)
	k.SetCidIscnId(ctx, cid, id)
	k.SetIscnIdCid(ctx, id, cid)
	k.SetIscnIdVersion(ctx, id, id.Version)
	k.SetIscnIdOwner(ctx, id, owner)
	event := sdk.NewEvent(
		types.EventTypeIscnRecord,
		sdk.NewAttribute(types.AttributeKeyIscnId, id.String()),
		sdk.NewAttribute(types.AttributeKeyIscnIdPrefix, id.Prefix()),
		sdk.NewAttribute(types.AttributeKeyIscnOwner, owner.String()),
		sdk.NewAttribute(types.AttributeKeyIscnRecordIpld, cid.String()),
	)
	for _, fingerprint := range fingerprints {
		k.AddFingerprintCid(ctx, fingerprint, cid)
		event.AppendAttributes(sdk.NewAttribute(types.AttributeKeyIscnContentFingerprint, fingerprint))
	}
	ctx.EventManager().EmitEvent(event)
	return &cid, nil
}
