package helper

import (
	matrix "pir/PIR/Matrix"
	// pir "pir/PIR/PIR"
)

// type Peer struct {
// 	peers  []p2p.Endpoint
// 	mu     sync.Mutex
// 	me     p2p.Endpoint
// 	Data   matrix.Matrix
// 	secret matrix.Matrix
// 	Hashes []byte
// 	Host   p2p.Endpoint
// }

/*
Args:
V: matrix.Matrix of size 256 x 1. Each value in `V` is a byte (i.e., 0 <= V.Get(i, 0) < 256 for all i

Output:
Returns indices of the files that the matrix represents. Expected size is 128-- return every other byte value.
*/
func MatrixToFileNames(M matrix.Matrix) []int {
	// check that all values v in M satisfy 0 <= v < 256
	for r := 0; r < 256; r++ {
		if M.Get(r, 0) < 0 || M.Get(r, 0) >= 256 {
			panic("Matrix values must be between 0 and 255")
		}
	}

	res := make([]int, 0)
	cur := 0
	lookVal := 0 // the val currently looking at when iterating through the matrix
	for r := 0; r < 256; r++ {
		lookVal = int(M.Get(r, 0))
		if r%2 == 0 {
			cur = lookVal * 256
		} else {
			cur += lookVal
			res = append(res, cur)
			cur = 0
		}
	}
	return res
}

/*
Args: []int `files` - list of file indices. Reverse the effect of MatrixToFileNames

Output: matrix.Matrix of size 256 x 1. Each value in the matrix represents a byte
*/
func FileNamestoMatrices(files []int) matrix.Matrix {
	res := matrix.MakeMatrix(256, 1, 0, 256)
	for i := 0; i < 128; i++ {
		res.Set(i*2, 0, int64(files[i]/256))
		res.Set(i*2+1, 0, int64(files[i]%256))
	}
	return res
}

func MatrixtoFileChunk(M []matrix.Matrix) []byte {
	res := make([]byte, 256*4)
	for i := 0; i < 4; i++ {
		for j := 0; j < 256; j++ {

			res[i*256+j] = byte(M[i].Get(j, 0))
		}
	}
	return res
}

func FileChunkToMatrix(M []byte) matrix.Matrix {
	res := matrix.MakeMatrix(256, 4, 0, 2147483647)
	for i := 0; i < 4; i++ {
		for j := 0; j < 256; j++ {
			res.Set(j, i, int64(M[i*256+j]))
		}
	}
	return res
}
