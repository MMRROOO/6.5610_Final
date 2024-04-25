package main

var DBCOLLUMNS = 100 // sqrt of DB
var DBROWS = 100     // sqrt of DB

func Query(row int, collumn int, secret Matrix) encryption {
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

func Reonstruct(ans Matrix, secret Matrix) Matrix {

	enc := EncryptionFromMatrix(ans)

	v := DEC(secret, enc.A, enc.b)

	return v
}
