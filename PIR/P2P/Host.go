package p2p

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
)

type Host struct {
	Peers []Endpoint
	Me    Endpoint
	// Hashes matrix.Matrix
	// Data   matrix.Matrix
}

// all files in Data must already be split by chunk and a multiple of 256*248
func MakeHost(Data []byte) *Host {
	H := Host{}
	fmt.Printf("before register")
	rpc.Register(&H)
	rpc.HandleHTTP()
	// Listen for requests on port 1234
	l, e := net.Listen("tcp", ":0")

	if e != nil {
		log.Fatal("listen error:", e)
	}
	fmt.Print("Endpoint: ", l, "\n")
	http.Serve(l, nil)
	H.Peers = MakeAllSeedPeers(H, Data)
	// H.Me = CreateEndpointSelf(10000)

	//TODO Make peers so that all data is owned

	// P = MakePeer(,)
	// H.Data = matrix.MakeMatrix(FileSize/int(matrix.DATA_SIZE), NFiles, 0, q)

	// H.Hashes = matrix.MakeMatrix(32, NFiles, 0, q)
	// FillWithHashes(H.Hashes, H.Data)
	return &H
}

func MakeAllSeedPeers(H Host, Data []byte) []Endpoint {
	endpoints := make([]Endpoint, 0)
	fmt.Printf("before Made peer")

	for i := 0; i < len(Data)/(256*248); i++ {
		endpoints = append(endpoints, MakeSeedPeer(Data[i*256*248:(i+1)*256*248], i))
		fmt.Printf("Made peer %d", i)
	}
	return endpoints

}

//	func (H *Host) GetFile(args *GetFileArgs, reply *GetFileReply) {
//		reply.Ans = PIR.Ans(H.Data, args.Qu)
//		H.Peers = append(H.Peers, args.Me)
//		reply.Peer = H.Peers[math.rand()%len(H.Peers)]
//	}
func (H *Host) SendPeers(args *SendPeersArgs, reply *SendPeersReply) error {
	reply.Peers = H.Peers

	H.Peers = append(H.Peers, args.Me)

	return nil
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
