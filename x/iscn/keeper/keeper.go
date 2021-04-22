package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	prefixstore "github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
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

func (k Keeper) prefixStore(ctx sdk.Context, prefix []byte) prefixstore.Store {
	return prefixstore.NewStore(ctx.KVStore(k.storeKey), prefix)
}

func (k Keeper) GetSequenceCount(ctx sdk.Context) uint64 {
	bz := ctx.KVStore(k.storeKey).Get(SequenceCountKey)
	return types.DecodeUint64(bz)
}

func (k Keeper) setSequenceCount(ctx sdk.Context, seq uint64) {
	bz := types.EncodeUint64(seq)
	ctx.KVStore(k.storeKey).Set(SequenceCountKey, bz)
}

func (k Keeper) GetStoreRecord(ctx sdk.Context, seq uint64) *StoreRecord {
	seqBytes := types.EncodeUint64(seq)
	recordBytes := k.prefixStore(ctx, SequenceToStoreRecordPrefix).Get(seqBytes)
	if recordBytes == nil {
		return nil
	}
	record := k.MustUnmarshalStoreRecord(recordBytes)
	return &record
}

func (k Keeper) AddStoreRecord(ctx sdk.Context, record StoreRecord) (seq uint64) {
	seq = k.GetSequenceCount(ctx)
	seq += 1
	k.setSequenceCount(ctx, seq)
	seqBytes := types.EncodeUint64(seq)
	recordBytes := k.MustMarshalStoreRecord(&record)
	k.prefixStore(ctx, SequenceToStoreRecordPrefix).Set(seqBytes, recordBytes)
	iscnIdBytes := k.MustMarshalIscnId(record.IscnId)
	k.prefixStore(ctx, IscnIdToSequencePrefix).Set(iscnIdBytes, seqBytes)
	k.prefixStore(ctx, CidToSequencePrefix).Set(record.CidBytes, seqBytes)
	return seq
}

func (k Keeper) IterateStoreRecords(ctx sdk.Context, f func(seq uint64, record StoreRecord) bool) {
	it := k.prefixStore(ctx, SequenceToStoreRecordPrefix).Iterator(nil, nil)
	defer it.Close()
	for ; it.Valid(); it.Next() {
		seq := types.DecodeUint64(it.Key())
		record := k.MustUnmarshalStoreRecord(it.Value())
		if f(seq, record) {
			break
		}
	}
}

func (k Keeper) GetIscnIdSequence(ctx sdk.Context, iscnId IscnId) uint64 {
	iscnIdBytes := k.MustMarshalIscnId(iscnId)
	seqBytes := k.prefixStore(ctx, IscnIdToSequencePrefix).Get(iscnIdBytes)
	return types.DecodeUint64(seqBytes)
}

func (k Keeper) GetCidSequence(ctx sdk.Context, cid CID) uint64 {
	seqBytes := k.prefixStore(ctx, CidToSequencePrefix).Get(cid.Bytes())
	return types.DecodeUint64(seqBytes)
}

func (k Keeper) AddFingerprintSequence(ctx sdk.Context, fingerprint string, seq uint64) {
	key := types.GetFingerprintSequenceKey(fingerprint, seq)
	k.prefixStore(ctx, FingerprintSequencePrefix).Set(key, []byte{0x01})
}

func (k Keeper) IterateAllFingerprints(ctx sdk.Context, f func(fingerprint string, seq uint64) bool) {
	it := k.prefixStore(ctx, FingerprintSequencePrefix).Iterator(nil, nil)
	defer it.Close()
	for ; it.Valid(); it.Next() {
		fingerprint, seq := types.ParseFingerprintSequenceBytes(it.Key())
		if f(fingerprint, seq) {
			break
		}
	}
}

func (k Keeper) IterateFingerprintSequencesWithStartingSequence(ctx sdk.Context, fingerprint string, seq uint64, f func(seq uint64) bool) {
	prefix := types.GetFingerprintPrefix(fingerprint)
	fromKey := types.EncodeUint64(seq)
	it := k.prefixStore(ctx, prefix).Iterator(fromKey, nil)
	defer it.Close()
	for ; it.Valid(); it.Next() {
		seq := types.DecodeUint64(it.Key())
		if f(seq) {
			break
		}
	}
}

func (k Keeper) IterateFingerprintSequences(ctx sdk.Context, fingerprint string, f func(seq uint64) bool) {
	k.IterateFingerprintSequencesWithStartingSequence(ctx, fingerprint, 0, f)
}

func (k Keeper) HasFingerprintSequence(ctx sdk.Context, fingerprint string, seq uint64) bool {
	key := types.GetFingerprintSequenceKey(fingerprint, seq)
	return k.prefixStore(ctx, FingerprintSequencePrefix).Has(key)
}

func (k Keeper) GetTracingIdRecord(ctx sdk.Context, iscnId IscnId) *TracingIdRecord {
	key := k.MustMarshalTracingId(iscnId)
	bz := k.prefixStore(ctx, TracingIdRecordPrefix).Get(key)
	if bz == nil {
		return nil
	}
	record := k.MustUnmarshalTracingIdRecord(bz)
	return &record
}

func (k Keeper) SetTracingIdRecord(ctx sdk.Context, iscnId IscnId, record *TracingIdRecord) {
	key := k.MustMarshalTracingId(iscnId)
	recordBytes := k.MustMarshalTracingIdRecord(record)
	k.prefixStore(ctx, TracingIdRecordPrefix).Set(key, recordBytes)
}

func (k Keeper) IterateTracingIdRecords(ctx sdk.Context, f func(iscnId IscnId, tracingIdRecord TracingIdRecord) bool) {
	it := k.prefixStore(ctx, TracingIdRecordPrefix).Iterator(nil, nil)
	defer it.Close()
	for ; it.Valid(); it.Next() {
		tracingId := k.MustUnmarshalIscnId(it.Key())
		record := k.MustUnmarshalTracingIdRecord(it.Value())
		if f(tracingId, record) {
			break
		}
	}
}

func (k Keeper) IterateIscnIds(ctx sdk.Context, f func(iscnId IscnId, tracingIdRecord TracingIdRecord) bool) {
	k.IterateTracingIdRecords(ctx, func(iscnId IscnId, tracingIdRecord TracingIdRecord) bool {
		for version := uint64(1); version <= tracingIdRecord.LatestVersion; version++ {
			iscnId.Version = version
			if f(iscnId, tracingIdRecord) {
				return true
			}
		}
		return false
	})
}

func (k Keeper) DeductFeeForIscn(ctx sdk.Context, feePayer sdk.AccAddress, data []byte) error {
	acc := k.accountKeeper.GetAccount(ctx, feePayer)
	if acc == nil {
		return fmt.Errorf("cannot find account %s for deducting fee", feePayer.String())
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
	ctx sdk.Context, iscnId IscnId, owner sdk.AccAddress, data []byte, fingerprints []string,
) (*CID, error) {
	if k.GetIscnIdSequence(ctx, iscnId) != 0 {
		return nil, sdkerrors.Wrapf(types.ErrReusingIscnId, "%s", iscnId.String())
	}
	cid := types.ComputeDataCid(data)
	if k.GetCidSequence(ctx, cid) != 0 {
		return nil, sdkerrors.Wrapf(types.ErrCidAlreadyExist, "%s", cid.String())
	}
	tracingIdRecord := k.GetTracingIdRecord(ctx, iscnId)
	if tracingIdRecord == nil {
		if iscnId.Version != 1 {
			return nil, sdkerrors.Wrapf(types.ErrInvalidIscnVersion, "expected version: 1")
		}
	} else {
		if iscnId.Version != tracingIdRecord.LatestVersion+1 {
			return nil, sdkerrors.Wrapf(types.ErrInvalidIscnVersion, "expected version: %d", tracingIdRecord.LatestVersion+1)
		}
		expectedOwner := tracingIdRecord.OwnerAddress()
		if !expectedOwner.Equals(owner) {
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "expected owner: %s", owner.String())
		}
	}
	err := k.DeductFeeForIscn(ctx, owner, data)
	if err != nil {
		return nil, sdkerrors.Wrapf(types.ErrDeductIscnFee, "%s", err.Error())
	}
	record := StoreRecord{
		IscnId:   iscnId,
		CidBytes: cid.Bytes(),
		Data:     data,
	}
	seq := k.AddStoreRecord(ctx, record)
	k.SetTracingIdRecord(ctx, iscnId, &TracingIdRecord{
		OwnerAddressBytes: owner.Bytes(),
		LatestVersion:     iscnId.Version,
	})
	event := sdk.NewEvent(
		types.EventTypeIscnRecord,
		sdk.NewAttribute(types.AttributeKeyIscnId, iscnId.String()),
		sdk.NewAttribute(types.AttributeKeyIscnIdPrefix, iscnId.Prefix()),
		sdk.NewAttribute(types.AttributeKeyIscnOwner, owner.String()),
		sdk.NewAttribute(types.AttributeKeyIscnRecordIpld, cid.String()),
	)
	for _, fingerprint := range fingerprints {
		k.AddFingerprintSequence(ctx, fingerprint, seq)
		event.AppendAttributes(sdk.NewAttribute(types.AttributeKeyIscnContentFingerprint, fingerprint))
	}
	ctx.EventManager().EmitEvent(event)
	return &cid, nil
}
