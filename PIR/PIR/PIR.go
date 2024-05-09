package PIR

import (
	"math"
	matrix "pir/PIR/Matrix"
)

var DBCOLUMNS = 256 // sqrt of DB
var DBROWS = 256    // sqrt of DB
var SMALLN = 16     // sqrt of DB
var q int64 = 2147483647
var logQ_logP = int((math.Log(float64(q)))/math.Log(float64(matrix.DATA_SIZE))) + 1

var MAX_ERROR int64 = 10

func Query(column int, secret matrix.Matrix) matrix.Encryption {
	qu := matrix.MakeMatrix(DBCOLUMNS, 1, 0, q)

	qu.Set(column, 0, 1)

	e := matrix.ENC(secret, qu)

	return e

}

func Ans(DB matrix.Matrix, qu matrix.Encryption) matrix.Matrix {
	quMatrix := matrix.MatrixFromEncryption(qu)
	a := matrix.MakeMatrix(DB.Rows, quMatrix.Columns, 0, q)
	a.Multiply(DB, quMatrix)

	return a
}

func Reconstruct(ans matrix.Matrix, secret matrix.Matrix) matrix.Matrix {

	enc := matrix.EncryptionFromMatrix(ans)

	v := matrix.DEC(secret, enc.A, enc.B)

	return v
}

func DoubleSetup(DB matrix.Matrix, A1 matrix.Matrix, A2 matrix.Matrix) (HintS matrix.Matrix, HintC matrix.Matrix) {
	A1xDB := matrix.MakeMatrix(A1.Columns, DB.Rows, 0, q)

	A1xDB.Multiply(A1.Transpose(), DB.Transpose())

	HintS = matrix.Decompose(A1xDB)
	HintC = matrix.MakeMatrix(HintS.Rows, A2.Columns, 0, q)
	HintC.Multiply(HintS, A2)

	return HintS, HintC
}

type Dqu struct {
	c1 matrix.Matrix
	c2 matrix.Matrix
}

func DoubleQuery(row int, column int, secretC matrix.Matrix, secretR matrix.Matrix, A1 matrix.Matrix, A2 matrix.Matrix) Dqu {

	c1 := matrix.MakeMatrix(A1.Rows, 1, 0, q)
	c1.Multiply(A1, secretC)
	c1.AddError(MAX_ERROR)
	columnVector := matrix.MakeMatrix(A1.Rows, 1, 0, q)
	columnVector.Set(column, 0, 1)
	columnVector.ScalarMultiply(q / matrix.DATA_SIZE)
	c1.Add(columnVector)

	c2 := matrix.MakeMatrix(A2.Rows, 1, 0, q)
	c2.Multiply(A2, secretR)
	c2.AddError(MAX_ERROR)
	rowVector := matrix.MakeMatrix(A2.Rows, 1, 0, q)
	rowVector.Set(row, 0, 1)
	rowVector.ScalarMultiply(q / matrix.DATA_SIZE)
	c2.Add(rowVector)

	return Dqu{c1: c1, c2: c2}
}

type Dans struct {
	H    matrix.Matrix
	AnsH matrix.Matrix
	Ans2 matrix.Matrix
}

func DoubleAns(DB matrix.Matrix, HintS matrix.Matrix, qu Dqu, A2 matrix.Matrix) Dans {
	c1 := qu.c1
	c2 := qu.c2

	ans1 := matrix.MakeMatrix(c1.Columns, DB.Rows, 0, q)
	ans1.Multiply(c1.Transpose(), DB.Transpose())
	ans1 = matrix.Decompose(ans1)

	H := matrix.MakeMatrix(ans1.Rows, A2.Columns, 0, q)

	H.Multiply(ans1, A2)

	ansH := matrix.MakeMatrix(HintS.Rows, 1, 0, q)
	ansH.Multiply(HintS, c2)

	ans2 := matrix.MakeMatrix(logQ_logP, 1, 0, q)

	ans2.Multiply(ans1, c2)

	return Dans{H: H, AnsH: ansH, Ans2: ans2}
}

func DoubleReconstruct(secret1 matrix.Matrix, secret2 matrix.Matrix, HintC matrix.Matrix, ans Dans) int64 {

	HintC_H := matrix.JoinVertical(HintC, ans.H)
	tmp := matrix.MakeMatrix(HintC_H.Rows, secret1.Columns, 0, q)
	tmp.Multiply(HintC_H, secret2)
	AnsH_Ans2 := matrix.JoinVertical(ans.AnsH, ans.Ans2)

	AnsH_Ans2.Subtract(tmp)
	AnsH_Ans2.LWERound()
	H1_a1 := matrix.Compose(AnsH_Ans2)
	H1, a1 := matrix.SplitVertical(H1_a1)
	retval := matrix.MakeMatrix(1, 1, 0, q)
	retval.Multiply(secret1.Transpose(), H1)
	a1.Subtract(retval)
	a1.LWERound()
	return a1.Get(0, 0)

}

// func main() {
// 	//square root pir
// 	// secret := MakeMatrix(N, 1, 1, q)

// 	// DB := MakeMatrix(DBROWS, DBCOLUMNS, 1, q)
// 	// DB.LWERound()

// 	// qu := Query(0, secret)
// 	// Ans := Ans(DB, qu)
// 	// out := Reconstruct(Ans, secret)

// 	// out.Print()
// 	// DB.PrintColumn(0)

// 	//DoublePIR

// 	DB := matrix.MakeMatrix(DBROWS, DBCOLUMNS, 1, q)
// 	DB.LWERound()

// 	A1 := matrix.MakeMatrix(DBROWS, SMALLN, 1, q)
// 	A2 := matrix.MakeMatrix(DBCOLUMNS, SMALLN, 1, q)

// 	Secret1 := matrix.MakeMatrix(SMALLN, 1, 1, q)
// 	Secret2 := matrix.MakeMatrix(SMALLN, 1, 1, q)

// 	HintS, HintC := DoubleSetup(DB, A1, A2)

// 	fmt.Print("After Setup\n")
// 	row, col := 100, 18

// 	qu := DoubleQuery(row, col, Secret1, Secret2, A1, A2)

// 	fmt.Print("After Query\n")

// 	Ans := DoubleAns(DB, HintS, qu, A2)

// 	fmt.Print("After Ans\n")

// 	val := DoubleReconstruct(Secret1, Secret2, HintC, Ans)

// 	fmt.Print("After Reconstruct\n")

// 	fmt.Print(val, DB.Get(row, col))
// 	// DB.Print()
// }
