package db

import (
	matrix "pir/PIR/Matrix"
	pir "pir/PIR/PIR"
)

type SingleHintReply struct {
	Hint matrix.Matrix
}

type DoubleHintReply struct {
	HintC matrix.Matrix
}

type SingleHintArgs struct {
}

type DoubleHintArgs struct {
}

type SendDoubleAnswerArgs struct {
	Query pir.Dqu
}

type SendDoubleAnswerReply struct {
	Ans pir.Dans
}
