package db

import matrix "pir/PIR/Matrix"

var q int64 = 2147483647

type Server struct {
	Data matrix.Matrix
	A1   matrix.Matrix
	A2   matrix.Matrix
}

func (S *Server) Hint(args *SingleHintArgs, reply *SingleHintReply) {
	reply.Hint.Multiply(S.Data, S.A1)
}

func (S *Server) DoubleHint(args *DoubleHintArgs, reply *DoubleHintReply) {
	preDecomp := matrix.MakeMatrix(S.Data.Rows, S.A1.Columns, 0, q)
	//Need to add decomposition/composition method?
	preDecomp.Multiply(S.Data, S.A1)

	HintS := matrix.Decompose(preDecomp)

	reply.HintC.Multiply(HintS, S.A2)

}
