package main

var M = 100
var N = 10
var q int64 = 65521

type encryption struct {
	A Matrix
	b Matrix
}

func ENC(secret Matrix, v Matrix) encryption {
	M := v.Rows
	N := secret.Rows
	A := MakeMatrix(M, N, 1)

	As := MakeMatrix(M, 1, 0)

	v_copy := Copy(v)

	As.Mupltiply(A, secret)
	As.AddError(q / 8)
	v_copy.ScalarMupltiply(q / 2) //change this to fix aliasing issue
	As.Add((v_copy))

	retval := encryption{A: A, b: As}
	return retval

}

func DEC(secret Matrix, A Matrix, b Matrix) Matrix {
	As := MakeMatrix(A.Rows, 1, 0)
	As.Mupltiply(A, secret)

	c := Copy(b)

	c.Subtract(As)

	c.LWERound()

	return c

}

// func main() {
// 	secret := MakeMatrix(N, 1, 1)

// 	v := MakeMatrix(M, 1, 1)
// 	v.LWERound()
// 	fmt.Print("before ENC\n")
// 	enc := ENC(secret, v)
// 	fmt.Print("after ENC\n")

// 	dec := DEC(secret, enc.A, enc.b)
// 	fmt.Print("after DEC\n")

// 	fmt.Print(dec.Data)
// 	fmt.Print("\n")
// 	fmt.Print(v.Data)

// 	fmt.Print("\n")
// 	enc.A.Print()
// 	enc.b.Print()

// }
