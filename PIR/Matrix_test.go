package main

import (
	lwe "pir/PIR/LWE"
	matrix "pir/PIR/Matrix"
	"testing"
)

var Rows = 30
var Cols = 40

// zero matrix
var Z = matrix.MakeMatrix(Rows, Cols, 0, q)

// random matrix
var A = matrix.MakeMatrix(Rows, Cols, 1, q)
var B = matrix.MakeMatrix(Rows, Cols, 1, q)

func TestMakeMatrix(t *testing.T) {
	// t.Parallel()
	if Z.Rows != Rows || Z.Columns != Cols {
		t.Errorf("MakeMatrix made a matrix of size %d x %d; got %d x %d", Rows, Cols, Z.Rows, Z.Columns)
	}

	// check that all elements are zero
	for r := 0; r < Rows; r++ {
		for c := 0; c < Cols; c++ {
			if Z.Get(r, c) != 0 {
				t.Errorf("MakeMatrix made a matrix with non-zero element at row %d, column %d", r, c)
			}
		}
	}

	// check that A not equal to B; MakeMatrix properly makes random matrices
	isEqual := true
	for r := 0; r < Rows; r++ {
		for c := 0; c < Cols; c++ {
			if A.Get(r, c) != B.Get(r, c) {
				isEqual = false
				break
			}
		}
	}
	if isEqual {
		t.Errorf("MakeMatrix with randomization made two identical matrices")
	}

}

func TestCopy(t *testing.T) {
	// t.Parallel()
	Acopy := matrix.Copy(A)
	if (A.Rows != Acopy.Rows) || (A.Columns != Acopy.Columns) {
		t.Errorf("Copy did not copy matrix size")
	}

	for r := 0; r < Rows; r++ {
		for c := 0; c < Cols; c++ {
			if A.Get(r, c) != Acopy.Get(r, c) {
				t.Errorf("Copy did not copy element at row %d, column %d", r, c)
			}
		}
	}
}

func TestEncryptionFromMatrix(t *testing.T) {
	// t.Parallel()
	Ans := matrix.MakeMatrix(Rows, Cols, 1, q)
	encryption := matrix.EncryptionFromMatrix(Ans)
	eA := encryption.A
	eb := encryption.B

	if eA.Rows != Rows || eA.Columns != Cols-1 {
		t.Errorf("EncryptionFromMatrix made A matrix of size %d x %d; got %d x %d", Rows, Cols-1, eA.Rows, eA.Columns)
	}

	if eb.Rows != Rows || eb.Columns != 1 {
		t.Errorf("EncryptionFromMatrix made b matrix of size %d x %d; got %d x %d", Rows, 1, eb.Rows, eb.Columns)
	}

	for r := 0; r < Rows; r++ {
		for c := 0; c < Cols-1; c++ {
			if eA.Get(r, c) != Ans.Get(r, c) {
				t.Errorf("EncryptionFromMatrix did not copy element at row %d, column %d from Ans to A", r, c)
			}
		}
	}

	for r := 0; r < Rows; r++ {
		if eb.Get(r, 0) != Ans.Get(r, Cols-1) {
			t.Errorf("EncryptionFromMatrix did not copy element at row %d, column %d from Ans to b", r, Cols-1)
		}
	}
}

func TestMatrixFromEncryption(t *testing.T) {
	// t.Parallel()
	e := lwe.Encryption{A: A, B: matrix.MakeMatrix(Rows, 1, 1, q)}
	C := matrix.MatrixFromEncryption(e)
	if C.Rows != Rows || C.Columns != Cols+1 {
		t.Errorf("MatrixFromEncryption made matrix of size %d x %d; got %d x %d", Rows, Cols+1, C.Rows, C.Columns)
	}

	for r := 0; r < Rows; r++ {
		for c := 0; c < Cols; c++ {
			if c == Cols {
				if C.Get(r, c) != e.B.Get(r, 0) {
					t.Errorf("MatrixFromEncryption did not copy element at row %d, column %d from b to C", r, c)
				}
			} else {
				if C.Get(r, c) != e.A.Get(r, c) {
					t.Errorf("MatrixFromEncryption did not copy element at row %d, column %d from A to C", r, c)
				}
			}
		}
	}
}
