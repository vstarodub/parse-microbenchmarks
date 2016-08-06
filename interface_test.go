package test

import (
	"testing"
)

type A struct {
	i int
	j int64
}
type I interface {
	F() int
	FInline() int
}

const Size = 12

//go:noinline
func (a *A) F() int {
	return a.i
}

func (a *A) FInline() int {
	return a.i
}

//go:noinline
func Noop(N int) int {
	for i := 0; i < N; i++ {
	}
	return 0
}

//go:noinline
func CallDirect(a *A, N int) int {
	ret := 0
	for i := 0; i < N; i++ {
		ret += a.F()
	}
	return ret
}

//go:noinline
func CallInterface(a I, N int) int {
	ret := 0
	for i := 0; i < N; i++ {
		ret += a.F()
	}
	return ret
}

//go:noinline
func CallInterfaceCast(a *A, N int) int {
	ret := 0
	for i := 0; i < N; i++ {
		var v I = a
		ret += v.F()
	}
	return ret
}

//go:noinline
func CallDirectI(a *A, N int) int {
	ret := 0
	for i := 0; i < N; i++ {
		ret += a.FInline()
	}
	return ret
}

//go:noinline
func CallInterfaceI(a I, N int) int {
	ret := 0
	for i := 0; i < N; i++ {
		ret += a.FInline()
	}
	return ret
}

//go:noinline
func CallInterfaceCastI(a *A, N int) int {
	ret := 0
	for i := 0; i < N; i++ {
		var v I = a
		ret += v.FInline()
	}
	return ret
}

func BenchmarkNoop(b *testing.B) {
	Noop(b.N)
	b.SetBytes(Size)
}

func BenchmarkCall_Direct(b *testing.B) {
	CallDirect(&A{1, 2}, b.N)
	b.SetBytes(Size)
}

func BenchmarkCall_Interface(b *testing.B) {
	CallInterface(&A{1, 2}, b.N)
	b.SetBytes(Size)
}

func BenchmarkCall_InterfaceCast(b *testing.B) {
	CallInterfaceCast(&A{1, 2}, b.N)
	b.SetBytes(Size)
}

func BenchmarkCall_DirectInline(b *testing.B) {
	CallDirectI(&A{1, 2}, b.N)
	b.SetBytes(Size)
}

func BenchmarkCall_InterfaceInline(b *testing.B) {
	CallInterfaceI(&A{1, 2}, b.N)
	b.SetBytes(Size)
}

func BenchmarkCall_InterfaceCastInline(b *testing.B) {
	CallInterfaceCastI(&A{1, 2}, b.N)
	b.SetBytes(Size)
}
