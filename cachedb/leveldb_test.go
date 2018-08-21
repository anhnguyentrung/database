package cachedb

import (
	"encoding/binary"
	"github.com/syndtr/goleveldb/leveldb/util"
	"testing"
)

func TestIterator(t *testing.T) {
	fileName := "test.db"
	db, err :=  NewLevelDB(fileName)
	if err != nil {
		panic(err)
	}
	for i := 1; i < 5; i++ {
		db.Put(encode(i), encode(i), false)
	}
	results := allValues(db)
	t.Log(results)
	db.Close()
}

func encode(i int) []byte {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(i))
	return buf
}

func decode(buf []byte) int {
	return int(binary.BigEndian.Uint64(buf))
}

func allValues(db *LevelDB) []int {
	result := make([]int, 0)
	slice := &util.Range{}
	slice.Start = nil
	slice.Limit = nil
	for itr := db.Iterator(slice); itr.Valid(); itr.Next() {
		result = append(result, decode(itr.Key()))
	}
	return result
}