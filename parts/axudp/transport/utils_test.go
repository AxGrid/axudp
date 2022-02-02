package transport

import (
	"bytes"
	"testing"
)

func BenchmarkAppendBytes(b *testing.B) {
	b1 := make([]byte, 2)
	b2 := make([]byte, 1024)
	for i := 0; i < b.N; i++ {
		_ = append(b1, b2...)
	}
}

func BenchmarkJoinBytes(b *testing.B) {
	b1 := make([]byte, 2)
	b2 := make([]byte, 1024)
	for i := 0; i < b.N; i++ {
		_ = bytes.Join([][]byte{b1, b2}, nil)
	}
}

func BenchmarkUtilsJoinBytes(b *testing.B) {
	b1 := make([]byte, 2)
	b2 := make([]byte, 1024)
	for i := 0; i < b.N; i++ {
		bytesJoin(b1, b2)
	}
}

func BenchmarkAddSize(b *testing.B) {
	b2 := make([]byte, 1024)
	for i := 0; i < b.N; i++ {
		addSize(b2)
	}
}
