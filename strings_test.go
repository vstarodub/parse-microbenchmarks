package test

import (
	"testing"
)

var testString = `an \"esca\"ped\"string of a great length"`

func runOnePass() []byte {
	buf := make([]byte, 0, 8)

	for i := 0; i < len(testString); i++ {
		c := testString[i]
		switch c {
		case '"':
			return buf
		case '\\':
			i++
			buf = append(buf, testString[i])
		default:
			buf = append(buf, c)
		}
	}
	return buf
}

func runTwoPass() []byte {
	l := 0

	for i := 0; i < len(testString); i++ {
		switch testString[i] {
		case '"':
			i = len(testString)

		case '\\':
			l++
			i++
		default:
			l++
		}
	}

	buf := make([]byte, 0, l)

	for i := 0; i < len(testString); i++ {
		c := testString[i]
		switch c {
		case '"':
			return buf
		case '\\':
			i++
			buf = append(buf, testString[i])
		default:
			buf = append(buf, c)
		}
	}
	return buf
}

func BenchmarkOnePass(b *testing.B) {
	for i := 0; i < b.N; i++ {
		runOnePass()
	}
	b.SetBytes(int64(len(testString)))
}

func BenchmarkTwoPass(b *testing.B) {
	for i := 0; i < b.N; i++ {
		runTwoPass()
	}
	b.SetBytes(int64(len(testString)))
}

//go:noinline
func runCopy(in, out []byte) {
	in = append(out, in...)
}

func BenchmarkCopy(b *testing.B) {
	buf1 := make([]byte, 11)
	buf2 := make([]byte, 0, 11)

	for i := 0; i < b.N; i++ {
		runCopy(buf1, buf2)
	}

	b.SetBytes(11)
}

//go:noinline
func runAppend(in, out []byte) {
	for i := 0; i < len(in); i++ {
		out = append(out, in[i])
	}
}

func BenchmarkAppend(b *testing.B) {
	buf1 := make([]byte, 11)
	buf2 := make([]byte, 0, 11)

	for i := 0; i < b.N; i++ {
		runAppend(buf1, buf2)
	}

	b.SetBytes(11)
}
