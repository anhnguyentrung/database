package cachedb

type Iterator interface {
	Next() bool
	Key() []byte
	Value() []byte
	Release()
}

