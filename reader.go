package patchwork

import "io"

const (
	whencesStart   = 0
	whencesCurrent = 1
	whencesEnd     = 2
)

type reader struct {
	rat    io.ReaderAt
	offset int64
	size   int64
}

func NewReader(r io.ReaderAt, size, offset int64) *reader {
	return &reader{
		rat:    r,
		size:   size,
		offset: offset,
	}
}

func (r *reader) Read(p []byte) (n int, err error) {
	buf := p

	if r.offset >= r.size {
		return 0, io.EOF
	}

	// ReadAt is stricter than Read, fails if p > size available
	if r.offset+int64(len(p)) > r.size {
		buf = p[r.offset:r.size]
		err = io.EOF
	}

	n, err = r.rat.ReadAt(buf, r.offset)
	r.offset += int64(n)
	return
}

func (r *reader) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	case whencesStart:
		r.offset = offset
	case whencesCurrent:
		r.offset += offset
	case whencesEnd:
		r.offset = r.size - offset
	}
	return r.offset, nil
}
