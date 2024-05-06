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

type PIRArgs struct {
	Qu matrix.Encryption
}

type PIRReply struct {
	Ans matrix.Matrix
}

func MakePeer(peers []*labrpc.ClientEnd, me int, knownPeers []int, Host *labrpc.ClientEnd) *Peer {
	P := new(Peer)
	P.peers = peers
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

func CheckHash(File Matrix, Hash [32]byte) bool {
	columnArray := File.GetColumn(0)
	CHash := sha256.Sum256([]byte(columnArray))

	return CHash == Hash
}

func (P *Peer) PIRAns(args *PIRArgs, reply *PIRReply) {
	reply.Ans = Ans(P.Data, args.Qu)
	return
}
