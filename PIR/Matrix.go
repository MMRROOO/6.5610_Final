package pir

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

type Matrix struct {
	Data     []int64
	Rows     int
	Collumns int
	q        int64
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
func MakeMatrix(Rows int, Collumns int, MType int) Matrix {
	if MType == 0 {
		return Matrix{Data: make([]int64, Rows*Collumns), Rows: Rows, Collumns: Collumns, q: 65521}
	} else if MType == 1 {
		M := Matrix{Data: make([]int64, Rows*Collumns), Rows: Rows, Collumns: Collumns, q: 65521}
		for i := 0; i < M.Rows; i++ {
			for j := 0; j < M.Collumns; j++ {
				M.Set(i, j, nrand()%M.q)
			}
		}
		return M
	}
	fmt.Print("incorrect MType")
	return Matrix{Data: make([]int64, 0), Rows: 0, Collumns: 0, q: 0}
}

func Copy(A Matrix) Matrix {
	C := Matrix{Data: make([]int64, A.Rows*A.Collumns), Rows: A.Rows, Collumns: A.Collumns, q: A.q}
	copy(C.Data, A.Data)
	return C
}

func (C *Matrix) Mupltiply(B Matrix, A Matrix) {

	// C := Matrix{Data: make([]int, A.Rows*B.Collumns), Rows: A.Rows, Collumns: B.Collumns}

	if A.Collumns != B.Rows || A.Rows != C.Rows || B.Collumns != C.Collumns {
		fmt.Printf("wrong Size Matrix")
		return
	}

	for j := 0; j < A.Rows; j++ {
		for i := 0; i < A.Collumns; i++ {
			for c := 0; c < B.Collumns; c++ {
				C.AddToIndex(j, i, A.Get(j, i)*B.Get(i, c))
			}
		}
	}
}

func (A *Matrix) Add(B Matrix) {
	if A.Collumns != B.Collumns || A.Rows != B.Rows {
		fmt.Printf("wrong Size Matrix")
		return
	}
	for i := 0; i < A.Rows; i++ {
		for j := 0; j < A.Collumns; j++ {
			A.AddToIndex(i, j, B.Get(i, j))
		}
	}
}
func (A *Matrix) Subtract(B Matrix) {
	if A.Collumns != B.Collumns || A.Rows != B.Rows {
		fmt.Printf("wrong Size Matrix")
		return
	}
	for i := 0; i < A.Rows; i++ {
		for j := 0; j < A.Collumns; j++ {
			A.AddToIndex(i, j, -B.Get(i, j))
		}
	}
}

func (A *Matrix) AddError(max_error int64) {
	for i := 0; i < A.Rows; i++ {
		for j := 0; j < A.Collumns; j++ {
			A.AddToIndex(i, j, nrand()%int64(max_error))
		}
	}
}

func (A *Matrix) LWERound() {
	for i := 0; i < A.Rows; i++ {
		for j := 0; j < A.Collumns; j++ {
			val := A.Get(i, j)
			if val < A.q/4 || val > 3*A.q/4 {
				A.Set(i, j, 0)
			} else {
				A.Set(i, j, 1)
			}
		}
	}
}

func (A *Matrix) ScalarMupltiply(value int64) {
	for j := 0; j < A.Rows; j++ {
		for i := 0; i < A.Collumns; i++ {
			A.MupltiplyToIndex(j, i, value)
		}
	}
}

func (A *Matrix) Get(row int, collumn int) int64 {
	return A.Data[row*A.Collumns+collumn]
}

func (A *Matrix) Set(row int, collumn int, value int64) {
	A.Data[row*A.Collumns+collumn] = value
}

func (A *Matrix) MupltiplyToIndex(row int, collumn int, value int64) {
	A.Data[row*A.Collumns+collumn] = (A.Data[row*A.Collumns+collumn] * value) % A.q
}
func (A *Matrix) AddToIndex(row int, collumn int, value int64) {
	A.Data[row*A.Collumns+collumn] = (A.Data[row*A.Collumns+collumn] + value) % A.q
}
