// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package crc64nvme

import (
	"bytes"
	"fmt"
	"hash"
	"hash/crc64"
	"io"
	"math/rand"
	"testing"
)

var crc64Table = crc64.MakeTable(NVME)

func TestChecksum(t *testing.T) {
	sizes := []int{0, 1, 3, 7, 8, 9, 15, 17, 127, 128, 129, 255, 256, 257, 1e3, 1e4, 1e5, 1e6}
	for _, size := range sizes {
		t.Run(fmt.Sprintf("size=%d", size), func(t *testing.T) {
			rng := rand.New(rand.NewSource(int64(size)))
			data := make([]byte, size)
			rng.Read(data)
			ref := crc64.Checksum(data, crc64Table)
			got := Checksum(data)
			if got != ref {
				t.Errorf("got 0x%x, want 0x%x", got, ref)
			}
		})
	}
}

func TestHasher(t *testing.T) {
	sizes := []int{0, 1, 3, 7, 8, 9, 15, 17, 127, 128, 129, 255, 256, 257, 1e3, 1e4, 1e5, 1e6}
	for _, size := range sizes {
		t.Run(fmt.Sprintf("size=%d", size), func(t *testing.T) {
			rng := rand.New(rand.NewSource(int64(size)))
			data := make([]byte, size)
			rng.Read(data)
			ref := crc64.Checksum(data, crc64Table)
			h := New()
			io.CopyBuffer(h, bytes.NewReader(data), make([]byte, 17))
			got := h.Sum64()
			if got != ref {
				t.Errorf("got 0x%x, want 0x%x", got, ref)
			}
		})
	}
}

func BenchmarkCrc64(b *testing.B) {
	b.Run("1MB", func(b *testing.B) {
		bench(b, New(), 1<<20)
	})
	b.Run("stdlib-1MB", func(b *testing.B) {
		bench(b, crc64.New(crc64Table), 1<<20)
	})
	b.Run("64KB", func(b *testing.B) {
		bench(b, New(), 64<<10)
	})
	b.Run("stdlib-64KB", func(b *testing.B) {
		bench(b, crc64.New(crc64Table), 64<<10)
	})
	b.Run("4KB", func(b *testing.B) {
		bench(b, New(), 4<<10)
	})
	b.Run("stdlib-4KB", func(b *testing.B) {
		bench(b, crc64.New(crc64Table), 4<<10)
	})
	b.Run("1KB", func(b *testing.B) {
		bench(b, New(), 1<<10)
	})
	b.Run("stdlib-1KB", func(b *testing.B) {
		bench(b, crc64.New(crc64Table), 1<<10)
	})
}

func bench(b *testing.B, h hash.Hash64, size int64) {
	b.SetBytes(size)
	data := make([]byte, size)
	for i := range data {
		data[i] = byte(i)
	}
	in := make([]byte, 0, h.Size())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h.Reset()
		h.Write(data)
		h.Sum(in)
	}
}
