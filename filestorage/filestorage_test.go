package filestorage

import (
	"testing"
	"bytes"
	"fmt"
)

func TestWriteData(t *testing.T) {
	fileStorage, err := NewFileStorage("test.dat", ReadAndWrite, 0)
	if err != nil {
		t.Error(err)
	}
	data := []byte("hello world")
	location := DataLocation{0, uint32(len(data))}
	fileStorage.Write(data, location)
	buf, err := fileStorage.Read(location)
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(buf, data) {
		t.Error("something is wrong")
	}
	fmt.Println("data: ", data, buf)
}
