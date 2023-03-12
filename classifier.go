package pa

import "fmt"

type Classifier[T Number] interface {
	Fit(X, y *Matrix[T]) (err error)
	Predict(X *Matrix[T]) (y_hat *Matrix[T], err error)
	Score(X, y_true *Matrix[T]) (float64, error)
}

type LinearRegression[T Number] struct {
	bhat *Matrix[T]
}

func (lr *LinearRegression[T]) Fit(X, y *Matrix[T]) (err error) {
	data := make([][]T, len(X.data))
	for i, row := range X.data {
		data[i] = append([]T{1}, row...)
	}
	X = NewMatrix(data, nil)
	inner := X.T().Mul(X)
	if inner.Err() != nil {
		return err
	}
	factored, err := Factor(inner)
	if err != nil {
		return err
	}
	lr.bhat = factored.Inverse(inner).Mul(X.T().Mul(y))
	if lr.bhat.Err() != nil {
		return err
	}
	return nil
}

func (lr *LinearRegression[T]) Predict(X *Matrix[T]) (y_hat *Matrix[T], err error) {
	data := make([][]T, len(X.data))
	for i, row := range X.data {
		data[i] = append([]T{1}, row...)
	}
	X = NewMatrix(data, nil)
	yh := X.Mul(lr.bhat)
	if yh.Err() != nil {
		return nil, yh.Err()
	}
	return yh, nil
}

func (lr *LinearRegression[T]) Coefficients() *Matrix[T] {
	return lr.bhat
}

func (lr *LinearRegression[T]) Score(X, y *Matrix[T]) (float64, error) {
	yh, err := lr.Predict(X)
	if err != nil {
		return -1, err
	}
	ybar := y.Mean(Column)
	yc, yr := y.Size()
	ybar = Filled(yc, yr, ybar.data[0][0])
	fmt.Println(ybar)
	rss := yh.Sub(ybar).Apply(func(t T) T { return t * t }).Sum(Column)
	tss := y.Sub(ybar).Apply(func(t T) T { return t * t }).Sum(Column)
	return float64(1 - (rss.data[0][0] / tss.data[0][0])), nil
}
