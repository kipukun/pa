package xgboost

//#cgo LDFLAGS: -L/usr/local/lib -lxgboost
//#cgo CFLAGS: -I/usr/local/include
//#include <assert.h>
//#include <stdio.h>
//#include <stdlib.h>
//#include <xgboost/c_api.h>
import "C"
import (
	"errors"
	"fmt"
	"unsafe"
)

var (
	errClosed = errors.New("closed")

	labelField = C.CString("label")
)

type cArray interface {
	*C.float | *C.DMatrixHandle
}

func arrtoC[T any, V cArray](arr []T) V {
	return (V)(unsafe.Pointer(&arr[0]))
}

type dMatrix struct {
	c             *C.DMatrixHandle
	rows, columns int
	missing       float64
	closed        bool
}

func (dmat *dMatrix) SetLabels(ls []float64) error {
	res := C.XGDMatrixSetFloatInfo(*dmat.c, C.CString("label"), arrtoC[float64, *C.float](ls), C.ulong(len(ls)))
	if res == -1 {
		lasterr := C.XGBGetLastError()
		err := fmt.Errorf("error setting labels for matrix: %s", C.GoString(lasterr))
		return err
	}
	return nil
}

// Close closes the handle to the underlying DMatrix. It is an error to operate on a closed DMatrix.
func (dmat *dMatrix) Close() error {
	if dmat.closed {
		return errors.New("called close on already closed DMatrix")
	}
	res := C.XGDMatrixFree(*dmat.c)
	if res == -1 {
		lasterr := C.XGBGetLastError()
		err := fmt.Errorf("error freeing matrix: %s", C.GoString(lasterr))
		return err
	}
	return nil
}

func CreateDMatrix(m []float64, rows, columns int, missing float64) (*dMatrix, error) {
	var mh C.DMatrixHandle
	res := C.XGDMatrixCreateFromMat(arrtoC[float64, *C.float](m), C.ulong(rows), C.ulong(columns), C.float(-1), &mh)
	if res == -1 {
		lasterr := C.XGBGetLastError()
		err := fmt.Errorf("error creating matrix: %s", C.GoString(lasterr))
		return nil, err
	}
	mat := &dMatrix{
		c:       &mh,
		rows:    rows,
		columns: columns,
		missing: missing,
	}
	return mat, nil
}

type Booster struct {
	c      *C.BoosterHandle
	closed bool
}

func (b *Booster) Close() error {
	if b.closed {
		return errors.New("called close on already closed Booster")
	}
	res := C.XGBoosterFree(*b.c)
	if res == -1 {
		lasterr := C.XGBGetLastError()
		err := fmt.Errorf("error freeing matrix: %s", C.GoString(lasterr))
		return err
	}
	return nil
}

func (b *Booster) Fit(m *dMatrix, iters int) error {
	for i := 0; i < iters; i++ {
		status := C.XGBoosterUpdateOneIter(*b.c, C.int(i), *m.c)
		if status == -1 {
			lasterr := C.GoString(C.XGBGetLastError())
			return fmt.Errorf("error iterating booster: %s", lasterr)
		}

		//status = C.XGBoosterEvalOneIter(booster, C.int(i), X_matrix, arrtoC([]), eval_dmats_size, &eval_result)

	}
	return nil
}

// CreateBooster initializes a Booster with training data m and any additional data ms.
func CreateBooster(m *dMatrix, ms ...*dMatrix) (*Booster, error) {
	cms := []C.DMatrixHandle{*m.c}
	for _, m := range ms {
		cms = append(cms, *m.c)
	}
	var booster C.BoosterHandle
	res := C.XGBoosterCreate(arrtoC[C.DMatrixHandle, *C.DMatrixHandle](cms), C.ulong(len(cms)), &booster)
	if res == -1 {
		lasterr := C.XGBGetLastError()
		return nil, fmt.Errorf("error creating booster: %s", C.GoString(lasterr))
	}
	b := &Booster{
		c: &booster,
	}
	return b, nil

}
