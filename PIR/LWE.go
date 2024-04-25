package main

import "fmt"

var M = 100
var N = 10
var q int64 = 65521

type encryption struct {
	A Matrix
	b Matrix
}

func ENC(secret Matrix, v Matrix) encryption {
	A := MakeMatrix(M, N, 1)

	b := MakeMatrix(1, N, 0)

	b.Mupltiply(A, secret)
	b.AddError(q / 4)
	v.ScalarMupltiply(q / 2)
	b.Add((v))

	retval := encryption{A: A, b: b}
	return retval

}

func DEC(secret Matrix, A Matrix, b Matrix) Matrix {
	As := MakeMatrix(1, N, 0)
	As.Mupltiply(A, secret)

	c := Copy(b)

	c.Subtract(As)

	c.LWERound()

	return c

}

func main() {
	secret := MakeMatrix(1, N, 0)

	v := MakeMatrix(1, N, 1)
	v.LWERound()

	enc := ENC(secret, v)

	dec := DEC(secret, enc.A, enc.b)

	fmt.Print(dec.Data[0])

}
