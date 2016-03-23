package patchwork

import "io"

type writer struct {
	offset int64
	wat    io.WriterAt
}

func NewWriter(w io.WriterAt, offset int64) io.Writer {
	return &writer{offset, w}
}

func (w *writer) Write(b []byte) (n int, err error) {
	n, err = w.wat.WriteAt(b, w.offset)
	w.offset += int64(n)
	return
}
