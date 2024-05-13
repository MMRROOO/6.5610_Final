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
	Me Endpoint
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

type Endpoint struct {
	ServerAddress string
	Port          string
}

// 16 bytes are reserved for the address
// well want to encode ip adress for for local POC we simply encode the port
func PeerEndpointToMatrix(e []Endpoint) matrix.Matrix {
	M := matrix.MakeMatrix(256, 1, 0, 2147483647)
	for i := 0; i < 256/16; i++ {
		if len(e) > i {
			Port := []byte(e[i].Port)
			for j := 0; j < len(Port); j++ {
				M.Set(i*16+j, 0, int64(Port[j]))
			}
		}
	}
	return M
}

func MatrixToEndpoint(M matrix.Matrix) []Endpoint {
	res := make([]Endpoint, 0)
	for i := 0; i < 256/16; i++ {
		curEnd := Endpoint{}
		curEnd.ServerAddress = "localhost"
		Port := make([]byte, 0)
		validport := false
		for j := 0; j < 16; j++ {
			if M.Get(i*16+j, 0) != 0 {
				Port = append(Port, byte(M.Get(i*16+j, 0)))
				validport = true
			}
		}
		if validport {
			curEnd.Port = string(Port[:])
			res = append(res, curEnd)
		}
	}
	return res
}
