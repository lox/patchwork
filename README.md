Patchwork
=========

[![GoDoc](https://godoc.org/github.com/lox/patchwork?status.svg)](https://godoc.org/github.com/lox/patchwork)
[![Software License](https://img.shields.io/badge/license-MIT-brightgreen.svg)](LICENSE.txt)

Golang provides asynchronous reading and writing interfaces in the form of [io.ReaderAt]() and [io.WriterAt](). Patchwork provides a [Pipe-like]() way to connect the two such that reading an unwritten range blocks until the write occurs.

Patchwork uses a buffer to cache the partially written stream (hence the name patchwork), and also provides a [io.ReadSeeker]() interface for it. This allows for use-cases like splitting an HTTP request into multiple range-based chunks and requesting them in parallel, but still using [http.ServeContent]() to serve the result downstream with support for arbitrary range requests (see [Examples]()).

Writes can either be co-ordinated outside of patchwork and written in via [WriteAt]() or a [Downloader]() can be used to trigger asynchronous download on-demand.

## Buffers

A [Buffer]() stores the remote stream as it's downloaded. It needs to provide asynchronous random-access ([ReaderWriterAt]). An implementation is provided that uses an underlying `os.File`, but anything that implements `BufferAt` from https://github.com/djherbis/buffer is also compatible.

## License

MIT Licensed (c) Lachlan Donald 2016

Inspire by https://github.com/djherbis/fscache.