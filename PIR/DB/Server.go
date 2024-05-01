package db

import matrix "pir/PIR/Matrix"

type Server struct {
	Data matrix.Matrix
	A1   matrix.Matrix
	A2   matrix.Matrix
}

func (S *Server) Hint(args SingleHintArgs, reply SingleHintReply) {
	reply.Hint.Multiply(S.Data, S.A1)
}

func (S *Server) DoubleHint(Args DoubleHintArgs, reply DoubleHintReply) {
	Hint1 := matrix.MakeMatrix(S.Data.Rows, S.A1.Columns, 0, q)
	//Need to add decomposition/composition method?
	Hint1.Multiply(S.Data, S.A1)

	reply.Hint.Multiply(Hint1, S.A2)

}
