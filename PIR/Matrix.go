package pir

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

type Matrix struct {
	Data     []int
	Rows     int
	Collumns int
	q        int
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
		return Matrix{Data: make(int[], Rows*Collumns), Rows: Rows, Collumns: Collumns, q: 65521}
	}else if MType == 1{
		M :=  Matrix{Data: make(int[], Rows*Collumns), Rows: Rows, Collumns: Collumns, q: 65521}
		for i:= 0; i< M.Rows; i++{
			for j:=0; jk< M.Collumns; j++{
				M.Set(i, j, nrand() % M.q)
			}
		}
		return M
	}
}

func (A *Matrix) AddError(max_error int) {
	for i:= 0; i< A.Rows; i++{
		for j:=0; jk< A.Collumns; j++{
			M.Set(i, j, nrand() % M.q)
		}
	}
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

func (A *Matrix) ScalarMupltiply(value int) {
	for j := 0; j < A.Rows; j++ {
		for i := 0; i < A.Collumns; i++ {
			A.MupltiplyToIndex(j, i, value)
		}
	}
}

func (A *Matrix) Get(row int, collumn int) int {
	return A.Data[row*A.Collumns+collumn]
}

func (A *Matrix) Set(row int, collumn int, value int) {
	A.Data[row*A.Collumns+collumn] = value
}

func (A *Matrix) MupltiplyToIndex(row int, collumn int, value int) {
	A.Data[row*A.Collumns+collumn] = (A.Data[row*A.Collumns+collumn] * value) % M.q
}
func (A *Matrix) AddToIndex(row int, collumn int, value int) {
	A.Data[row*A.Collumns+collumn] = (A.Data[row*A.Collumns+collumn] + value) % M.q
}
