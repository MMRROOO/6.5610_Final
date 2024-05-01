package db

import matrix "pir/PIR/Matrix"

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
