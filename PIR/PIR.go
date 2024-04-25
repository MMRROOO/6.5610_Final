package main

var DBCOLUMNS = 256 // sqrt of DB
var DBROWS = 256    // sqrt of DB

func Query(column int, secret Matrix) encryption {
	qu := MakeMatrix(DBCOLUMNS, 1, 0)

	qu.Set(column, 0, 1)

	e := ENC(secret, qu)

	return e

}

func Ans(DB Matrix, qu encryption) Matrix {
	quMatrix := MatrixFromEncryption(qu)
	a := MakeMatrix(DB.Rows, quMatrix.Columns, 0)
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
	secret := MakeMatrix(N, 1, 1)

	DB := MakeMatrix(DBROWS, DBCOLUMNS, 1)
	DB.LWERound()

	qu := Query(0, secret)
	Ans := Ans(DB, qu)
	out := Reconstruct(Ans, secret)

	out.Print()
	DB.PrintColumn(0)
}
