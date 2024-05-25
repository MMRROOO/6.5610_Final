package p2p

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"net/rpc"
	matrix "pir/PIR/Matrix"
	PIR "pir/PIR/PIR"
	"pir/PIR/helper"
	"sync"
)

type Peer struct {
	peers           []Endpoint
	mu              sync.Mutex
	me              Endpoint
	Data            matrix.Matrix
	secret          matrix.Matrix
	Hashes          []byte
	Host            Endpoint
	DesiredChunkName int
	DesiredChunk     []byte
	OwnedChunks     []int
}

func contains(s []int, e int) int {
	for i, a := range s {
		if a == e {
			return i
		}
	}
	return -1
}

var q int64 = 2147483647

func MakePeer(Host Endpoint, Hashes []byte, ChunkToDownload int) *Peer {
	P := Peer{}

	// Listen for requests on port 1234
	// l, e := net.Listen("tcp", "0")
	// fmt.Print("Endpoint: ", l.Addr().String(), "\n")
	// if e != nil {
	// 	log.Fatal("listen error:", e)
	// }
	// http.Serve(l, nil)

	P.peers = make([]Endpoint, 0)
	P.me = CreateEndpointSelf(&P)
	P.Data = matrix.MakeMatrix(256, 256, 0, q)
	P.secret = matrix.MakeMatrix(16, 1, 1, q)
	P.Host = Host
	P.Hashes = Hashes
	P.DesiredChunkName = ChunkToDownload
	P.OwnedChunks = make([]int, 0)

	go P.ticker()

	return &P

}

//Data must be in following format:
// total bytes: 248*256 bytes
// each chunk is continuous 4*256 bytes in Data

func MakeSeedPeer(Data []byte, i int) Endpoint {
	P := Peer{}
	fmt.Print("before endpoint\n")

	P.me = CreateEndpointSelf(&P)
	fmt.Print("made endpoint\n")
	P.Data = matrix.MakeMatrix(256, 256, 0, q)
	chunks := make([]int, 128)
	for f := i; f < i+62; f++ {
		chunks[f-i] = f
	}

	firstFileColumn := helper.FileNamestoMatrices(chunks)
	// for f := i + 128; f < i+248; f++ {
	// 	files[f-(i+128)] = f
	// }

	// for f := i + 248; f < i+256; f++ {
	// 	files[f-(i+128)] = 256 * 255
	// }

	// secondFileColumn := helper.FileNamestoMatrices(files)

	P.Data.CopyColumn(firstFileColumn, 4)
	// P.Data.CopyColumn(secondFileColumn, 5)

	FillMatrix(P.Data, Data)
	P.secret = matrix.MakeMatrix(16, 1, 1, q)
	return P.me
}

func FillMatrix(Matrix matrix.Matrix, Data []byte) {
	for i := 0; i < 248; i++ {
		for j := 0; j < 256; j++ {
			Matrix.Set(j, i+8, int64(Data[i*256+j]))
		}
	}
}

var ipAddress string = "localhost"

func handler(w http.ResponseWriter, req *http.Request) {
	/*
		get its own IP address so it can send to tracker
	*/
	ipAddress = req.Header.Get("X-Real-Ip") // Store the IP address from the request header
	if ipAddress == "" {
		ipAddress = req.Header.Get("X-Forwarded-For")
	}
	if ipAddress == "" {
		ipAddress = req.RemoteAddr
	}
	ipAddress = "localhost"

}
func nrand() int64 { //secure implementation
	max := big.NewInt(int64(1) << 62)
	bigx, _ := rand.Int(rand.Reader, max)
	x := bigx.Int64()
	return x
}

// TODO: Create Peers own endpoint
func CreateEndpointSelf(P *Peer) Endpoint {
	/*
		i is ID of the peer (unique)
	*/

	port := fmt.Sprint(nrand()%1000 + 1000)
	ownEndpoint := Endpoint{}

	ownEndpoint.Port = port
	ownEndpoint.ServerAddress = ipAddress
	fmt.Print(ownEndpoint, "\n")
	go RegisterWithEndpoint(P, ownEndpoint)

	return ownEndpoint
}

func RegisterWithEndpoint(P *Peer, e Endpoint) {
	rpc.Register(P)
	http.HandleFunc("/"+fmt.Sprint(e.Port), handler)
	address := e.ServerAddress + ":" + e.Port
	err := http.ListenAndServe(address, nil)
	if err != nil {
		fmt.Print("network error ", err, "\n")
	}
}

// TODO: given vector of file names return list of file names it represents
func MatrixToChunkNames(M matrix.Matrix) []int {
	return helper.MatrixToFileNames(M)
}

// TODO: given 4 matrixes return file data
func ChunkFromMatrices(M []matrix.Matrix) []byte {
	return helper.MatrixtoFileChunk(M)
}

// TODO: given vector of peers return list of file names it represents
func MatrixToPeers(M matrix.Matrix) []int {
	return make([]int, 0)
}

func (P *Peer) GetChunkNames(server int) []int {
	qu1 := PIR.Query(4, P.secret)
	args := PIRArgs{Qu: qu1}
	reply := PIRReply{}

	client, err := rpc.DialHTTP("tcp", P.peers[server].ServerAddress+":"+P.peers[server].Port)
	if err != nil {
		log.Fatal("dialing:", err)
	}
	call_err := client.Call("Peer.PIRAns", &args, &reply)
	if call_err != nil {
		log.Fatal("arith error:", call_err)
	}

	ChunkNames := MatrixToChunkNames(PIR.Reconstruct(reply.Ans, P.secret))

	qu2 := PIR.Query(5, P.secret)
	args = PIRArgs{Qu: qu2}
	reply = PIRReply{}

	pir_client, pir_err := rpc.DialHTTP("tcp", P.peers[server].ServerAddress+":"+P.peers[server].Port)
	if pir_err != nil {
		log.Fatal("dialing:", pir_err)
	}
	pir_call_err := pir_client.Call("Peer.PIRAns", &args, &reply)
	if pir_call_err != nil {
		log.Fatal("arith error:", pir_call_err)
	}

	ChunkNames = append(ChunkNames, MatrixToChunkNames(PIR.Reconstruct(reply.Ans, P.secret))...)

	return ChunkNames
}

func (P *Peer) GetPeers(server int) []int {

	knownPeers := make([]int, 0)
	for i := 0; i < 4; i++ {
		qu1 := PIR.Query(i, P.secret)
		args := PIRArgs{Qu: qu1}
		reply := PIRReply{}
		pir_client, pir_err := rpc.DialHTTP("tcp", P.peers[server].ServerAddress+":"+P.peers[server].Port)
		if pir_err != nil {
			log.Fatal("dialing:", pir_err)
		}
		pir_call_err := pir_client.Call("Peer.PIRAns", &args, &reply)
		if pir_call_err != nil {
			log.Fatal("arith error:", pir_call_err)
		}

		knownPeers = append(knownPeers, MatrixToPeers(PIR.Reconstruct(reply.Ans, P.secret))...)
	}

	return knownPeers
}

//asks for chunk at index from server
func (P *Peer) GetChunk(server int, index int) []byte {

	fileMatrixes := make([]matrix.Matrix, 0)
	for i := 0; i < 4; i++ {

		qu1 := PIR.Query(index*4+i+8, P.secret)
		args := PIRArgs{Qu: qu1, Me: P.me}

		reply := PIRReply{}
		pir_client, pir_err := rpc.DialHTTP("tcp", P.peers[server].ServerAddress+":"+P.peers[server].Port)
		if pir_err != nil {
			log.Fatal("dialing:", pir_err)
		}
		pir_call_err := pir_client.Call("Peer.PIRAns", &args, &reply)
		if pir_call_err != nil {
			log.Fatal("arith error:", pir_call_err)
		}

		fileMatrixes = append(fileMatrixes, PIR.Reconstruct(reply.Ans, P.secret))
	}
	file := ChunkFromMatrices(fileMatrixes)
	return file
}

func CheckHash(chunk matrix.Matrix, Hash [32]byte) bool {
	columnArray := intToByte(chunk.GetColumn(0))

	CHash := sha256.Sum256([]byte(columnArray))
	return CHash == Hash
}

func (P *Peer) PIRAns(args *PIRArgs, reply *PIRReply) error {
	P.addPeerToMatrix(args.Me, len(P.peers))
	reply.Ans = PIR.Ans(P.Data, args.Qu)
	return nil
}

func intToByte(s []int) []byte {
	r := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		r[i] = byte(s[i])
	}
	return r
}

//puts bytes into internal matrix data structure
func (P *Peer) PutChunk(chunk []byte, index int) {
	P.Data.PutFile(helper.FileChunkToMatrix(chunk), index)
}

//puts flile name into index in interal matrix data stucture
func (P *Peer) PutChunkName(chunkName int, index int) {
	P.Data.Set((index%128)*2, index/128+4, int64(chunkName/256))
	P.Data.Set((index%128)*2+1, index/128+4, int64(chunkName%256))
}

func (P *Peer) GetNewPeer(e Endpoint) {
	endpoints := make([]Endpoint, 0)
	for i := 0; i < 4; i++ {
		fmt.Print(i, "index\n")

		qu1 := PIR.Query(i, P.secret)
		args := PIRArgs{Qu: qu1}
		reply := PIRReply{}
		pir_client, pir_err := rpc.DialHTTP("tcp", e.ServerAddress+":"+e.Port)
		if pir_err != nil {
			log.Fatal("dialing:", pir_err)
		}
		pir_call_err := pir_client.Call("Peer.PIRAns", &args, &reply)
		if pir_call_err != nil {
			log.Fatal("arith error:", pir_call_err)
		}
		tmp := MatrixToEndpoint(PIR.Reconstruct(reply.Ans, P.secret))
		fmt.Print(PIR.Reconstruct(reply.Ans, P.secret), e, "\n")

		endpoints = append(endpoints, tmp...)
	}

	if len(endpoints) == 0 {
		return
	}
	newRandomPeer := endpoints[int(nrand())%len(endpoints)]
	fmt.Print(newRandomPeer, "\n")

	P.addPeerToMatrix(newRandomPeer, len(P.peers))

	P.peers = append(P.peers, newRandomPeer)

}

func (P *Peer) addPeerToMatrix(e Endpoint, index int) {
	if e.Port == P.me.Port && e.ServerAddress == P.me.ServerAddress {
		return
	}
	for i := 0; i < len(P.peers); i++ {
		if e.Port == P.peers[i].Port && e.ServerAddress == P.peers[i].ServerAddress {
			return
		}
	}
	col := index / (256 / 16)
	dataCol := P.Data.GetColumnMatrix(col)
	endpoints := MatrixToEndpoint(dataCol)
	endpoints = append(endpoints, e)
	M := PeerEndpointToMatrix(endpoints)

	P.Data.CopyColumn(M, col)
}

func (P *Peer) ticker() {
	args := SendPeersArgs{}
	args.Me = P.me
	reply := SendPeersReply{}

	client, err := rpc.DialHTTP("tcp", P.Host.ServerAddress+":"+P.Host.Port)
	if err != nil {
		log.Fatal("dialing:", err)
	}
	call_err := client.Call("Host.SendPeers", &args, &reply)
	if call_err != nil {
		log.Fatal("arith error:", call_err)
	}
	for i := 0; i < len(reply.Peers); i++ {
		P.addPeerToMatrix(reply.Peers[i], i)
	}

	P.peers = reply.Peers
	done := false
	for {

		// fmt.Print(P.peers[0], "\n")

		for i := 0; i < len(P.peers); i++ {

			ChunkNames := P.GetChunkNames(i)
			c := contains(ChunkNames, P.DesiredChunkName)

			if c != -1 && done == false {
				P.DesiredChunk = P.GetChunk(i, c)
				done = true
			} else {
				randChunkidx := int(nrand()) % 62
				randChunkName := ChunkNames[randChunkidx]
				newChunk := P.GetChunk(i, randChunkidx)
				P.mu.Lock()
				P.OwnedChunks = append(P.OwnedChunks, randChunkName)
				P.PutChunk(newChunk, len(P.OwnedChunks))
				P.PutChunkName(randChunkName, len(P.OwnedChunks))
				P.mu.Unlock()
			}


			if nrand()%10 == 0 {
				fmt.Print("me: ", P.me, " peers: ", P.peers, "\n")

				P.GetNewPeer(P.peers[i])
				fmt.Print("me: ", P.me, " peers: ", P.peers, "\n")

			}

		}
	}
}
