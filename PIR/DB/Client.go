package db

import matrix "pir/PIR/Matrix"

type Client struct {
	server *Server
}

var HintS matrix.Matrix
var HintC matrix.Matrix
var Hint matrix.Matrix
var singleHint = SingleHintReply{Hint: Hint}

var hasHint = false

func GetSingleHint(server *Server) {
	if hasHint {
		server.Hint(singleHint)
	}

}

func (C *Client) MakeSingleQuery(r int, c int) {

}

func (C *Client) MakeDoubleQuery(r int, c int) {

}
