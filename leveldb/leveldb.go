package leveldb

import (
	"github.com/syndtr/goleveldb/leveldb"
	"fmt"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/util"
	"bytes"
)

type LevelDB struct {
	db 		*leveldb.DB
	batch 	*leveldb.Batch
}

func NewLevelDB(path string) (levelDB *LevelDB, err error)  {
	levelDB = nil
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	batch := new(leveldb.Batch)
	levelDB = &LevelDB{
		db:		db,
		batch:	batch,
	}
	return
}

func (levelDB *LevelDB) Close() {
	err := levelDB.db.Close()
	if err != nil {
		fmt.Println(err)
	}
}

func (levelDB *LevelDB) Delete(key []byte, sync bool) {
	if len(key) == 0 {
		fmt.Println("length of key should be greater than 0")
		return
	}
	wo := &opt.WriteOptions{Sync: sync}
	err := levelDB.db.Delete(key, wo)
	if err != nil {
		panic(err)
	}
}

func (levelDB *LevelDB) Get(key []byte) (value []byte) {
	if len(key) == 0 {
		return nil
	}
	value, err := levelDB.db.Get(key, nil)
	if err != nil {
		if err == leveldb.ErrNotFound {
			value = nil
			return
		}
		panic(err)
	}
	return
}

func (levelDB *LevelDB) Has(key []byte) (ret bool) {
	ret = levelDB.Get(key) != nil
	return
}

func (levelDB *LevelDB) Put(key, value []byte, sync bool) {
	if len(key) == 0 || len(value) == 0 {
		fmt.Println("lengths of key and value should be greater than 0")
		return
	}
	wo := &opt.WriteOptions{Sync: sync}
	err := levelDB.db.Put(key, value, wo)
	if err != nil {
		panic(err)
	}
}

func (levelDB *LevelDB) DeleteBatch(key []byte) {
	levelDB.batch.Delete(key)
}

func (levelDB *LevelDB) PutBatch(key, value []byte) {
	levelDB.batch.Put(key, value)
}

func (levelDB *LevelDB) WriteBatch(sync bool) {
	wo := &opt.WriteOptions{Sync: sync}
	err := levelDB.db.Write(levelDB.batch, wo)
	if err != nil {
		panic(err)
	}
}

func (levelDB *LevelDB) Iterator(slice *util.Range) Iterator {
	itr := levelDB.db.NewIterator(slice, nil)
	return newLevelDBIterator(itr, slice)
}

type levelDBIterator struct {
	iterator iterator.Iterator
	slice *util.Range
	forwards bool
}

var _ Iterator = (*levelDBIterator)(nil)

func newLevelDBIterator(itr iterator.Iterator, slice *util.Range) *levelDBIterator {
	forwards := true
	start := slice.Start
	limit := slice.Limit
	if bytes.Compare(start, limit) > 0 {
		forwards = false
	}
	return &levelDBIterator{
		itr,
		slice,
		forwards,
	}
}

func (ldbIterator *levelDBIterator) Next() bool {
	if err := ldbIterator.iterator.Error(); err != nil {
		panic(err)
	}
	if !ldbIterator.iterator.Valid() {
		panic("invalid iterator")
	}
	if ldbIterator.forwards {
		return ldbIterator.iterator.Next()
	} else {
		return ldbIterator.iterator.Prev()
	}
}

func (ldbIterator *levelDBIterator) Key() []byte {
	if err := ldbIterator.iterator.Error(); err != nil {
		panic(err)
	}
	if !ldbIterator.iterator.Valid() {
		panic("invalid iterator")
	}
	return ldbIterator.iterator.Key()
}

func (ldbIterator *levelDBIterator) Value() []byte {
	if err := ldbIterator.iterator.Error(); err != nil {
		panic(err)
	}
	if !ldbIterator.iterator.Valid() {
		panic("invalid iterator")
	}
	return ldbIterator.iterator.Value()
}

func (ldbIterator *levelDBIterator) Release() {
	ldbIterator.iterator.Release()
}

