package main

var DBCOLLUMNS = 10 // sqrt of DB
var DBROWS = 10     // sqrt of DB

func Query(collumn int, secret Matrix) encryption {
	qu := MakeMatrix(DBCOLLUMNS, 1, 0)

	qu.Set(collumn, 1, 1)

	e := ENC(secret, qu)

	return e

}

func Ans(DB Matrix, qu encryption) Matrix {
	quMatrix := MatrixFromEncryption(qu)
	a := MakeMatrix(DB.Rows, quMatrix.Collumns, 0)
	a.Mupltiply(DB, quMatrix)

	return a
}

func Reconstruct(ans Matrix, secret Matrix) Matrix {

	enc := EncryptionFromMatrix(ans)

	v := DEC(secret, enc.A, enc.b)

	return v
}

func main() {
	secret := MakeMatrix(N, 1, 1)

	DB := MakeMatrix(DBROWS, DBCOLLUMNS, 1)
	DB.LWERound()

	qu := Query(0, secret)
	Ans := Ans(DB, qu)
	out := Reconstruct(Ans, secret)

	out.Print()
	DB.Print()
}
