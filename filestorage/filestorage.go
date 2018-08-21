package filestorage

import (
	"os"
	"fmt"
	"syscall"
	"io"
)

type FileAccessMode uint8

const (
	Read FileAccessMode = iota
	ReadAndWrite
)

type DataLocation struct {
	offset uint32
	length uint32
}

type FileStorage struct {
	path		string
	file 		*os.File
	mmapData 	[]byte
}

func NewFileStorage(path string, accessMode FileAccessMode, dataSize int64) (*FileStorage, error) {
	flag := os.O_RDONLY
	if accessMode == ReadAndWrite {
		flag = os.O_CREATE | os.O_RDWR
	}
	f, err := os.OpenFile(path, flag, 0666)
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
		if dataSize > size {
			size = dataSize
			err = syscall.Ftruncate(int(f.Fd()), size)
			if err != nil {
				fmt.Println(err)
				return nil, err
			}
		}
	}
	data, err := syscall.Mmap(int(f.Fd()), 0, int(size), prot, syscall.MAP_SHARED)
	fileStorage := &FileStorage{
		path: path,
		file: 		f,
		mmapData: 	data,
	}
	return fileStorage, nil
}

func (fileStorage *FileStorage) Close() {
	err := fileStorage.file.Close()
	if err != nil {
		panic(err)
	}
	if fileStorage.mmapData == nil {
		return
	}
	err = syscall.Munmap(fileStorage.mmapData)
	if err != nil {
		panic(err)
	}
	fileStorage.mmapData = nil
}

func (fileStorage *FileStorage) Read(location DataLocation) ([]byte, error) {
	if fileStorage.mmapData == nil {
		return nil, fmt.Errorf("nil data")
	}
	if location.offset + location.length > fileStorage.size() {
		return nil, io.EOF
	}
	start := location.offset
	limit := location.offset + location.length
	return fileStorage.mmapData[start:limit], nil
}

func (fileStorage *FileStorage) Write(data []byte, location DataLocation) error {
	if uint32(len(data)) != location.length {
		return fmt.Errorf("wrong data")
	}
	if location.offset + location.length > fileStorage.size() {
		newSize := int64(location.offset + location.length)
		err := fileStorage.resize(newSize)
		if err != nil {
			return err
		}
	}
	start := location.offset
	limit := location.offset + location.length
	fmt.Println("start limit: ", start, limit)
	copy(fileStorage.mmapData[start:limit], data)
	return nil
}

func (fileStorage *FileStorage) resize(newSize int64) error {
	path := fileStorage.path
	fileStorage.Close()
	fileStorage, err := NewFileStorage(path, ReadAndWrite, newSize)
	return err
}

func (fileStorage *FileStorage) size() uint32 {
	return uint32(len(fileStorage.mmapData))
}


