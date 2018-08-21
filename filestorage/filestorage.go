package filestorage

import (
	"os"
	"fmt"
	"syscall"
)

type FileAccessMode uint8

const (
	Read FileAccessMode = iota
	ReadAndWrite
)

type FileStorage struct {
	file *os.File
	data []byte
}

func NewFileStorage(path string, accessMode FileAccessMode) (*FileStorage, error) {
	flag := os.O_RDONLY
	if accessMode == ReadAndWrite {
		flag = os.O_CREATE | os.O_RDWR
	}
	f, err := os.OpenFile(path, flag, 0)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	fileStats, err := f.Stat()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	size := fileStats.Size()
	prot := syscall.PROT_READ
	if accessMode == ReadAndWrite {
		prot = syscall.PROT_WRITE | syscall.PROT_READ
	}
	data, err := syscall.Mmap(int(f.Fd()), 0, int(size), prot, syscall.MAP_SHARED)
	fileStorage := &FileStorage{
		file: f,
		data: data,
	}
	return fileStorage, nil
}

func (fileStorage *FileStorage) Seek

func (fileStorage *FileStorage) Close() {
	err := fileStorage.file.Close()
	if err != nil {
		panic(err)
	}
	if fileStorage.data == nil {
		return
	}
	err = syscall.Munmap(fileStorage.data)
	if err != nil {
		panic(err)
	}
	fileStorage.data = nil
}


