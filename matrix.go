package pa

import (
	"fmt"
	"reflect"
	"strings"
)

type Axis int

const (
	Column Axis = iota
	Row
	All
)

type Matrix[T Number] struct {
	data       [][]T
	rows, cols int
	index      []int
	columns    []string
	err        error
}

// Err returns the first error encountered from an operation performed on m.
// Users should check Err() after performing operations, as no further operations are performed after the first error.
// Users should not trust the result of any operation chain if Err() is not nil.
func (m *Matrix[T]) Err() error {
	return m.err
}

func NewMatrix[T Number](d [][]T, columns []string) *Matrix[T] {
	m := new(Matrix[T])
	m.data = d
	if len(d) < 1 {
		return m
	}
	l := len(d[0])
	for _, row := range d {
		if l != len(row) {
			panic("all rows must be the same size")
		}
		l = len(row)
	}
	m.rows, m.cols = len(d), len(d[0])
	m.index = []int{}
	for i := 0; i < len(d); i++ {
		m.index = append(m.index, i)
	}
	m.columns = columns
	return m
}

func (m *Matrix[T]) Size() (int, int) {
	return m.rows, m.cols
}

func (m *Matrix[T]) sizet() (T, T) {
	return any(m.rows).(T), any(m.cols).(T)
}

func (m *Matrix[T]) String() string {
	var b strings.Builder
	for i, col := range m.columns {
		if i == len(m.columns)-1 {
			b.WriteString(fmt.Sprintf("%s\n", col))
		} else {
			b.WriteString(fmt.Sprintf("%s,", col))
		}
	}
	for i, row := range m.data {
		b.WriteString(fmt.Sprintf("%d: ", i))
		for j, col := range row {
			if j == len(row)-1 {
				b.WriteString(fmt.Sprintf("%v\n", col))
			} else {
				b.WriteString(fmt.Sprintf("%v,", col))
			}
		}
	}
	return b.String()
}

func (m *Matrix[T]) Square() bool {
	return m.cols == m.rows
}

func (m *Matrix[T]) T() *Matrix[T] {
	if m.err != nil {
		return m
	}
	if m.rows == 1 && m.cols == 1 {
		return m
	}

	data := make([][]T, m.cols)
	for i := range data {
		data[i] = make([]T, m.rows)
	}
	for i := 0; i < m.cols; i++ {
		for j := 0; j < m.rows; j++ {
			data[i][j] = m.data[j][i]
		}
	}

	t := NewMatrix(data, nil)

	return t
}

func (m *Matrix[T]) Add(b *Matrix[T]) *Matrix[T] {
	if m.err != nil {
		return m
	}

	mm, mn := m.Size()
	bm, bn := b.Size()

	if mm != bm || mn != bn {
		m.err = fmt.Errorf("m.Add: addition of two unequal sized matrices is undefined: (%d x %d) and (%d x %d)", mm, mn, bm, bn)
		return m
	}

	data := make([][]T, mm)
	for i := range data {
		data[i] = make([]T, mn)
	}
	for i, row := range m.data {
		for j := range row {
			data[i][j] = m.data[i][j] + b.data[i][j]
		}
	}

	return NewMatrix(data, m.columns)
}

func Empty[T Number](r, c int) *Matrix[T] {
	d := make([][]T, r)
	for i := range d {
		d[i] = make([]T, c)
	}
	return NewMatrix(d, nil)
}

func Filled[T Number](r, c int, s T) *Matrix[T] {
	e := Empty[T](r, c)
	e.Fill(s)
	return e
}

func (m *Matrix[T]) Fill(s T) {
	for i := range m.data {
		for j := range m.data[i] {
			m.data[i][j] = s
		}
	}
}

func (m *Matrix[T]) Sub(b *Matrix[T]) *Matrix[T] {
	var negone T
	negone++
	negone = -negone
	return m.Add(b.Product(negone))
}

func (m *Matrix[T]) Sum(ax Axis) *Matrix[T] {
	mn, mm := m.Size()
	var data [][]T
	var l int
	switch ax {
	case Row:
		data = m.data
		l = mn
	case Column:
		data = m.T().data
		l = mm
	}
	r := make([]T, l)
	for i, row := range data {
		r[i] = sum(row)
	}
	return NewMatrix(append(make([][]T, 0), r), nil)
}

func (m *Matrix[T]) Mean(ax Axis) *Matrix[T] {
	mn, mm := m.Size()
	var d []T
	switch ax {
	case Row:
		d = make([]T, mn)
		for i, row := range m.data {
			d[i] = mean(row)
		}
	case Column:
		d = make([]T, mm)
		for i, row := range m.T().data {
			d[i] = mean(row)
		}
	default:
		return m
	}
	return NewMatrix(append(make([][]T, 0), d), nil)
}

func (m *Matrix[T]) Product(scalar T) *Matrix[T] {
	if m.err != nil {
		return m
	}
	mn, mm := m.Size()
	r := make([][]T, mn)
	for i := range r {
		r[i] = make([]T, mm)
	}
	for i, row := range r {
		for j := range row {
			r[i][j] = m.data[i][j] * scalar
		}
	}
	return NewMatrix(r, nil)
}

func (m *Matrix[T]) Usage() int {
	return int(reflect.TypeOf(m).Size())
}

func (m *Matrix[T]) Mul(b *Matrix[T]) *Matrix[T] {
	if m.err != nil {
		return m
	}
	mi, mj := m.Size()
	bi, bj := b.Size()
	if mj != bi {
		m.err = fmt.Errorf("matrix multiplcation where m.rows = %d and b.columns = %d is undefined", mj, bi)
		return m
	}

	data := make([][]T, mi)
	for i := range data {
		data[i] = make([]T, bj)
	}

	for i := 0; i < mi; i++ {
		for j := 0; j < bj; j++ {
			var sum T
			for k := 0; k < mj; k++ {
				sum = sum + m.data[i][k]*b.data[k][j]
			}
			data[i][j] = sum
		}
	}
	return NewMatrix(data, nil)
}

// Delete drops the nth member of axis and returns the modified matrix.
func (m *Matrix[T]) Delete(n int, axis Axis) *Matrix[T] {
	if m.err != nil {
		return m
	}
	mm, mn := m.Size()
	data := make([][]T, mm)
	for i := range data {
		data[i] = make([]T, mn)
	}
	switch axis {
	case Row:
		if n > mm {
			m.err = fmt.Errorf("row number out of bounds: %d > %d", n, mm)
			return m
		}
		data = append(m.data[:n], m.data[n+1:]...)
	case Column:
		if n > mn {
			m.err = fmt.Errorf("column number out of bounds: %d > %d", n, mn)
			return m
		}
		for i, row := range m.data {
			data[i] = append(row[:n], row[n+1:]...)
		}
	}
	return NewMatrix(data, nil)
}

// Minor returns m with the ith row and jth column removed
func (m *Matrix[T]) Minor(i, j int) *Matrix[T] {
	if m.err != nil {
		return m
	}
	return m.Delete(i, Row).Delete(j, Column)
}

// Applies f element wise
func (m *Matrix[T]) Apply(f func(T) T) *Matrix[T] {
	mn, mm := m.Size()
	data := make([][]T, mn)
	for i := range data {
		data[i] = make([]T, mm)
	}
	for i, row := range data {
		for j := range row {
			data[i][j] = f(m.data[i][j])
		}
	}
	return NewMatrix(data, nil)
}

func (m *Matrix[T]) Reverse(axis Axis) *Matrix[T] {
	if m.err != nil {
		return m
	}
	mn, mm := m.Size()
	data := make([][]T, mn)
	for i := range data {
		data[i] = make([]T, mm)
	}
	switch axis {
	case Row:
		for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
			data[i], data[j] = m.data[j], m.data[i]
		}
	case Column:
		for i := range data {
			for j, k := 0, len(data[i])-1; j < k; j, k = j+1, k-1 {
				data[i][j], data[i][k] = m.data[i][k], m.data[i][j]
			}
		}
	}

	return NewMatrix(data, nil)
}

func NewIdentity[T Number](size int) *Matrix[T] {
	data := make([][]T, size)
	for i := range data {
		data[i] = make([]T, size)
	}
	for i := range data {
		data[i][i] = 1
	}

	return NewMatrix(data, nil)
}

// use an interface here, since there are many implementations of matrix factoring
type factored[T Number] interface {
	factor(m *Matrix[T])
	Det() T
	Solve(m *Matrix[T]) *Matrix[T]
	Inverse(m *Matrix[T]) *Matrix[T]
}

type lupDecomp[T Number] struct {
	p, u, l *Matrix[T]
	lowers  [][]T
	s       int
}

func (lupd *lupDecomp[T]) factor(m *Matrix[T]) {
	mi, _ := m.Size()
	data := make([][]T, len(m.data))
	copy(data, m.data)
	lupd.u = NewMatrix(data, nil)
	lupd.p = NewIdentity[T](mi)
	lupd.lowers = make([][]T, mi)
	for i := range lupd.lowers {
		lupd.lowers[i] = make([]T, mi)
	}

	for step := 0; step < mi-1; step++ {
		// perform swapping
		lupd.u.data[step+1], lupd.u.data[step] = lupd.u.data[step], lupd.u.data[step+1]
		lupd.p.data[step+1], lupd.p.data[step] = lupd.p.data[step], lupd.p.data[step+1]

		for i, row := range lupd.u.data[step+1:] {
			base := lupd.u.data[step]
			lower := row[step] / base[step]
			lupd.u.data[step+1+i] = Array[T](row).Sub(Array[T](base).Scale(lower))
			lupd.lowers[step][i] = lower
		}

		lupd.s++
	}

	l := NewMatrix(lupd.lowers, nil).T().Reverse(Row)
	for i := range l.data {
		l.data[i][i] = 1
	}

	lupd.l = l
}

func Factor[T Number](m *Matrix[T]) (factored[T], error) {
	lupd := new(lupDecomp[T])
	lupd.factor(m)
	return lupd, nil
}

func (lupd *lupDecomp[T]) Det() T {
	var lp, up T
	lp = 1
	up = 1
	for i := 0; i < len(lupd.u.data); i++ {
		up *= lupd.u.data[i][i]
		lp *= lupd.l.data[i][i]
	}

	if lupd.s%2 == 0 {
		return lp * up
	}

	return -lp * up
}

func (lupd *lupDecomp[T]) Solve(b *Matrix[T]) *Matrix[T] {

	rhs := lupd.p.Mul(b)

	xs := make([]T, len(lupd.l.data))
	for i, row := range lupd.l.data {
		before := make([]T, len(row))
		for j, v := range row {
			before[j] = xs[j] * v
		}
		xs[i] = (rhs.data[i][0] - sum(before)) / lupd.l.data[i][i]
	}
	Y := NewMatrix(append(make([][]T, 0), xs), nil).T()
	xs = make([]T, len(lupd.u.data))
	for i := len(lupd.u.data) - 1; i >= 0; i-- {
		front := make([]T, len(lupd.u.data[i]))
		for j := len(lupd.u.data[i]) - 1; j >= 0; j-- {
			front[j] = xs[j] * lupd.u.data[i][j]
		}
		xs[i] = (Y.data[i][0] - sum(front)) / lupd.u.data[i][i]
	}

	return NewMatrix(append(make([][]T, 0), xs), nil).T()
}

func (lupd *lupDecomp[T]) Inverse(m *Matrix[T]) *Matrix[T] {
	cols := [][]T{}
	b := NewIdentity[T](len(m.data))
	for _, row := range b.T().data {
		res := lupd.Solve(NewMatrix(append(make([][]T, 0), row), nil).T())
		cols = append(cols, res.T().data...)
	}
	return NewMatrix(cols, nil).T()
}
