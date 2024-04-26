package main

import (
	"crypto/sha256"
	"sync"
)

type Host struct {
	Peers      []*labrpc.ClientEnd
	knownPeers []int
	mu         sync.Mutex
	Hashes     Matrix
	Data       Matrix
}

func MakeHost(peers []*labrpc.ClientEnd, NFiles int, HashSize int, FileSize int) *Host {
	H := new(Host)
	H.Peers = peers
	H.Data = MakeMatrix(FileSize/DATA_SIZE, NFiles, 0, q)
	H.Hashes = MakeMatrix(HashSize/DATA_SIZE, NFiles, 0, q)
	FillWithHashes(H.Hashes, H.Data)
	return H
}

func (H *Host) GetFile(args *GetFileArgs, reply *GetFileReply) {
	reply.Ans = Ans(H.Data, args.qu)
	H.knownPeers = append(H.knownPeers, args.me)
	reply.Peer = H.knownPeers[rand()%len(H.knownPeers)]
}

func FillWithHashes(Hash Matrix, Data Matrix) {
	for c := 0; c < Data.Columns; c++ {
		columnArray := Data.GetColumn(c)
		Hash = sha256.Sum256([]byte(columnArray))
	}
}

func (H *Host) GetHash(args *GetFileArgs, reply *GetFileReply) {
	reply.Ans = Ans(H.Hashes, args.qu)
	H.knownPeers = append(H.knownPeers, args.me)
	reply.Peer = H.knownPeers[rand()%len(H.knownPeers)]
}
