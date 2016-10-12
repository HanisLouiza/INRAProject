// matrix_test.go
package matrix

import (
	"testing"
)

func TestNewMatrix(t *testing.T) {
	Matrix := NewMatrix(5, 5)
	if len(Matrix.samples) == 5 {
		t.Log("Matrix samples has the good size!\n")
	} else {
		t.Errorf("Matrix samples has the wrong size : %v\n", len(Matrix.samples))
	}
	if len(Matrix.data) == 25 {
		t.Log("Matrix data has the good size!\n")
	} else {
		t.Errorf("Matrix data has the wrong size : %v\n", len(Matrix.data))
	}
}

func TestMatrixGetLine(t *testing.T) {
	MatrixOut := newMatrixOut(5, 5)
	if len(MatrixOut.GetDataLine(1)) == 5 {
		t.Log("MatrixOut samples has the good size!\n")
	} else {
		t.Errorf("MatrixOut has the wrong size : %v\n", len(MatrixOut.GetDataLine(1)))
	}

}
