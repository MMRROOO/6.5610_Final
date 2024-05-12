package helper

import (
	matrix "pir/PIR/Matrix"
	"testing"
)

func TestMatrixToFileNames(t *testing.T) {

	/*
		Input:
		M = [
			255,
			254,
			253,
			...
			0
		]

		Expected:
		E = [255*256+254, 253*256+252, ... 1*256+0]
	*/
	M := matrix.MakeMatrix(256, 1, 0, 256)
	for i := 0; i < 256; i++ {
		M.Set(i, 0, int64(255-i))
	}
	R := MatrixToFileNames(M)
	E := make([]int, 0)
	for i := 128; i > 0; i-- {
		val := 256*(2*i-1) + 2*(i-1)
		E = append(E, val)
	}

	for i := 0; i < 128; i++ {
		if R[i] != E[i] {
			t.Errorf("Expected %v, got %v", E, R)
			return
		}
	}
}

func TestFileFromMatrices(t *testing.T) {
	// effectively reverse the effect of Matrix to FileNames
	M := matrix.MakeMatrix(256, 1, 0, 256)
	for i := 0; i < 256; i++ {
		M.Set(i, 0, int64(255-i))
	}
	R := FileNamestoMatrices(MatrixToFileNames(M))

	for i := 0; i < 256; i++ {
		if R.Get(i, 0) != M.Get(i, 0) {
			t.Errorf("Expected %v, got %v", M, R)
			return
		}
	}
}
