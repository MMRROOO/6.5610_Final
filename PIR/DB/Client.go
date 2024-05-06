package db

import (
	matrix "pir/PIR/Matrix"
	pir "pir/PIR/PIR"
	"pir/PIR/labrpc"
)

type Client struct {
	server  *labrpc.ClientEnd
	HintC   matrix.Matrix
	SecretS matrix.Matrix
	SecretC matrix.Matrix
	A1      matrix.Matrix
	A2      matrix.Matrix
}

func (C *Client) GetSingleHint() {
	args := DoubleHintArgs{}
	reply := DoubleHintReply{}
	ok := C.server.Call("Server.SendDoubleHint", &args, &reply)
	if ok {
		C.HintC = reply.HintC
	}

}

func (C *Client) MakeSingleQuery(r int, c int) {
}

func (C *Client) MakeDoubleQuery(r int, c int) pir.Dqu {
	return (pir.DoubleQuery(r, c, C.SecretC, C.SecretS, C.A1, C.A2))

}
