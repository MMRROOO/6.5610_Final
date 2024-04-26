package main

import (
	"sync"
)

type Peer struct {
	peers      []*labrpc.ClientEnd
	knownPeers []int
	mu         sync.Mutex
	me         int
	Data       Matrix
	secret     Matrix
	Host       *labrpc.ClientEnd
}

type PIRArgs struct {
	Qu encryption
}

type PIRReply struct {
	Ans Matrix
}

func MakePeer(peers []*labrpc.ClientEnd, me int, knownPeers []int, Host *labrpc.ClientEnd) *Peer {
	P := new(Peer)
	P.peers = peers
	P.me = metrics
	P.Data = MakeMatrix(256, 256, 0, q)
	P.secret = MakeMatrix(256, 1, 1, q)
	P.knownPeers = knownPeers
	P.Host = Host

	return P
}

// TODO: given vector of file names return list of file names it represents
func MatrixToFileNames(M Matrix) []int {}

// TODO: given 4 matrixes return file data
func FileFromMatrixes(M []Matrix) []int {}

// TODO: given vector of peers return list of file names it represents
func MatrixToPeers(M Matrix) []int {}

func (P *Peers) GetFileNames(server int) []int {

	qu1 = Query(4, P.secret)
	args := PIRArgs{Qu: qu1}
	reply := PIRReply
	ok := P.peers[server].Call("Peer.PIRAns", &args, &reply)

	FileNames := MatrixToFileNames(reconstruct(reply.Ans, P.secret))

	qu2 = Query(5, P.secret)
	args = PIRArgs{Qu: qu1}
	reply = PIRReply
	ok = P.peers[server].Call("Peer.PIRAns", &args, &reply)

	FileNames = append(Filenames, MatrixToFileNames(reconstruct(reply.Ans, P.secret)))

	return FileNames
}

func (P *Peers) GetPeers(server int) []int {
	secret := MakeMatrix(N, 1, 1)

	knownPeers := []int
	for i := 0; i < 4; i++ {
		qu1 = Query(i, secret)
		args := PIRArgs{Qu: qu1}
		reply := PIRReply
		ok := P.peers[server].Call("Peer.PIRAns", &args, &reply)

		knownPeers = append(knownPeers, MatrixToPeers(reconstruct(reply.Ans, P.secret)))
	}

	return knownPeers
}

func (P *Peers) GetFile(server int, index int) []int {

	fileMatrixes := []Matrix
	for i := 0; i < 4; i++ {
		qu1 = Query(i, secret)
		args := PIRArgs{Qu: qu1}
		reply := PIRReply
		ok := P.peers[server].Call("Peer.PIRAns", &args, &reply)

		fileMatrixes = append(fileMatrixes, reconstruct(reply.Ans, P.secret))
	}

	return FileFromMatrixes(fileMatrixes)
}

func (P *Peer) PIRAns(args *PIRArgs, reply *PIRReply) {
	reply.Ans = Ans(P.Data, args.Qu)
	return
}
