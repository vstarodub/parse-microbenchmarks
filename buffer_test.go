package test

import (
	"bytes"
	"sync"
	"testing"
)

const Size = 1048576
const ChunkSize = 512
const MaxChunkSize = 16384

func runBuffer() {
	var b bytes.Buffer
	for i := 0; i < Size; i++ {
		b.WriteByte('x')
	}
}

func runPreallocatedCopy(buf []byte) {
	for i := 0; i < len(buf); i++ {
		buf[i] = 'x'
	}
}

func runAllocatedCopy() []byte {
	out := make([]byte, Size)

	for i := 0; i < Size; i++ {
		out[i] = 'x'
	}
	return out
}

func runAllocatedAppend() []byte {
	out := make([]byte, 0, Size)

	for i := 0; i < Size; i++ {
		out = append(out, 'x')
	}
	return out
}

func runAppend() []byte {
	var out []byte

	for i := 0; i < Size; i++ {
		out = append(out, 'x')
	}
	return out
}

func runChunks() [][]byte {
	var out [][]byte

	for i := 0; i < Size; {
		chunk := make([]byte, 0, ChunkSize)
		left := ChunkSize
		if left > Size-i {
			left = Size - i
		}
		for j := 0; j < left; j++ {
			chunk = append(chunk, 'x')
		}

		i += left
		out = append(out, chunk)
	}
	return out
}

func runExponentialChunks() [][]byte {
	var out [][]byte

	size := ChunkSize
	for i := 0; i < Size; {
		chunk := make([]byte, 0, size)

		left := size
		if left > Size-i {
			left = Size - i
		}

		for j := 0; j < left; j++ {
			chunk = append(chunk, 'x')
		}

		i += left

		size *= 2
		if size > MaxChunkSize {
			size = MaxChunkSize
		}

		out = append(out, chunk)
	}
	return out
}

var pool = sync.Pool{
	New: func() interface{} {
		return make([]byte, 0, MaxChunkSize)
	},
}

func runExponentialPooledChunks(Size int) [][]byte {
	out := make([][]byte, 0, 10)

	size := ChunkSize
	for i := 0; i < Size; {
		var chunk []byte
		if size == MaxChunkSize {
			chunk = pool.Get().([]byte)
		} else {
			chunk = make([]byte, 0, size)
		}

		left := size
		if left > Size-i {
			left = Size - i
		}

		for j := 0; j < left; j++ {
			chunk = append(chunk, 'x')
		}

		i += left

		size *= 2
		if size > MaxChunkSize {
			size = MaxChunkSize
		}

		out = append(out, chunk)
	}
	return out
}

func poolChunks(chunks [][]byte) {
	for _, c := range chunks {
		if cap(c) == MaxChunkSize {
			pool.Put(c[:0])
		}
	}
}

func BenchmarkBytesBuffer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		runBuffer()
	}
	b.SetBytes(Size)
}

func BenchmarkPreallocatedCopy(b *testing.B) {
	buf := make([]byte, Size)
	for i := 0; i < b.N; i++ {
		runPreallocatedCopy(buf)
	}
	b.SetBytes(Size)
}

func BenchmarkAllocatedCopy(b *testing.B) {
	for i := 0; i < b.N; i++ {
		runAllocatedCopy()
	}
	b.SetBytes(Size)
}

func BenchmarkAllocatedAppend(b *testing.B) {
	for i := 0; i < b.N; i++ {
		runAllocatedAppend()
	}
	b.SetBytes(Size)
}

func BenchmarkAppend(b *testing.B) {
	for i := 0; i < b.N; i++ {
		runAppend()
	}
	b.SetBytes(Size)
}

func BenchmarkChunks(b *testing.B) {
	for i := 0; i < b.N; i++ {
		runChunks()
	}
	b.SetBytes(Size)
}

func BenchmarkExponentialChunks(b *testing.B) {
	for i := 0; i < b.N; i++ {
		runExponentialChunks()
	}
	b.SetBytes(Size)
}

func BenchmarkExponentialPooledChunks(b *testing.B) {
	for i := 0; i < b.N; i++ {
		chunks := runExponentialPooledChunks(Size)
		poolChunks(chunks)
	}
	b.SetBytes(Size)
}

func BenchmarkParalleExponentialPooledChunks(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			chunks := runExponentialPooledChunks(Size)
			poolChunks(chunks)
		}
	})
	b.SetBytes(Size)
}
func BenchmarkParalleExponentialPooledChunksSmall(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			chunks := runExponentialPooledChunks(SmallSize)
			poolChunks(chunks)
		}
	})
	b.SetBytes(SmallSize)
}

func concatChunks(chunks [][]byte) []byte {
	size := 0
	for _, b := range chunks {
		size += len(b)
	}

	out := make([]byte, 0, size)
	for _, b := range chunks {
		out = append(out, b...)
	}

	return out
}

func BenchmarkExponentialPooledChunksConcat(b *testing.B) {
	for i := 0; i < b.N; i++ {
		chunks := runExponentialPooledChunks(Size)
		out := concatChunks(chunks)
		poolChunks(chunks)

		b.SetBytes(int64(len(out)))
	}
}

const SmallSize = 127

func runSmallAppend() []byte {
	buf := make([]byte, 0, SmallSize)

	for i := 0; i < SmallSize; i++ {
		buf = append(buf, 'x')
	}

	return buf
}

type Serializer struct {
	initial [128]byte
}

func (s *Serializer) runWithBuf(buf []byte) {
	for i := 0; i < SmallSize; i++ {
		buf = append(buf, 'x')
	}
}

func (s *Serializer) Run() {
	buf := s.initial[:0]
	s.runWithBuf(buf)
}

func BenchmarkParallelSmallAppend(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			runSmallAppend()
		}
	})
	b.SetBytes(SmallSize)
}

func BenchmarkParallelSmallStackBuf(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			var s Serializer
			s.Run()
		}
	})
	b.SetBytes(SmallSize)

}
