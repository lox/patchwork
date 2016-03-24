package patchwork

import (
	"io"
	"io/ioutil"
	"os"
)

// BufferAt is a buffer which supports io.ReaderAt and io.WriterAt
type Buffer interface {
	io.ReaderAt
	io.WriterAt

	// The capacity of the buffer
	Cap() int64
}

type File struct {
	*os.File
	RemoveOnClose bool
	Size          int64
}

var _ Buffer = &File{}

func NewFileBuffer(size int64) (*File, error) {
	f, err := ioutil.TempFile("", "patchwork")
	if err != nil {
		return nil, err
	}
	if err := os.Truncate(f.Name(), size); err != nil {
		return nil, err
	}
	return &File{File: f, RemoveOnClose: true, Size: size}, nil
}

func NewFileBufferString(s string) (*File, error) {
	f, err := NewFileBuffer(int64(len(s)))
	if err != nil {
		return nil, err
	}

	_, err = f.WriteAt([]byte(s), 0)
	return f, err
}

func (f *File) Cap() int64 {
	return f.Size
}

func (f *File) Close() error {
	err := f.File.Close()
	if err != nil {
		return err
	}
	if f.RemoveOnClose {
		err = os.Remove(f.File.Name())
	}
	return err
}
