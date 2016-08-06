package test

import (
	"reflect"
	"testing"
)

type SubS1 struct {
	A, B string
}

type SubS2 struct {
	C string
	S SubS1
}

type S struct {
	D string

	S1 SubS1
	S2 SubS2
}

var testReflectIndex = [][]int{
	[]int{0},
	[]int{1, 0},
	[]int{1, 1},
	[]int{2, 0},
	[]int{2, 1, 0},
	[]int{2, 1, 1},
}

var testValue = S{"this", SubS1{"is", "a"}, SubS2{"test", SubS1{"of", "generation"}}}

//go:noinline
func length(value int, s string) int {
	return value + len(s)
}

func generatedLengthSubS1(l int, v *SubS1) int {
	l = length(l, v.A)
	l = length(l, v.B)
	return l
}

func generatedLengthSubS2(l int, v *SubS2) int {
	l = length(l, v.C)
	l = generatedLengthSubS1(l, &v.S)
	return l
}

func generatedLengthS(l int, v *S) int {
	l = length(l, v.D)
	l = generatedLengthSubS1(l, &v.S1)
	l = generatedLengthSubS2(l, &v.S2)
	return l
}

func BenchmarkGenerated(b *testing.B) {
	var l int
	for i := 0; i < b.N; i++ {
		l = generatedLengthS(0, &testValue)
	}
	b.SetBytes(int64(l))
}

func reflectLengthSimpleImpl(l int, obj reflect.Value) int {
	v := reflect.Indirect(obj)

	switch v.Kind() {
	case reflect.String:
		l = length(l, v.String())

	case reflect.Struct:
		n := v.NumField()
		for i := 0; i < n; i++ {
			l = reflectLengthSimpleImpl(l, v.Field(i))
		}
	default:
		panic("error in test")
	}

	return l
}

func reflectLengthSimple(l int, obj interface{}) int {
	return reflectLengthSimpleImpl(l, reflect.ValueOf(obj))
}

func reflectLengthIndex(l int, obj interface{}) int {
	v := reflect.Indirect(reflect.ValueOf(obj))

	for _, key := range testReflectIndex {
		f := v.FieldByIndex(key)

		switch f.Kind() {
		case reflect.String:
			l = length(l, f.String())
		default:
			panic("error in test")
		}
	}

	return l
}

func BenchmarkReflect(b *testing.B) {
	var l int
	for i := 0; i < b.N; i++ {
		l = reflectLengthSimple(0, &testValue)
	}
	b.SetBytes(int64(l))
}

func BenchmarkReflectIndex(b *testing.B) {
	var l int
	for i := 0; i < b.N; i++ {
		l = reflectLengthIndex(0, &testValue)
	}
	b.SetBytes(int64(l))
}
