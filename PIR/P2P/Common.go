package p2p

import (
	matrix "pir/PIR/Matrix"
	"pir/PIR/labrpc"
)

type PIRArgs struct {
	Qu matrix.Encryption
}

type PIRReply struct {
	Ans matrix.Matrix
}

type GetFileArgs struct {
	Me *labrpc.ClientEnd
	Qu matrix.Encryption
}

type GetFileReply struct {
	Ans  matrix.Matrix
	Peer *labrpc.ClientEnd
}

type Endpoint struct {
	ServerAddress string
	Port          string
}
