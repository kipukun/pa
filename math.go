package pa

import "fmt"

func Dot[T Number, A Array[T]](a, b A) (T, error) {
	var sum T
	if len(a) != len(b) {
		return sum, fmt.Errorf("dot product of arrays with different lengths is undefined: %d vs. %d", len(a), len(b))
	}
	for i := range a {
		sum += a[i] * b[i]
	}

	return sum, nil
}
