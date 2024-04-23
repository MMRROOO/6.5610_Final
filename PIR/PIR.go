// package pir

// import (
// 	"crypto/rand"
// 	"math/big"

// 	"gonum.org/v1/gonum/mat"
// )

// var DB_SIZE = 100 // sqrt of DB values
// func nrand() int64 {
// 	max := big.NewInt(int64(1) << 62)
// 	bigx, _ := rand.Int(rand.Reader, max)
// 	x := bigx.Int64()
// 	return x
// }

// func Query(index int) {
// 	key := nrand()
// 	queryVector := mat.NewDense(1, DB_SIZE, nil)

// }
