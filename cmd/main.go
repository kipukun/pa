package main

import (
	"fmt"
	"log"

	"github.com/kipukun/pa"
)

func main() {
	h := [][]float64{{1.47, 1.50, 1.52, 1.55, 1.57, 1.60, 1.63, 1.65, 1.68, 1.70, 1.73, 1.75, 1.78, 1.80, 1.83}}
	X := pa.NewMatrix(h, nil)
	y := pa.NewMatrix([][]float64{{52.21, 53.12, 54.48, 55.84, 57.20, 58.57, 59.93, 61.29, 63.11, 64.47, 66.28, 68.10, 69.92, 72.19, 74.46}}, nil)

	clf := new(pa.LinearRegression[float64])
	err := clf.Fit(X.T(), y.T())
	if err != nil {
		log.Fatalln(err)
	}

	//yh, err := clf.Predict(pa.NewMatrix([][]float64{{1.66}, {2}}, nil))
	//if err != nil {
	//	log.Fatalln(err)
	//}
	r2, err := clf.Score(X.T(), y.T())
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(r2)
}
