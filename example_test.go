package patchwork_test

import (
	"io/ioutil"
	"log"
	"sync"
	"testing"

	"github.com/lox/patchwork"
)

func TestReadmeExample(t *testing.T) {
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
			if len(b) != 10 {
				t.Fatalf("Expected 10 bytes, got %d bytes", len(b))
			}
			wg.Done()
		}()
	}

	// launch 10 concurrent readers, reading chunks as they are ready
	for i := 0; i < 10; i += 2 {
		wg.Add(1)
		go func(offset int) {
			b := make([]byte, 2)
			_, err := pw.ReadAt(b, int64(offset))
			if err != nil {
				log.Fatal(err)
			}
			wg.Done()
		}(i)
	}

	wg.Wait()
}
