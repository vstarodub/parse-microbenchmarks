package test

import (
	"reflect"
	"testing"
)

type V struct {
	AAAA string
	BBBB string
	CCCC string
	DDDD string
}

type testStreamItem struct {
	Name, Value string
}

var testStream = []testStreamItem{
	{"CCCC", "this"},
	{"BBBB", "is"},
	{"DDDD", "a"},
	{"AAAA", "test"},
}

func generatedParse(v *V) int {
	var ret int
	for _, it := range testStream {
		ret += len(it.Name) + len(it.Value)
		switch it.Name {
		case "AAAA":
			v.AAAA = it.Value
		case "BBBB":
			v.BBBB = it.Value
		case "CCCC":
			v.CCCC = it.Value
		case "DDDD":
			v.DDDD = it.Value
		}
	}
	return ret
}

var states [256 * 6]int

func init() {
	// default: 0, error state

	// only valid transitions from 1, 2, 3, 4 are to self
	states['A'+256*1] = 1
	states['B'+256*2] = 2
	states['C'+256*3] = 3
	states['D'+256*4] = 4

	// initial state
	states['A'+256*5] = 1
	states['B'+256*5] = 2
	states['C'+256*5] = 3
	states['D'+256*5] = 4
}

func generatedParseDFA(v *V) int {
	var ret int
	for _, it := range testStream {
		ret += len(it.Name) + len(it.Value)

		n := 5
		field := it.Name

		for i := 0; i < len(field); i++ {
			c := int(field[i])
			n = states[c+n*256]
		}

		switch n {
		case 1:
			v.AAAA = it.Value
		case 2:
			v.BBBB = it.Value
		case 3:
			v.CCCC = it.Value
		case 4:
			v.DDDD = it.Value
		}

	}
	return ret
}

func BenchmarkGenerated(b *testing.B) {
	var v V
	var l int

	for i := 0; i < b.N; i++ {
		l = generatedParse(&v)
	}
	b.SetBytes(int64(l))
}

func BenchmarkGeneratedDFA(b *testing.B) {
	var v V
	var l int

	for i := 0; i < b.N; i++ {
		l = generatedParseDFA(&v)
	}
	b.SetBytes(int64(l))
}

func reflectSimple(obj interface{}) int {
	v := reflect.Indirect(reflect.ValueOf(obj))

	var ret int
	for _, it := range testStream {
		ret += len(it.Name) + len(it.Value)

		f := v.FieldByName(it.Name)
		f.SetString(it.Value)
	}
	return ret
}

var index = map[string]int{
	"AAAA": 0,
	"BBBB": 1,
	"CCCC": 2,
	"DDDD": 3,
}

func reflectIndex(obj interface{}) int {
	v := reflect.Indirect(reflect.ValueOf(obj))

	var ret int
	for _, it := range testStream {
		ret += len(it.Name) + len(it.Value)
		f := v.Field(index[it.Name])
		f.SetString(it.Value)
	}
	return ret
}

func BenchmarkReflect(b *testing.B) {
	var v V
	var l int
	for i := 0; i < b.N; i++ {
		l = reflectSimple(&v)
	}
	b.SetBytes(int64(l))
}

func BenchmarkReflectIndex(b *testing.B) {
	var v V
	var l int
	for i := 0; i < b.N; i++ {
		l = reflectIndex(&v)
	}
	b.SetBytes(int64(l))
}
