package main

var DBCOLUMNS = 256 // sqrt of DB
var DBROWS = 256    // sqrt of DB

var MAX_ERROR int64 = 5

func Query(column int, secret Matrix) encryption {
	qu := MakeMatrix(DBCOLUMNS, 1, 0, q)

	qu.Set(column, 0, 1)

	e := ENC(secret, qu)

	return e

}

func Ans(DB Matrix, qu encryption) Matrix {
	quMatrix := MatrixFromEncryption(qu)
	a := MakeMatrix(DB.Rows, quMatrix.Columns, 0, q)
	a.Mupltiply(DB, quMatrix)

	return a
}

func Reconstruct(ans Matrix, secret Matrix) Matrix {

	enc := EncryptionFromMatrix(ans)
	ans.Print()
	enc.b.Print()
	v := DEC(secret, enc.A, enc.b)

	return v
}

func main() {
	secret := MakeMatrix(N, 1, 1, q)

	DB := MakeMatrix(DBROWS, DBCOLUMNS, 1, q)
	DB.LWERound()

	qu := Query(0, secret)
	Ans := Ans(DB, qu)
	out := Reconstruct(Ans, secret)

	out.Print()
	DB.PrintColumn(0)
}

type Dqu struct {
	c1 Matrix
	c2 Matrix
}

func DoubleQuery(row int, column int, secretC Matrix, secretR Matrix, A1 Matrix, A2 Matrix) Dqu {

	c1 := MakeMatrix(A1.Rows, 1, 0, q)
	c1.Multiply(A1, secretC)
	c1.AddError(MAX_ERROR)
	columnVector := MakeMatrix(A1.Rows, 1, 0, q)
	columnVector.Set(column, 0, 1)
	columnVector.ScalarMultiply(q / DATA_SIZE)
	c1.Add(columnVector)

	c2 := MakeMatrix(A2.Rows, 1, 0, q)
	c2.Multiply(A2, secretR)
	c2.AddError(MAX_ERROR)
	rowVector := MakeMatrix(A2.Rows, 1, 0, q)
	rowVector.Set(row, 0, 1)
	rowVector.ScalarMultiply(q / DATA_SIZE)
	c2.Add(rowVector)

	return Dqu{c1: c1, c2: c2}
}

type Dans struct {
	H    Matrix
	AnsH Matrix
	Ans2 Matrix
}

func DoubleAns(DB Matrix, HintS Matrix, qu Dqu, A2 Matrix) Dans {
	c1 := qu.c1
	c2 := qu.c2

	ans1 := MakeMatrix(c1.Columns, DB.Rows, 0, q)
	ans1.Multiply(c1.Transpose(), DB.Transpose())
	ans1 = Decompose(ans1)

	H := MakeMatrix(ans1.Rows, A2.Columns, 0, q)

	H.Multiply(ans1, A2)

	ansH := MakeMatrix(HintS.Rows, 1, 9, q)
	ansH.Multiply(HintS, c2)

	ans2 := MakeMatrix(1, 1, 0, q)

	ans2.Multiply(ans1, c2)

	return Dans{H: H, AnsH: ansH, Ans2: ans2}
}

func DoubleReconstruct(secret1 Matrix, secret2 Matrix, HintC Matrix, ans Dans) int64 {

	HintC_H := JoinVertical(HintC, ans.H)
	AnsH_Ans2 := JoinVertical(ans.AnsH, ans.Ans2)

	AnsH_Ans2.Subtract(HintC_H)
	AnsH_Ans2.LWERound()
	H1_a1 := Compose(AnsH_Ans2)
	H1, a1 := SplitVertical(H1_a1)
	retval := MakeMatrix(1, 1, 0, q)
	retval.Multiply(secret1.Transpose(), H1)
	a1.Subtract(retval)

	return a1.Get(0, 0)

}
