package p2p

import (
	matrix "pir/PIR/Matrix"
)

type SendPeersArgs struct {
	Me Endpoint
}

type SendPeersReply struct {
	Peers []Endpoint
}
type PIRArgs struct {
	Qu matrix.Encryption
}

type PIRReply struct {
	Ans matrix.Matrix
}

type GetFileArgs struct {
	Me *Endpoint
	Qu matrix.Encryption
}

type GetFileReply struct {
	Ans  matrix.Matrix
	Peer *Endpoint
}

type  Endpoint struct {
	ServerAddress string
	Port          string
}
