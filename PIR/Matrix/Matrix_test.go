package matrix

import (
	"math"
	"testing"
)

var Rows = 30
var Cols = 40

func isEqual(A Matrix, B Matrix) (bool) {
	/*
	Tells whether matrices A, B are equal
	*/
	
	if A.Rows != B.Rows || A.Columns != B.Columns {
		return false
	}

	for r := 0; r < A.Rows; r++ {
		for c := 0; c < A.Columns; c++ {
			if A.Get(r, c) != B.Get(r, c) {
				return false
			}
		}
	}
	return true
}
func TestMakeMatrix(t *testing.T) {
	Z := MakeMatrix(Rows, Cols, 0, q)
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
	A := MakeMatrix(Rows, Cols, 1, q)
	B := MakeMatrix(Rows, Cols, 1, q)

	// check that A not equal to B; MakeMatrix properly makes random matrices
	if isEqual(A, B) {
		t.Errorf("MakeMatrix made two equal matrices, expected randomization to make them different")
	}

}

func TestCopy(t *testing.T) {
	// t.Parallel()
	A := MakeMatrix(Rows, Cols, 1, q)
	Acopy := Copy(A)
	if (A.Rows != Acopy.Rows) || (A.Columns != Acopy.Columns) {
		t.Errorf("Copy did not copy matrix size")
	}

	if !isEqual(A, Acopy) {
		t.Errorf("Copy did not copy matrix elements")
	}
}

func TestEncryptionFromMatrix(t *testing.T) {
	// t.Parallel()
	Ans := MakeMatrix(Rows, Cols, 1, q)
	encryption := EncryptionFromMatrix(Ans)
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
	A := MakeMatrix(Rows, Cols, 1, q)
	// t.Parallel()
	e := Encryption{A: A, B: MakeMatrix(Rows, 1, 1, q)}
	C := MatrixFromEncryption(e)
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

func TestDecompose(t *testing.T){
	/*
	k = ceiling(log q / log p)
	Maps an `a x b` mod q matrix into a `ka x b`mod q matrix as its base-p decomposition
	*/
	A := MakeMatrix(Rows, Cols, 1, q)

	p := DATA_SIZE
	k := logQ_logP

	decomp := Decompose(A)
	if decomp.Rows != k*Rows || decomp.Columns != Cols {
		t.Errorf("Decompose made matrix of size %d x %d; got %d x %d", k*Rows, Cols, decomp.Rows, decomp.Columns)
	}

	for r := 0; r < A.Rows; r++ {
		for c := 0; c < A.Columns; c++ {
			for i := 0; i < k; i++ {
				orig := A.Get(r, c)
				num := orig % int64(math.Pow(float64(p), float64(i + 1)))
				denom := (int64(math.Pow(float64(p), float64(i))))

				DecompVal := num / denom

				if decomp.Get(r*k+i, c) != DecompVal {
					t.Errorf("Decompose did not copy element at row %d, column %d from A to decomp", r, c)
					t.Log("Expected ", DecompVal, " got ", decomp.Get(r*k+i, c))
					return
				}


			}
		}
	}

}

func TestCompose(t *testing.T){
	A := MakeMatrix(Rows, Cols, 1, q)

	// interprets each 𝜅 × 1 submatrix of its input as a base-𝑝 decomposition of a Z𝑞 element and outputs the matrix of these elements.
	
	decomp := Decompose(A)
	composed := Compose(decomp)

	if !isEqual(A, composed) {
		t.Errorf("Compose did not return original matrix")
	}
}