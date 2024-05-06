package p2p

import (
	matrix "pir/PIR/Matrix"
	"pir/PIR/PIR"
	"pir/PIR/labrpc"

	"sync"
)

type Host struct {
	Peers  []*labrpc.ClientEnd
	mu     sync.Mutex
	Hashes matrix.Matrix
	Data   matrix.Matrix
}

func MakeHost(NFiles int, HashSize int, FileSize int) *Host {
	H := new(Host)
	H.Peers = make([]*labrpc.ClientEnd, 0)

	//TODO Make peers so that all data is owned
	// P = MakePeer(,)
	H.Data = matrix.MakeMatrix(FileSize/int(matrix.DATA_SIZE), NFiles, 0, q)
	H.Hashes = matrix.MakeMatrix(32, NFiles, 0, q)
	// FillWithHashes(H.Hashes, H.Data)
	return H
}

func (H *Host) GetFile(args *GetFileArgs, reply *GetFileReply) {
	reply.Ans = PIR.Ans(H.Data, args.Qu)
	H.Peers = append(H.Peers, args.Me)
	reply.Peer = H.Peers[math.rand()%len(H.Peers)]
}

// func FillWithHashes(Hash matrix.Matrix, Data matrix.Matrix) {
// 	for c := 0; c < Data.Columns; c++ {
// 		columnArray := Data.GetColumn(c)
// 		CHash := sha256.Sum256([]byte(columnArray))
// 		for r := 0; r < Hash.Rows; r++ {
// 			Hash.Set(r, c, CHash[r])
// 		}
// 	}
// }

// func (H *Host) GetHash(args *GetFileArgs, reply *GetFileReply) {
// 	reply.Ans = PIR.Ans(H.Hashes, args.Qu)
// 	H.knownPeers = append(H.knownPeers, args.me)
// 	reply.Peer = H.knownPeers[rand()%len(H.knownPeers)]
// }
