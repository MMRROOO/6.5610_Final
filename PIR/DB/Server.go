package db

import (
	matrix "pir/PIR/Matrix"
	pir "pir/PIR/PIR"
)

var q int64 = 2147483647

type Server struct {
	Data  matrix.Matrix
	A1    matrix.Matrix
	A2    matrix.Matrix
	HintS matrix.Matrix
	HintC matrix.Matrix
}

func (S *Server) SendHint(args *SingleHintArgs, reply *SingleHintReply) {
	reply.Hint.Multiply(S.Data, S.A1)
}

func (S *Server) SendDoubleHint(args *DoubleHintArgs, reply *DoubleHintReply) {
	preDecomp := matrix.MakeMatrix(S.Data.Rows, S.A1.Columns, 0, q)
	//Need to add decomposition/composition method?
	preDecomp.Multiply(S.Data, S.A1)

	HintS := matrix.Decompose(preDecomp)

	reply.HintC = matrix.MakeMatrix(HintS.Rows, S.A2.Columns, 0, q)
	reply.HintC.Multiply(HintS, S.A2)
}

func (S *Server) SendDoubleAnswer(args *SendDoubleAnswerArgs, reply *SendDoubleAnswerReply) {
	reply.Ans = pir.DoubleAns(S.Data, S.HintS, args.Query, S.A2)
}
