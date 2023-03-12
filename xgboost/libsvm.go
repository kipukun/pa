package xgboost

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

func DMatrixFromLibSVM(r io.Reader) ([][]float64, error) {
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanWords)
	run := 0
	mat := [][]float64{}
	working := []float64{}
	for scanner.Scan() {
		cur := scanner.Text()
		if !strings.Contains(cur, ":") { // at the beginning of a line
			run = 0
			mat = append(mat, working)
			working = []float64{}
			f, err := strconv.ParseFloat(cur, 64)
			if err != nil {
				return nil, fmt.Errorf("error parsing float: %w", err)
			}
			working = append(working, f)
			run++
		} else {
			split := strings.Split(cur, ":")
			idx, err := strconv.Atoi(split[0])
			if err != nil {
				return nil, err
			}
			val, err := strconv.ParseFloat(split[1], 64)
			if err != nil {
				return nil, err
			}
			length := idx - run
			if length == 0 {
				working = append(working, val)
			} else if length < 0 {
				return nil, errors.New("negative index found relative to file")
			} else {
				working = append(working, append(make([]float64, length-1), val)...)
			}
			run = idx

		}
	}
	err := scanner.Err()
	if err != nil {
		return nil, fmt.Errorf("error scanning from reader: %w", err)
	}
	return mat, nil
}
