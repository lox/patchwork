package patchwork

import (
	"bytes"
	"fmt"
	"io"
	"sync"
)

// Patchwork glues together io.ReaderAt and io.WriterAt to allow for asynchronous reading and
// writing of a resource to gradually coalesce into a whole. Reads for missing portions block,
// writes aren't synchronized, so it's up to the caller to ensure they aren't overlapping.
type Patchwork struct {
	cond       *sync.Cond
	buf        Buffer
	size       int64
	written    *intervalSet
	Downloader Downloader
}

// New creates a Patchwork with the given Buffer
func New(buf Buffer) *Patchwork {
	return &Patchwork{
		buf:     buf,
		size:    int64(buf.Cap()),
		cond:    sync.NewCond(&sync.Mutex{}),
		written: &intervalSet{},
	}
}

// WriteAt fulfils io.WriterAt and allows concurrent writes
func (p *Patchwork) WriteAt(data []byte, off int64) (int, error) {
	p.cond.L.Lock()
	defer p.cond.L.Unlock()

	written, err := io.Copy(NewWriter(p.buf, off), bytes.NewReader(data))
	if written > 0 {
		p.written.Add(interval{off, off + int64(written)})
		p.cond.Broadcast()
	}

	return int(written), err
}

// ReadAt fulfils io.ReaderAt and allows concurrent reads
func (p *Patchwork) ReadAt(b []byte, off int64) (n int, err error) {
	if off+int64(len(b)) > p.size {
		return 0, fmt.Errorf("ReadAt goes past end of patchwork")
	}

	if !p.written.Contains(interval{off, off + int64(len(b))}) {
		p.cond.L.Lock()
		for !p.written.Contains(interval{off, off + int64(len(b))}) {
			if p.Downloader != nil {
				if err = p.Downloader.Download(off, off+int64(len(b))); err != nil {
					return n, err
				}
			}
			p.cond.Wait()
		}
		p.cond.L.Unlock()
	}

	return p.buf.ReadAt(b, off)
}

// Closes the patchwork and if the buffer is an io.Closer, closes it too
func (p *Patchwork) Close() error {
	if closer, ok := p.buf.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}

func (p *Patchwork) Reader() io.Reader {
	return io.NewSectionReader(p, 0, p.size)
}

func (p *Patchwork) Writer() io.Writer {
	return NewWriter(p, 0)
}

type Downloader interface {
	Download(from, to int64) error
}

type DownloaderFunc func(from, to int64) error

func (df DownloaderFunc) Download(from, to int64) error {
	return df(from, to)
}
