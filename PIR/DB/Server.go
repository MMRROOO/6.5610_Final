package db

type Server struct {
	Data Matrix
	A1   Matrix
	A2   Matrix
}

func (S *Server) Hint(args SingleHintArgs, reply SingleHintReply) {
	reply.Hint.Multiply(S.Data, S.A1)
}

func (S *Server) DoubleHint(Args DoubleHintArgs, reply DoubleHintReply) {
	Hint1 := MakeMatrix(S.Data.Rows, S.A1.Columns, 0, q)
	//Need to add decomposition/composition method?
	Hint1.Multiply(S.Data, S.A1)

	reply.Hint.Multiply(Hint1, S.A2)

}
