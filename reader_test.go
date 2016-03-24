package patchwork

import "testing"

func TestReaderReadsToEOF(t *testing.T) {
	data, err := NewFileBufferString("the lovely llama loves lettuce")
	if err != nil {
		t.Fatal(err)
	}

	r := NewReader(data, data.Cap(), 0)
	buf := make([]byte, data.Cap())
	n, err := r.Read(buf)
	if err != nil {
		t.Fatal(err)
	}

	if int64(n) != data.Cap() {
		t.Fatalf("Was expecting %d, got %d", data.Cap(), n)
	}
}
