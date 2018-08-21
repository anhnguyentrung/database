package main

import (
	"encoding/binary"
	"db/cachedb"
	"fmt"
	"github.com/syndtr/goleveldb/leveldb/util"
)

func main() {
	fileName := "test.db"
	db, err :=  cachedb.NewLevelDB(fileName)
	if err != nil {
		panic(err)
	}
	for i := 1; i < 5; i++ {
		db.Put(encode(i), encode(i), false)
	}
	fmt.Println(allValues(db))
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

func allValues(db *cachedb.LevelDB) []int {
	result := make([]int, 0)
	slice := &util.Range{}
	slice.Start = nil
	slice.Limit = nil
	for itr := db.Iterator(slice); itr.Next(); {
		result = append(result, decode(itr.Value()))
	}
	return result
}