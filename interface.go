package database

type Iterator interface {
	Next()
	Key() []byte
	Value() []byte
	Release()
}
