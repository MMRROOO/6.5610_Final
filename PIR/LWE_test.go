package main

import (
	lwe "pir/PIR/LWE"
	matrix "pir/PIR/Matrix"
	"testing"
)

func TestEncDec(t *testing.T) {

	secret := matrix.MakeMatrix(lwe.N, 1, 1, q)
	pt := matrix.MakeMatrix(lwe.M, 1, 1, q)
	pt.LWERound()
	enc := lwe.ENC(secret, pt)
	dec := lwe.DEC(secret, enc.A, enc.B)
	for i := 0; i < lwe.M; i++ {
		if pt.Get(i, 0) != dec.Get(i, 0) {
			t.Errorf("decryption failed at row %d", i)
			t.Log("Expected ", pt.Get(i, 0), " got ", dec.Get(i, 0))
			// leave early to avoid spamming the output
			return
		}
	}
}

func TestLWERound(t *testing.T) {
	pt := matrix.MakeMatrix(lwe.M, 1, 1, q)
	pt.LWERound()
	for i := 0; i < lwe.M; i++ {
		if pt.Get(i, 0) >= q {
			t.Errorf("LWERound failed at row %d", i)
			t.Log("Expected ", pt.Get(i, 0), " to be less than ", q)
			// leave early to avoid spamming the output
			return
		}
	}
}
