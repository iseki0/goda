package goda

import (
	"cmp"
	"encoding"
	"math"
	"strconv"
)

func parseInt64(input []byte) (int64, error) {
	return strconv.ParseInt(string(input), 10, 64)
}

func parseInt(input []byte) (i int, e error) {
	i, e = strconv.Atoi(string(input))
	return
}

func unmarshalJsonImpl[T encoding.TextUnmarshaler](ref T, data []byte) error {
	if len(data) < 2 || data[0] != '"' || data[len(data)-1] != '"' {
		return parseFailedError(data)
	}
	return ref.UnmarshalText(data[1 : len(data)-1])
}

func marshalJsonImpl[T encoding.TextAppender](ref T) ([]byte, error) {
	d, e := ref.AppendText([]byte{'"'})
	if e != nil {
		return nil, e
	}
	d = append(d, '"')
	return d, nil
}

func marshalTextImpl[T encoding.TextAppender](ref T) ([]byte, error) {
	return ref.AppendText(nil)
}

func stringImpl[T encoding.TextAppender](ref T) string {
	b, e := ref.AppendText(nil)
	if e != nil {
		panic(e)
	}
	return bytes2string(b)
}

func bytes2string(b []byte) string {
	return string(b)
}

func floorDiv(x, y int64) (q int64) {
	q = x / y
	if (x^y) < 0 && (q*y != x) {
		q = q - 1
	}
	return
}

func floorMod(x, y int64) (r int64) {
	r = x % y
	if (x^y) < 0 && r != 0 {
		r = r + y
	}
	return
}

type comparable0[T any] interface {
	Compare(T) int
	IsZero() bool
}

func comparing1[E any, T comparable0[T]](f func(E) T) func(E, E) int {
	return func(a, b E) int {
		return f(a).Compare(f(b))
	}
}

func comparing[E any, T cmp.Ordered](f func(E) T) func(E, E) int {
	return func(a, b E) int {
		return cmp.Compare(f(a), f(b))
	}
}

func doCompare[E any](a, b E, f ...func(E, E) int) int {
	for _, it := range f {
		var i = it(a, b)
		if i != 0 {
			return i
		}
	}
	return 0
}

func compareZero[T interface{ IsZero() bool }](a, b T) int {
	if a.IsZero() {
		if b.IsZero() {
			return 0
		}
		return -1
	}
	if b.IsZero() {
		return 1
	}
	return 0
}

func mustValue[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

func addExactly(x, y int64) (r int64, overflow bool) {
	r = x + y
	overflow = ((x ^ r) & (y ^ r)) < 0
	return
}

func mulExact(x, y int64) (r int64, overflow bool) {
	r = x * y
	overflow = ((y != 0) && (r/y != x)) || (x == math.MinInt64 && y == -1)
	return
}
