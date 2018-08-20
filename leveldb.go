package database

import (
	"github.com/syndtr/goleveldb/leveldb"
	"fmt"
	"github.com/syndtr/goleveldb/leveldb/opt"
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

