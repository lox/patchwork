Patchwork
=========

[![GoDoc](https://godoc.org/github.com/lox/patchwork?status.svg)](https://godoc.org/github.com/lox/patchwork)
[![Software License](https://img.shields.io/badge/license-MIT-brightgreen.svg)](LICENSE.txt)

Golang provides asynchronous reading and writing interfaces in the form of [io.ReaderAt] and [io.WriterAt]. Patchwork provides a [Pipe-like] way to connect the two such that reading an unwritten range blocks until the write occurs.

Patchwork uses a buffer with a known maximum capacity to cache the partially written stream (hence the name patchwork), and also provides a [io.ReadSeeker] interface for it. This allows for use-cases like splitting an HTTP request into multiple range-based chunks and requesting them in parallel, but still using [http.ServeContent] to serve the result downstream with support for arbitrary range requests (see [Examples](#examples)).

Writes can either be co-ordinated outside of patchwork and written in via [WriteAt] or a [Downloader] can be used to trigger asynchronous download on-demand.

## Buffers

A [Buffer]() stores the remote stream as it's downloaded. It needs to provide asynchronous random-access ([ReaderWriterAt]). An implementation is provided that uses an underlying `os.File`, but anything that implements `BufferAt` from https://github.com/djherbis/buffer is also compatible.

## Examples

```go
buf, err := patchwork.NewFileBuffer(10)
if err != nil {
	log.Fatal(err)
}

pw := patchwork.New(buf)

// writes can happen asynchronously, in chunks
go func() {
	for i := 0; i < 10; i += 2 {
		pw.WriteAt([]byte("xy"), int64(i))
	}
}()

var wg sync.WaitGroup

// launch 10 concurrent readers all waiting for the full stream
for i := 0; i < 10; i++ {
	wg.Add(1)
	go func() {
		b, err := ioutil.ReadAll(pw.Reader())
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Full stream: %q", b)
		wg.Done()
	}()
}

// launch 10 concurrent readers, reading chunks as they are ready
for i := 0; i < 10; i += 2 {
	wg.Add(1)
	go func(offset int) {
		log.Printf("Reading at %d", offset)
		b := make([]byte, 2)
		_, err := pw.ReadAt(b, int64(offset))
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Chunk: %q", b)
		wg.Done()
	}(i)
}

wg.Wait()
```

## License

MIT Licensed (c) Lachlan Donald 2016

Inspire by https://github.com/djherbis/fscache.

[io.ReaderAt]: https://golang.org/pkg/io/#ReaderAt
[io.WriterAt]: https://golang.org/pkg/io/#WriterAt
[io.ReadSeeker]: https://golang.org/pkg/io/#ReadSeeker
[http.ServeContent]: https://golang.org/pkg/http/#ServeContent
[Downloader]: https://godoc.org/github.com/lox/patchwork#Downloader
[Buffer]: https://godoc.org/github.com/lox/patchwork#Buffer
[WriteAt]: https://godoc.org/github.com/lox/patchwork#Patchwork.WriteAt
[Pipe-link]: https://golang.org/pkg/io/#Pipe