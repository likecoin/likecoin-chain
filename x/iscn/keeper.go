package iscn

import (
	"encoding/binary"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/tendermint/tendermint/crypto/tmhash"
)

const (
	DefaultParamspace = ModuleName
)

// TODO: move to individual file
type AccountKeeper interface {
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) auth.Account
}

type Keeper struct {
	storeKey      sdk.StoreKey
	cdc           *codec.Codec
	paramstore    params.Subspace
	codespace     sdk.CodespaceType
	accountKeeper AccountKeeper
	supplyKeeper  authTypes.SupplyKeeper
}

func NewKeeper(
	cdc *codec.Codec, key sdk.StoreKey, accountKeeper AccountKeeper,
	supplyKeeper authTypes.SupplyKeeper, paramstore params.Subspace,
	codespace sdk.CodespaceType,
) Keeper {
	return Keeper{
		storeKey:      key,
		cdc:           cdc,
		paramstore:    paramstore.WithKeyTable(ParamKeyTable()),
		codespace:     codespace,
		accountKeeper: accountKeeper,
		supplyKeeper:  supplyKeeper,
	}
}

func (keeper Keeper) Codespace() sdk.CodespaceType {
	return keeper.codespace
}

func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&Params{})
}

func (k Keeper) FeePerByte(ctx sdk.Context) (res sdk.DecCoin) {
	k.paramstore.Get(ctx, KeyFeePerByte, &res)
	return
}

func (k Keeper) GetParams(ctx sdk.Context) Params {
	return Params{
		FeePerByte: k.FeePerByte(ctx),
	}
}

func (k Keeper) SetParams(ctx sdk.Context, params Params) {
	k.paramstore.SetParamSet(ctx, &params)
}

func (k Keeper) SetIscnRecord(ctx sdk.Context, iscnId []byte, record *IscnRecord) {
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(record) // TODO: cbor?
	key := GetIscnRecordKey(iscnId)
	ctx.KVStore(k.storeKey).Set(key, bz)
}

func (k Keeper) GetIscnRecord(ctx sdk.Context, iscnId []byte) *IscnRecord {
	key := GetIscnRecordKey(iscnId)
	bz := ctx.KVStore(k.storeKey).Get(key)
	if bz == nil {
		return nil
	}
	record := IscnRecord{}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &record)
	return &record
}

func (k Keeper) IterateIscnRecords(ctx sdk.Context, f func(iscnId []byte, iscnRecord *IscnRecord) bool) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, IscnRecordKey)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		id := iterator.Key()[len(IscnRecordKey):]
		record := IscnRecord{}
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &record)
		stop := f(id, &record)
		if stop {
			break
		}
	}
}

func (k Keeper) AddIscnRecord(ctx sdk.Context, feePayer sdk.AccAddress, record *IscnRecord) (iscnId []byte, err sdk.Error) {
	acc := k.accountKeeper.GetAccount(ctx, feePayer)
	if acc == nil {
		return nil, nil // TODO: error
	}
	// TODO: checkings (in handler?)
	feePerByte := k.GetParams(ctx).FeePerByte
	feeAmount := feePerByte.Amount.MulInt64(int64(len(ctx.TxBytes())))
	fees := sdk.NewCoins(sdk.NewCoin(feePerByte.Denom, feeAmount.Ceil().RoundInt()))
	result := auth.DeductFees(k.supplyKeeper, ctx, acc, fees)
	if !result.IsOK() {
		return nil, nil // TODO: error
	}
	hasher := tmhash.New()
	hasher.Write(ctx.BlockHeader().LastBlockId.Hash)
	// TODO: counter for number of CreateIscn txs to prevent collision
	iscnCount := uint64(0)
	binary.Write(hasher, binary.BigEndian, iscnCount)
	id := hasher.Sum(nil)
	k.SetIscnRecord(ctx, id, record)
	return id, nil
}
