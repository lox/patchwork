package patchwork

import (
	"bytes"
	"io"
	"io/ioutil"
	"math/rand"
	"strings"
	"sync"
	"testing"
)

var (
	rng *rand.Rand
)

func init() {
	rng = rand.New(rand.NewSource(24))
}

func TestWritingIntoPatchwork(t *testing.T) {
	data := strings.Repeat("llamas", 10)

	buf, err := NewFileBufferString(data)
	if err != nil {
		t.Fatal(err)
	}

	pw, err := New(buf)
	if err != nil {
		t.Fatal(err)
	}
	defer pw.Close()

	if _, err = randCopy(pw, bytes.NewBufferString(data)); err != nil && err != io.EOF {
		t.Fatal(err)
	}

	result, err := ioutil.ReadAll(pw.Reader())
	if err != nil {
		t.Fatal(err)
	}

	if string(result) != data {
		t.Fatalf("Expected %q, got %q", data, result)
	}
}

func TestReadingFromPatchwork(t *testing.T) {
	data := strings.Repeat("llamas", 10)

	buf, err := NewFileBufferString(data)
	if err != nil {
		t.Fatal(err)
	}

	pw, err := New(buf)
	if err != nil {
		t.Fatal(err)
	}
	defer pw.Close()

	if _, err := pw.WriteAt([]byte(data), 0); err != nil {
		t.Fatal(err)
	}

	result := make([]byte, len(data))
	if _, err = pw.ReadAt(result, 0); err != nil {
		t.Fatal(err)
	}

	if string(result) != string(data) {
		t.Fatalf("Expected %q, got %q", data, result)
	}
}

func TestBlockingReadsFromPatchwork(t *testing.T) {
	data := strings.Repeat("llamas", 6000)

	buf, err := NewFileBufferString(data)
	if err != nil {
		t.Fatal(err)
	}

	pw, err := New(buf)
	if err != nil {
		t.Fatal(err)
	}
	defer pw.Close()

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		_, err := randCopy(pw, strings.NewReader(data))
		if err != nil {
			t.Fatal(err)
		}
		wg.Done()
	}()

	var reads = []struct {
		from, to int64
	}{
		{0, 100},
		{100, 5000},
		{0, 6000},
	}

	for _, r := range reads {
		expected := data[r.from:r.to]

		result := make([]byte, len(expected))
		if _, err = pw.ReadAt(result, r.from); err != nil {
			t.Fatal(err)
		}

		if string(result) != string(expected) {
			t.Fatalf("Expected %q, got %q", expected, result)
		}
	}

	wg.Wait()
}

func TestConcurrentStreamingReads(t *testing.T) {
	data := strings.Repeat("llamas", 6000)

	buf, err := NewFileBufferString(data)
	if err != nil {
		t.Fatal(err)
	}

	pw, err := New(buf)
	if err != nil {
		t.Fatal(err)
	}
	defer pw.Close()

	var wg sync.WaitGroup

	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func() {
			r := pw.Reader()
			result := make([]byte, len(data))

			_, err := io.ReadFull(r, result)
			if err != nil {
				t.Fatal(err)
			}

			if string(result) != string(data) {
				t.Fatalf("Expected %q, got %q", data, result)
			}

			wg.Done()
		}()
	}

	go func() {
		if _, err = io.Copy(pw.Writer(), strings.NewReader(data)); err != nil && err != io.EOF {
			t.Fatal(err)
		}
	}()

	wg.Wait()
}

func randBuffer() []byte {
	var blockSize = rng.Intn(1024)

	if blockSize < 64 {
		blockSize = 64
	}

	return make([]byte, blockSize)
}

func randCopy(w io.WriterAt, r io.Reader) (int64, error) {
	var off, written int64

	for {
		buf := randBuffer()
		n, err := r.Read(buf)
		if n > 0 {
			n, err = w.WriteAt(buf[:n], off)
			written += int64(n)
		}
		if err != nil {
			break
		}
		off += int64(n)
	}

	return written, nil
}
