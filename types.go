package pa

import (
	"golang.org/x/exp/constraints"
)

type Number interface {
	constraints.Signed | constraints.Float | constraints.Unsigned
}

type Data[T Number] interface {
	Array[T] | Matrix[T] | Number
}

type Array[T Number] []T

func (a Array[T]) Scale(s T) Array[T] {
	r := make(Array[T], len(a))
	for i, v := range a {
		r[i] = v * s
	}

	return r
}

func (a Array[T]) Sub(b Array[T]) Array[T] {
	r := make(Array[T], len(a))
	for i := range a {
		r[i] = a[i] - b[i]
	}

	return r
}
