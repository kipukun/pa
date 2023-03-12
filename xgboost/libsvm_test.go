package xgboost

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

const (
	svm1 = "0 1:5 5:27.1 20:52.5\n1 2:5 5:98 53:5"
)

func TestDMatrixFromLibSVM(t *testing.T) {
	matrix, err := DMatrixFromLibSVM(strings.NewReader(svm1))
	if err != nil {
		t.Errorf("error creating matrix: %s", err.Error())
	}
	b, _ := json.Marshal(matrix)
	fmt.Println(string(b))
}
