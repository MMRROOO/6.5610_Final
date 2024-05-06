package matrix

var M = 100
var N = 10

// var q int64 = 2147483647
// var DATA_SIZE int64 = 256

// var logQ_logP = int((math.Log(float64(q)))/math.Log(float64(DATA_SIZE))) + 1

// type Encryption struct {
// 	A matrix.Matrix
// 	B matrix.Matrix //testing
// }

func ENC(secret Matrix, v Matrix) Encryption {
	M := v.Rows
	N := secret.Rows
	A := MakeMatrix(M, N, 1, q)

	As := MakeMatrix(M, 1, 0, q)

	v_copy := Copy(v)

	As.Multiply(A, secret)
	As.AddError(4)
	v_copy.ScalarMultiply(q / DATA_SIZE) //change this to fix aliasing issue
	As.Add((v_copy))

	retval := Encryption{A: A, B: As}
	return retval

}

func DEC(secret Matrix, A Matrix, b Matrix) Matrix {
	As := MakeMatrix(A.Rows, 1, 0, q)
	As.Multiply(A, secret)

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
