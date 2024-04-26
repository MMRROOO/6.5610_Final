package main

import (
	"testing"
)

func TestEncDec(t *testing.T){

	secret := MakeMatrix(N, 1, 1, q)
	pt := MakeMatrix(M, 1, 1, q)
	pt.LWERound()
	enc := ENC(secret, pt)
	dec := DEC(secret, enc.A, enc.b)
	for i := 0; i < M; i++ {
		if pt.Get(i, 0) != dec.Get(i, 0) {
			t.Errorf("decryption failed at row %d", i)
			t.Log("Expected ", pt.Get(i, 0), " got ", dec.Get(i, 0))
			// leave early to avoid spamming the output
			return
		}
	}
}

func TestLWERound(t *testing.T){
	pt := MakeMatrix(M, 1, 1, q)
	pt.LWERound()
	for i := 0; i < M; i++ {
		if pt.Get(i, 0) >= q {
			t.Errorf("LWERound failed at row %d", i)
			t.Log("Expected ", pt.Get(i, 0), " to be less than ", q)
			// leave early to avoid spamming the output
			return
		}
	}
}
