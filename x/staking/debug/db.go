package debug

import (
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"
)

type DBDebugWrap struct {
	dbm.DB
	Logger log.Logger
}

func (d *DBDebugWrap) Set(key []byte, value []byte) error {
	return d.DB.Set(key, value)
}

func (d *DBDebugWrap) Delete(key []byte) error {
	return d.DB.Delete(key)
}

func (d *DBDebugWrap) Close() error {
	d.Logger.Debug("DB Close")
	return d.DB.Close()
}

func (d *DBDebugWrap) NewBatch() dbm.Batch {
	return &DBDebugBatch{
		Batch:  d.DB.NewBatch(),
		Logger: d.Logger,
	}
}

func (d *DBDebugWrap) Iterator(start, end []byte) (dbm.Iterator, error) {
	// d.Logger.Debug("DB Iterator", "start", start, "end", end)
	it, err := d.DB.Iterator(start, end)
	if err != nil {
		return nil, err
	}
	return &DBDebugIterator{
		Iterator: it,
		Logger:   d.Logger,
	}, nil
}

type DBDebugBatch struct {
	dbm.Batch
	Logger log.Logger
}

func (b *DBDebugBatch) Set(key, value []byte) error {
	// b.Logger.Debug("DB Batch Set", "key", key, "value", value)
	return b.Batch.Set(key, value)
}

func (b *DBDebugBatch) Delete(key []byte) error {
	// b.Logger.Debug("DB Batch Set", "key", key)
	return b.Batch.Delete(key)
}

func (b *DBDebugBatch) Write() error {
	// b.Logger.Debug("DB Batch Write")
	return b.Batch.Write()
}

func (b *DBDebugBatch) WriteSync() error {
	// b.Logger.Debug("DB Batch WriteSync")
	return b.Batch.WriteSync()
}

type DBDebugIterator struct {
	dbm.Iterator
	Logger log.Logger
}

func (it *DBDebugIterator) Key() []byte {
	key := it.Iterator.Key()
	// it.Logger.Debug("DB Iterator Key", "key", key)
	return key
}

func (it *DBDebugIterator) Value() []byte {
	value := it.Iterator.Value()
	// it.Logger.Debug("DB Iterator value", "value", value)
	return value
}
