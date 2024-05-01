package db

import matrix "pir/PIR/Matrix"

var q int64 = 2147483647

type Server struct {
	Data matrix.Matrix
	A1   matrix.Matrix
	A2   matrix.Matrix
}

type SingleHintReply struct {
	Hint matrix.Matrix
}

type DoubleHintReply struct {
	HintS matrix.Matrix
	HintC matrix.Matrix
}

func (S *Server) Hint(reply SingleHintReply) {
	reply.Hint.Multiply(S.Data, S.A1)
}

func (S *Server) DoubleHint(reply DoubleHintReply) {
	preDecomp := matrix.MakeMatrix(S.Data.Rows, S.A1.Columns, 0, q)
	//Need to add decomposition/composition method?
	preDecomp.Multiply(S.Data, S.A1)

	reply.HintS = matrix.Decompose(preDecomp)

	reply.HintC.Multiply(reply.HintS, S.A2)

}
