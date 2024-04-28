package main

import (
	"crypto/rand"
	"fmt"
	"math"
	"math/big"
)

type Matrix struct {
	Data    []int64
	Rows    int
	Columns int
	q       int64
}

func nrand() int64 {
	max := big.NewInt(int64(1) << 62)
	bigx, _ := rand.Int(rand.Reader, max)
	x := bigx.Int64()
	return x
}

// Mtypes:
// 0 = Zero Matrix
// 1 = Random Matrix
func MakeMatrix(Rows int, Columns int, MType int, q int64) Matrix {

	if MType == 0 {
		return Matrix{Data: make([]int64, Rows*Columns), Rows: Rows, Columns: Columns, q: q}
	} else if MType == 1 {
		M := Matrix{Data: make([]int64, Rows*Columns), Rows: Rows, Columns: Columns, q: q}
		for i := 0; i < M.Rows; i++ {
			for j := 0; j < M.Columns; j++ {
				M.Set(i, j, nrand()%M.q)
			}
		}
		return M
	}
	fmt.Print("incorrect MType")
	return Matrix{Data: make([]int64, 0), Rows: 0, Columns: 0, q: 0}
}

func Copy(A Matrix) Matrix {
	C := Matrix{Data: make([]int64, A.Rows*A.Columns), Rows: A.Rows, Columns: A.Columns, q: A.q}
	copy(C.Data, A.Data)
	return C
}

func EncryptionFromMatrix(Ans Matrix) encryption {
	A := MakeMatrix(Ans.Rows, Ans.Columns-1, 0, Ans.q)

	for i := 0; i < Ans.Rows; i++ {
		for j := 0; j < Ans.Columns-1; j++ {
			A.Set(i, j, Ans.Get(i, j))
		}
	}
	b := MakeMatrix(Ans.Rows, 1, 0, Ans.q)
	for i := 0; i < Ans.Rows; i++ {
		b.Set(i, 0, Ans.Get(i, Ans.Columns-1))
	}

	return encryption{A: A, b: b}

}

func MatrixFromEncryption(E encryption) Matrix {

	C := Matrix{Data: make([]int64, E.A.Rows*(E.A.Columns+1)), Rows: E.A.Rows, Columns: E.A.Columns + 1, q: E.A.q}
	for i := 0; i < E.A.Rows; i++ {
		for j := 0; j < (E.A.Columns + 1); j++ {
			if j == E.A.Columns {
				C.Data[i*(E.A.Columns+1)+j] = E.b.Get(i, 0)
			} else {
				C.Data[i*(E.A.Columns+1)+j] = E.A.Get(i, j)
			}
		}
	}

	return C
}

func (C *Matrix) Multiply(A Matrix, B Matrix) {

	// C := Matrix{Data: make([]int, A.Rows*B.Columns), Rows: A.Rows, Columns: B.Columns}

	if A.Columns != B.Rows || A.Rows != C.Rows || B.Columns != C.Columns {
		fmt.Printf("wrong Size Matrix, A rows, columns: %d, %d, B rows columns: %d, %d, C rows, columns: %d, %d\n", A.Rows, A.Columns, B.Rows, B.Columns, C.Rows, C.Columns)
		return
	}

	for j := 0; j < A.Rows; j++ {
		for i := 0; i < B.Columns; i++ {
			for c := 0; c < A.Columns; c++ {
				C.AddToIndex(j, i, A.Get(j, c)*B.Get(c, i))
			}
		}
	}
}

func (A *Matrix) Add(B Matrix) {
	if A.Columns != B.Columns || A.Rows != B.Rows {
		fmt.Printf("wrong Size Matrix\n")
		return
	}
	for i := 0; i < A.Rows; i++ {
		for j := 0; j < A.Columns; j++ {
			A.AddToIndex(i, j, B.Get(i, j))
		}
	}
}
func (A *Matrix) Subtract(B Matrix) {
	if A.Columns != B.Columns || A.Rows != B.Rows {
		fmt.Printf("wrong Size Matrix\n")
		return
	}
	for i := 0; i < A.Rows; i++ {
		for j := 0; j < A.Columns; j++ {
			A.SubtractFromIndex(i, j, B.Get(i, j))
		}
	}
}

func (A *Matrix) AddError(max_error int64) {
	for i := 0; i < A.Rows; i++ {
		for j := 0; j < A.Columns; j++ {
			A.AddToIndex(i, j, nrand()%int64(max_error)-(int64(max_error)/2))
		}
	}
}

func (A *Matrix) ScalarMultiply(value int64) {
	for j := 0; j < A.Rows; j++ {
		for i := 0; i < A.Columns; i++ {
			A.MultiplyToIndex(j, i, value)
		}
	}
}

func (A *Matrix) Get(row int, column int) int64 {
	return A.Data[row*A.Columns+column] % A.q
}

func (A *Matrix) Set(row int, column int, value int64) {
	A.Data[row*A.Columns+column] = value % A.q
}

func (A *Matrix) MultiplyToIndex(row int, column int, value int64) {
	A.Data[row*A.Columns+column] = (A.Data[row*A.Columns+column] * value) % A.q
}
func (A *Matrix) AddToIndex(row int, column int, value int64) {
	A.Data[row*A.Columns+column] = (A.Data[row*A.Columns+column] + value) % A.q
}

func (A *Matrix) SubtractFromIndex(row int, column int, value int64) {
	A.Data[row*A.Columns+column] = (A.Data[row*A.Columns+column] + (A.q - value)) % A.q
}

func (A *Matrix) Print() {
	for i := 0; i < A.Columns; i++ {
		fmt.Print("______")
	}
	fmt.Print("__")

	for i := 0; i < A.Rows; i++ {
		fmt.Print("\n|")
		for j := 0; j < A.Columns; j++ {
			val := A.Get(i, j)
			if val < 1000 {
				fmt.Print(" ")
			}
			if val < 10 {
				fmt.Print(" ")
			}
			fmt.Printf("%d", A.Get(i, j))
			if val < 10000 {
				fmt.Print(" ")
			}
			if val < 100 {
				fmt.Print(" ")
			}
			fmt.Print("|")
		}
	}
	fmt.Print("\n")

}

func (A *Matrix) GetColumn(column int) []int {
	C := make([]int, A.Rows)
	for i := 0; i < A.Rows; i++ {
		C[i] = int(A.Get(i, column))
	}
	return C
}
func (A *Matrix) PrintColumn(column int) {
	C := MakeMatrix(A.Rows, 1, 0, A.q)

	for i := 0; i < A.Rows; i++ {
		C.Set(i, 0, A.Get(i, column))
	}

	C.Print()
}

// each column of length k stores one value from original matrix
func Decompose(A Matrix) Matrix {
	DecompSize := logQ_logP
	B := MakeMatrix(A.Rows*DecompSize, A.Columns, 0, A.q)

	for r := 0; r < A.Rows; r++ {
		for c := 0; c < A.Columns; c++ {
			for k := 0; k < DecompSize; k++ {
				DecompVal := (A.Get(r, c) % (int64(math.Pow(float64(DATA_SIZE), float64(k+1))))) / (int64(math.Pow(float64(DATA_SIZE), float64(k))))
				B.Set(r*DecompSize+k, c, DecompVal)

			}
		}
	}
	return B
}

func Compose(A Matrix) Matrix {
	DecompSize := logQ_logP
	B := MakeMatrix(A.Rows/DecompSize, A.Columns, 0, A.q)

	for r := 0; r < B.Rows; r++ {
		for c := 0; c < B.Columns; c++ {
			for k := 0; k < DecompSize; k++ {
				CompVal := (A.Get(r*DecompSize+k, c) * (int64(math.Pow(float64(DATA_SIZE), float64(k)))))
				B.AddToIndex(r, c, CompVal)
			}
		}
	}
	return B
}

func (A *Matrix) Transpose() Matrix {
	T := MakeMatrix(A.Columns, A.Rows, 0, A.q)
	for r := 0; r < A.Rows; r++ {
		for c := 0; c < A.Columns; c++ {
			T.Set(c, r, A.Get(r, c))
		}
	}
	return T
}

func JoinVertical(A Matrix, B Matrix) Matrix {
	if A.Columns != B.Columns {
		fmt.Print("Different number of columns")
		return MakeMatrix(0, 0, 0, A.q)
	}
	J := MakeMatrix(A.Rows+B.Rows, A.Columns, 0, q)
	for r := 0; r < A.Rows+B.Rows; r++ {
		for c := 0; c < A.Columns; c++ {
			if r < A.Rows {
				J.Set(r, c, A.Get(r, c))
			} else {
				J.Set(r, c, B.Get(r-A.Rows, c))
			}
		}
	}
	return J
}

func SplitVertical(A Matrix) (Matrix, Matrix) {
	T := MakeMatrix(A.Rows-1, A.Columns, 0, A.q)
	B := MakeMatrix(1, A.Columns, 0, A.q)
	T.Data = A.Data[0 : A.Columns*(A.Rows-1)]
	B.Data = A.Data[A.Columns*(A.Rows-1) : A.Columns*(A.Rows)]

	return T, B
}
