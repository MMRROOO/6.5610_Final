package p2p

import (
	"crypto/sha256"
	matrix "pir/PIR/Matrix"
	"sync"
)

type Host struct {
	Peers      []*labrpc.ClientEnd
	knownPeers []int
	mu         sync.Mutex
	Hashes     matrix.Matrix
	Data       matrix.Matrix
}

func MakeHost(peers []*labrpc.ClientEnd, NFiles int, HashSize int, FileSize int) *Host {
	H := new(Host)
	H.Peers = peers
	H.Data = matrix.MakeMatrix(FileSize/int(matrix.DATA_SIZE), NFiles, 0, q)
	H.Hashes = matrix.MakeMatrix(32, NFiles, 0, q)
	FillWithHashes(H.Hashes, H.Data)
	return H
}

func (H *Host) GetFile(args *GetFileArgs, reply *GetFileReply) {
	reply.Ans = Ans(H.Data, args.qu)
	H.knownPeers = append(H.knownPeers, args.me)
	reply.Peer = H.knownPeers[rand()%len(H.knownPeers)]
}

func FillWithHashes(Hash matrix.Matrix, Data matrix.Matrix) {
	for c := 0; c < Data.Columns; c++ {
		columnArray := Data.GetColumn(c)
		CHash := sha256.Sum256([]byte(columnArray))
		for r := 0; r < Hash.Rows; r++ {
			Hash.Set(r, c, CHash[r])
		}
	}
}

func (H *Host) GetHash(args *GetFileArgs, reply *GetFileReply) {
	reply.Ans = Ans(H.Hashes, args.qu)
	H.knownPeers = append(H.knownPeers, args.me)
	reply.Peer = H.knownPeers[rand()%len(H.knownPeers)]
}
