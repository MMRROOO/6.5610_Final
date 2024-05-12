package matrix

import (
	"crypto/rand"
	"fmt"
	"math"
	"math/big"
	insecure_rand "math/rand" // Not cryptographically secure change back at some point
	"time"
)

var DATA_SIZE int64 = 256
var q int64 = 2147483647

type Matrix struct {
	Data    []int64
	Rows    int
	Columns int
	Q       int64
}

type Encryption struct {
	A Matrix
	B Matrix //testing
}

var logQ_logP = int((math.Log(float64(q)))/math.Log(float64(DATA_SIZE))) + 1

func nrand() int64 { //secure implementation
	max := big.NewInt(int64(1) << 62)
	bigx, _ := rand.Int(rand.Reader, max)
	x := bigx.Int64()
	return x
}

func insecure_nrand(seed int64) int64 { //unsecure for testing only
	insecure_rand.New(insecure_rand.NewSource(seed)) // Seed the default Source in the math/rand package

	max_val := big.NewInt(1)           // Create a big.Int
	max_val = max_val.Lsh(max_val, 62) // Left shift to get 2^62

	// Generate a random *big.Int in range [0, max_val)
	return insecure_rand.New(insecure_rand.NewSource(time.Now().UnixNano())).Int63n(max_val.Int64())
}

// Mtypes:
// 0 = Zero Matrix
// 1 = Random Matrix
func MakeMatrix(Rows int, Columns int, MType int, q int64) Matrix {

	seed := 1079

	if MType == 0 {
		return Matrix{Data: make([]int64, Rows*Columns), Rows: Rows, Columns: Columns, Q: q}
	} else if MType == 1 {
		M := Matrix{Data: make([]int64, Rows*Columns), Rows: Rows, Columns: Columns, Q: q}
		for i := 0; i < M.Rows; i++ {
			for j := 0; j < M.Columns; j++ {
				M.Set(i, j, insecure_nrand(int64(seed)%M.Q))
			}
		}
		return M
	}
	fmt.Print("incorrect MType")
	return Matrix{Data: make([]int64, 0), Rows: 0, Columns: 0, Q: 0}
}

func Copy(A Matrix) Matrix {
	C := Matrix{Data: make([]int64, A.Rows*A.Columns), Rows: A.Rows, Columns: A.Columns, Q: A.Q}
	copy(C.Data, A.Data)
	return C
}

func EncryptionFromMatrix(Ans Matrix) Encryption {
	A := MakeMatrix(Ans.Rows, Ans.Columns-1, 0, Ans.Q)

	for i := 0; i < Ans.Rows; i++ {
		for j := 0; j < Ans.Columns-1; j++ {
			A.Set(i, j, Ans.Get(i, j))
		}
	}
	b := MakeMatrix(Ans.Rows, 1, 0, Ans.Q)
	for i := 0; i < Ans.Rows; i++ {
		b.Set(i, 0, Ans.Get(i, Ans.Columns-1))
	}

	return Encryption{A: A, B: b}

}

func MatrixFromEncryption(E Encryption) Matrix {

	C := Matrix{Data: make([]int64, E.A.Rows*(E.A.Columns+1)), Rows: E.A.Rows, Columns: E.A.Columns + 1, Q: E.A.Q}
	for i := 0; i < E.A.Rows; i++ {
		for j := 0; j < (E.A.Columns + 1); j++ {
			if j == E.A.Columns {
				C.Data[i*(E.A.Columns+1)+j] = E.B.Get(i, 0)
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
	return A.Data[row*A.Columns+column] % A.Q
}

func (A *Matrix) Set(row int, column int, value int64) {
	A.Data[row*A.Columns+column] = value % A.Q
}

func (A *Matrix) MultiplyToIndex(row int, column int, value int64) {
	A.Data[row*A.Columns+column] = (A.Data[row*A.Columns+column] * value) % A.Q
}
func (A *Matrix) AddToIndex(row int, column int, value int64) {
	A.Data[row*A.Columns+column] = (A.Data[row*A.Columns+column] + value) % A.Q
}

func (A *Matrix) SubtractFromIndex(row int, column int, value int64) {
	A.Data[row*A.Columns+column] = (A.Data[row*A.Columns+column] + (A.Q - value)) % A.Q
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
	C := MakeMatrix(A.Rows, 1, 0, A.Q)

	for i := 0; i < A.Rows; i++ {
		C.Set(i, 0, A.Get(i, column))
	}

	C.Print()
}

// each column of length k stores one value from original matrix
func Decompose(A Matrix) Matrix {
	DecompSize := logQ_logP
	B := MakeMatrix(A.Rows*DecompSize, A.Columns, 0, A.Q)

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
	B := MakeMatrix(A.Rows/DecompSize, A.Columns, 0, A.Q)

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
	T := MakeMatrix(A.Columns, A.Rows, 0, A.Q)
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
		return MakeMatrix(0, 0, 0, A.Q)
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

func (A *Matrix) CopyColumn(B Matrix, C int) {
	for i := 0; i < A.Rows; i++ {
		A.Set(i, C, B.Get(i, 0))
	}
}

func SplitVertical(A Matrix) (Matrix, Matrix) {
	T := MakeMatrix(A.Rows-1, A.Columns, 0, A.Q)
	B := MakeMatrix(1, A.Columns, 0, A.Q)
	T.Data = A.Data[0 : A.Columns*(A.Rows-1)]
	B.Data = A.Data[A.Columns*(A.Rows-1) : A.Columns*(A.Rows)]

	return T, B
}

func (A *Matrix) LWERound() {
	for i := 0; i < A.Rows; i++ {
		for j := 0; j < A.Columns; j++ {

			val := A.Get(i, j)
			roundedVal := (((val + (q/(DATA_SIZE))/2) / (q / DATA_SIZE)) % A.Q) % (DATA_SIZE)

			A.Set(i, j, roundedVal)
		}
	}
}

func IsEqual(A Matrix, B Matrix) bool {
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
