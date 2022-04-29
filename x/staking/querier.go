package staking

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"

	dbm "github.com/tendermint/tm-db"

	"github.com/cosmos/cosmos-sdk/codec"
	dbstore "github.com/cosmos/cosmos-sdk/store/dbadapter"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"

	"github.com/cosmos/cosmos-sdk/x/staking/keeper"
	"github.com/cosmos/cosmos-sdk/x/staking/types"
)

var ValidatorDelegationsIndexPrefix = GetTableKey("ValidatorDelegationsIndexPrefix")
var ValidatorDelegationsIndexHeightKey = GetTableKey("ValidatorDelegationsIndexHeightKey")

func GetTableKey(name string) []byte {
	hash := sha256.Sum256([]byte(name))
	return hash[:8]
}

// `append` has concurrency issue so can't be used here
func getIndexValidatorPrefix(valAddr sdk.ValAddress) []byte {
	buf := &bytes.Buffer{}
	buf.Write(ValidatorDelegationsIndexPrefix)
	buf.Write(address.MustLengthPrefix(valAddr))
	return buf.Bytes()
}

// `append` has concurrency issue so can't be used here
func getIndexKey(valAddr sdk.ValAddress, delAddr sdk.AccAddress) []byte {
	buf := &bytes.Buffer{}
	buf.Write(ValidatorDelegationsIndexPrefix)
	buf.Write(address.MustLengthPrefix(valAddr))
	buf.Write(address.MustLengthPrefix(delAddr))
	return buf.Bytes()
}

func encodeUint64(n uint64) []byte {
	heightBz := make([]byte, 8)
	binary.BigEndian.PutUint64(heightBz, n)
	return heightBz
}

type Querier struct {
	types.QueryServer
	keeper     *keeper.Keeper
	cdc        codec.BinaryCodec
	indexingDB dbm.DB
	batch      dbm.Batch
	readStore  dbstore.Store
}

func NewQuerier(k *keeper.Keeper, cdc codec.BinaryCodec, indexingDB dbm.DB) *Querier {
	readStore := dbstore.Store{DB: indexingDB}
	return &Querier{
		QueryServer: keeper.Querier{Keeper: *k},
		keeper:      k,
		cdc:         cdc,
		indexingDB:  indexingDB,
		batch:       indexingDB.NewBatch(),
		readStore:   readStore,
	}
}

func (q *Querier) batchSet(key, value []byte) {
	err := q.batch.Set(key, value)
	if err != nil {
		panic(err)
	}
}

func (q *Querier) batchDelete(key []byte) {
	err := q.batch.Delete(key)
	if err != nil {
		panic(err)
	}
}

func (q *Querier) GetHeight() uint64 {
	bz := q.readStore.Get(ValidatorDelegationsIndexHeightKey)
	if bz == nil {
		return 0
	}
	return binary.BigEndian.Uint64(bz)
}

func (q *Querier) setHeight(newHeight uint64) {
	q.batchSet(ValidatorDelegationsIndexHeightKey, encodeUint64(newHeight))
}

func (q *Querier) clearIndex() {
	it, err := q.indexingDB.Iterator(nil, nil)
	if err != nil {
		panic(err)
	}
	defer it.Close()
	for ; it.Valid(); it.Next() {
		err := q.batch.Delete(it.Key())
		if err != nil {
			panic(err)
		}
	}
}

func (q *Querier) BuildIndex(ctx sdk.Context) {
	blockHeight := ctx.BlockHeight()
	ctx.Logger().Debug("Rebuilding index for ValidatorDelegations", "block_height", blockHeight)
	q.clearIndex()
	// TODO: this could be slow (30s up), need to change this to async
	q.keeper.IterateAllDelegations(ctx, func(delegation types.Delegation) bool {
		key := getIndexKey(delegation.GetValidatorAddr(), delegation.GetDelegatorAddr())
		q.batchSet(key, types.MustMarshalDelegation(q.cdc, delegation))
		return false
	})
	q.batchSet(ValidatorDelegationsIndexHeightKey, encodeUint64(uint64(ctx.BlockHeight())))
	q.flushBatch()
}

func (q *Querier) BeginWriteIndex(ctx sdk.Context) {
	currentHeight := q.GetHeight()
	blockHeight := uint64(ctx.BlockHeight())
	ctx.Logger().Debug(
		"Beginning write index",
		"index_height", currentHeight,
		"block_height", blockHeight,
	)
	if currentHeight == 0 || blockHeight != currentHeight+1 {
		q.BuildIndex(ctx)
	}
}

func (q *Querier) RemoveIndex(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) {
	ctx.Logger().Debug(
		"Removing index",
		"block_height", ctx.BlockHeight(),
		"delegator", delAddr.String(),
		"validator", valAddr.String(),
	)
	key := getIndexKey(valAddr, delAddr)
	q.batchDelete(key)
}

func (q *Querier) UpdateIndex(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) {
	ctx.Logger().Debug(
		"Updating index",
		"block_height", ctx.BlockHeight(),
		"delegator", delAddr.String(),
		"validator", valAddr.String(),
	)
	key := getIndexKey(valAddr, delAddr)
	delegation, found := q.keeper.GetDelegation(ctx, delAddr, valAddr)
	if !found {
		ctx.Logger().Debug("Delegation not found, deleting index")
		q.batchDelete(key)
	} else {
		ctx.Logger().Debug("Delegation found, updating index", "delegation", delegation)
		q.batchSet(key, types.MustMarshalDelegation(q.cdc, delegation))
	}
}

func (q *Querier) flushBatch() {
	b := q.batch
	q.batch = q.indexingDB.NewBatch()
	err := b.Write()
	if err != nil {
		panic(err)
	}
	err = b.Close()
	if err != nil {
		panic(err)
	}
}

func (q *Querier) CommitWriteIndex(ctx sdk.Context) {
	blockHeight := uint64(ctx.BlockHeight())
	ctx.Logger().Debug(
		"Committing index writes",
		"block_height", blockHeight,
	)
	q.setHeight(blockHeight)
	q.flushBatch()
}
