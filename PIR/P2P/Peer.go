package p2p

import (
	"crypto/sha256"
	matrix "pir/PIR/Matrix"
	"pir/PIR/PIR"
	"pir/PIR/labrpc"
	"sync"
)

type Peer struct {
	peers      []*labrpc.ClientEnd
	knownPeers []int
	mu         sync.Mutex
	me         int
	Data       matrix.Matrix
	secret     matrix.Matrix
	Hashes     [][32]byte
	Host       *labrpc.ClientEnd
}

var q int64 = 2147483647

func MakePeer(me int, knownPeers []int, Host *labrpc.ClientEnd) *Peer {
	P := Peer{}
	P.peers = make([]*labrpc.ClientEnd, 0)
	P.me = me
	P.Data = matrix.MakeMatrix(256, 256, 0, q)
	P.secret = matrix.MakeMatrix(256, 1, 1, q)
	P.knownPeers = knownPeers
	P.Host = Host

	return P
}

// TODO: given vector of file names return list of file names it represents
func MatrixToFileNames(M matrix.Matrix) []int {}

// TODO: given 4 matrixes return file data
func FileFromMatrixes(M []matrix.Matrix) []int {}

// TODO: given vector of peers return list of file names it represents
func MatrixToPeers(M matrix.Matrix) []int {}

func (P *Peer) GetFileNames(server int) []int {

	qu1 := PIR.Query(4, P.secret)
	args := PIRArgs{Qu: qu1}
	reply := PIRReply{}
	ok := P.peers[server].Call("Peer.PIRAns", &args, &reply)

	FileNames := MatrixToFileNames(PIR.Reconstruct(reply.Ans, P.secret))

	qu2 := PIR.Query(5, P.secret)
	args = PIRArgs{Qu: qu2}
	reply = PIRReply{}
	ok = P.peers[server].Call("Peer.PIRAns", &args, &reply)

	FileNames = append(FileNames, MatrixToFileNames(PIR.Reconstruct(reply.Ans, P.secret))...)

	return FileNames
}

func (P *Peer) GetPeers(server int) []int {

	knownPeers := make([]int, 0)
	for i := 0; i < 4; i++ {
		qu1 := PIR.Query(i, P.secret)
		args := PIRArgs{Qu: qu1}
		reply := PIRReply{}
		ok := P.peers[server].Call("Peer.PIRAns", &args, &reply)

		knownPeers = append(knownPeers, MatrixToPeers(PIR.Reconstruct(reply.Ans, P.secret))...)
	}

	return knownPeers
}

func (P *Peer) GetFile(server int, index int) []int {

	fileMatrixes := make([]matrix.Matrix, 0)
	for i := 0; i < 4; i++ {
		qu1 := PIR.Query(i, P.secret)
		args := PIRArgs{Qu: qu1}
		reply := PIRReply{}
		ok := P.peers[server].Call("Peer.PIRAns", &args, &reply)

		fileMatrixes = append(fileMatrixes, PIR.Reconstruct(reply.Ans, P.secret))
	}

	return FileFromMatrixes(fileMatrixes)
}

func CheckHash(File matrix.Matrix, Hash [32]byte) bool {
	columnArray := File.GetColumn(0)
	CHash := sha256.Sum256([]byte(columnArray))

	return CHash == Hash
}

func (P *Peer) PIRAns(args *PIRArgs, reply *PIRReply) {
	reply.Ans = PIR.Ans(P.Data, args.Qu)
	return
}
