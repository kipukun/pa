package pa

import (
	"reflect"
	"testing"
)

func TestMatrix_T(t *testing.T) {
	tests := []struct {
		name        string
		given, want *Matrix[float64]
	}{
		{
			name:  "square",
			given: NewMatrix([][]float64{{1, 2}, {3, 4}}, nil),
			want:  NewMatrix([][]float64{{1, 3}, {2, 4}}, nil),
		},
		{
			name:  "unit",
			given: NewMatrix([][]float64{{1}}, nil),
			want:  NewMatrix([][]float64{{1}}, nil),
		},
		{
			name:  "rectangle",
			given: NewMatrix([][]float64{{1, 2, 3, 4}, {5, 6, 7, 8}}, nil),
			want:  NewMatrix([][]float64{{1, 5}, {2, 6}, {3, 7}, {4, 8}}, nil),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.given.T(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("\nMatrix.T()\n%swant\n%s", got, tt.want)
			}
		})
	}
}

func TestMatrix_Add(t *testing.T) {
	tests := []struct {
		name       string
		m, n, want *Matrix[float64]
	}{
		{
			name: "square",
			m:    NewMatrix([][]float64{{1, 1}, {1, 1}}, []string{"one", "two"}),
			n:    NewMatrix([][]float64{{2, 2}, {2, 2}}, []string{"one", "two"}),
			want: NewMatrix([][]float64{{3, 3}, {3, 3}}, []string{"one", "two"}),
		},
		{
			name: "unit",
			m:    NewMatrix([][]float64{{1}}, []string{"one"}),
			n:    NewMatrix([][]float64{{1}}, []string{"one"}),
			want: NewMatrix([][]float64{{2}}, []string{"one"}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.Add(tt.n); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Matrix.Add() = %s, want %s", got, tt.want)
			}
		})
	}
}

func TestMatrix_Mul(t *testing.T) {
	tests := []struct {
		name       string
		m, n, want *Matrix[float64]
	}{
		{
			name: "simple",
			m:    NewMatrix([][]float64{{1, 2}, {3, 4}}, nil),
			n:    NewMatrix([][]float64{{5, 6}, {7, 8}}, nil),
			want: NewMatrix([][]float64{{19, 22}, {43, 50}}, nil),
		},
		{
			name: "inverse",
			m:    NewMatrix([][]float64{{-1, 3.0 / 2.0}, {1, -1}}, nil),
			n:    NewMatrix([][]float64{{2, 3}, {2, 2}}, nil),
			want: NewIdentity[float64](2),
		},
		{
			name: "lu decomp",
			m:    NewMatrix([][]float64{{0, 1, 0}, {0, 0, 1}, {1, 0, 0}}, nil),
			n:    NewMatrix([][]float64{{0, 5, 22.0 / 3.0}, {4, 2, 1}, {2, 7, 9}}, nil),
			want: NewMatrix([][]float64{{4, 2, 1}, {2, 7, 9}, {0, 5, 22.0 / 3.0}}, nil),
		},
		{
			name: "clf1",
			n:    NewMatrix([][]float64{{1, 1, 1, 1}, {1, 2, 3, 4}}, nil).T(),
			m:    NewMatrix([][]float64{{1, 1, 1, 1}, {1, 2, 3, 4}}, nil),
			want: NewMatrix([][]float64{{4, 10}, {10, 30}}, nil),
		},
		{
			name: "clf2",
			m:    NewMatrix([][]float64{{1, 1, 1, 1}, {1, 2, 3, 4}}, nil),
			n:    NewMatrix([][]float64{{1, 3, 3, 5}}, nil).T(),
			want: NewMatrix([][]float64{{12}, {36}}, nil),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.Mul(tt.n); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("\nMatrix.Mul(n):\n%s\nwant:\n%s", got, tt.want)
			}
		})
	}
}

func TestMatrix_Delete(t *testing.T) {
	type args struct {
		n    int
		axis Axis
	}
	tests := []struct {
		name   string
		source *Matrix[int]
		args   args
		want   *Matrix[int]
	}{
		{
			"one",
			NewMatrix([][]int{{1}}, nil),
			args{
				0,
				Row,
			},
			NewMatrix([][]int{}, nil),
		},
		{
			"two",
			NewMatrix([][]int{{1, 0}, {0, 1}}, nil),
			args{
				1,
				Column,
			},
			NewMatrix([][]int{{1}, {0}}, nil),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.source.Delete(tt.args.n, tt.args.axis); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("\nM.Delete():\n%v\nwant:\n%v", got, tt.want)
			}
		})
	}
}

func TestNewIdentity(t *testing.T) {
	type args struct {
		size int
	}
	tests := []struct {
		name string
		args args
		want *Matrix[int]
	}{
		{
			"one",
			args{
				1,
			},
			NewMatrix([][]int{{1}}, nil),
		},
		{
			"two",
			args{
				2,
			},
			NewMatrix([][]int{{1, 0}, {0, 1}}, nil),
		},
		{
			"three",
			args{
				3,
			},
			NewMatrix([][]int{{1, 0, 0}, {0, 1, 0}, {0, 0, 1}}, nil),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewIdentity[int](tt.args.size); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewIdentity() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMatrix_Sub(t *testing.T) {
	tests := []struct {
		name       string
		m, n, want *Matrix[float64]
	}{
		{
			name: "simple",
			m:    NewMatrix([][]float64{{5, 6}, {7, 8}}, []string{"one", "two"}),
			n:    NewMatrix([][]float64{{1, 2}, {3, 4}}, []string{"one", "two"}),
			want: NewMatrix([][]float64{{4, 4}, {4, 4}}, []string{"one", "two"}),
		},
		{
			name: "single",
			m:    NewMatrix([][]float64{{1}}, nil),
			n:    NewMatrix([][]float64{{0}}, nil),
			want: NewMatrix([][]float64{{1}}, nil),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.Sub(tt.n); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("\nMatrix.Mul(n):\n%s\nwant:\n%s", got, tt.want)
			}
		})
	}
}

func TestMatrix_Sum(t *testing.T) {
	tests := []struct {
		name    string
		axis    Axis
		m, want *Matrix[float64]
	}{
		{
			name: "simple row",
			m:    NewMatrix([][]float64{{0.5, 1.5}}, nil),
			axis: Row,
			want: NewMatrix([][]float64{{2.0}}, nil),
		},
		{
			name: "simple col",
			m:    NewMatrix([][]float64{{0.5, 1.5}}, nil),
			axis: Column,
			want: NewMatrix([][]float64{{0.5, 1.5}}, nil),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.Sum(tt.axis); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("\nMatrix.Sum(axis):\n%s\nwant:\n%s", got, tt.want)
			}
		})
	}
}

func TestMatrix_Product(t *testing.T) {
	tests := []struct {
		name    string
		n       float64
		m, want *Matrix[float64]
	}{
		{
			name: "x2",
			m:    NewMatrix([][]float64{{1, 2, 3}}, nil),
			n:    2,
			want: NewMatrix([][]float64{{2, 4, 6}}, nil),
		},
		{
			name: "x10",
			m:    NewMatrix([][]float64{{1, 2, 3}, {4, 5, 6}}, nil),
			n:    10,
			want: NewMatrix([][]float64{{10, 20, 30}, {40, 50, 60}}, nil),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.Product(tt.n); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("\nMatrix.Product(n):\n%s\nwant:\n%s", got, tt.want)
			}
		})
	}
}

func TestMatrix_Reverse(t *testing.T) {
	tests := []struct {
		name    string
		axis    Axis
		m, want *Matrix[float64]
	}{
		{
			name: "simple row",
			m:    NewMatrix([][]float64{{1, 2, 3}, {4, 5, 6}}, nil),
			axis: Row,
			want: NewMatrix([][]float64{{4, 5, 6}, {1, 2, 3}}, nil),
		},
		{
			name: "simple col",
			m:    NewMatrix([][]float64{{1, 2, 3}, {4, 5, 6}}, nil).T(),
			axis: Column,
			want: NewMatrix([][]float64{{4, 5, 6}, {1, 2, 3}}, nil).T(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.Reverse(tt.axis); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("\nMatrix.Reverse(ax):\n%s\nwant:\n%s", got, tt.want)
			}
		})
	}
}

func TestMatrix_Mean(t *testing.T) {
	tests := []struct {
		name    string
		axis    Axis
		m, want *Matrix[float64]
	}{
		{
			name: "simple row",
			m:    NewMatrix([][]float64{{1, 2, 3}, {4, 5, 6}}, nil),
			axis: Row,
			want: NewMatrix([][]float64{{2, 5}}, nil),
		},
		{
			name: "simple col",
			m:    NewMatrix([][]float64{{1, 2, 3}, {4, 5, 6}}, nil),
			axis: Column,
			want: NewMatrix([][]float64{{2.5, 3.5, 4.5}}, nil),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.Mean(tt.axis); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("\nMatrix.Mean(ax):\n%s\nwant:\n%s", got, tt.want)
			}
		})
	}
}
