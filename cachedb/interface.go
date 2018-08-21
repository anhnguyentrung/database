package cachedb

type Iterator interface {
	Valid() bool
	Next()
	Key() []byte
	Value() []byte
	Release()
}

