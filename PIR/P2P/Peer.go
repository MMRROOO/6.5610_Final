package p2p

import (
	"crypto/sha256"
	"fmt"
	"log"
	"net/http"
	"net/rpc"
	matrix "pir/PIR/Matrix"
	"pir/PIR/PIR"
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
	DesiredFileName int
	DesiredFile     []byte
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

func MakePeer(Host Endpoint, Hashes []byte, FileToDownload int) *Peer {
	P := Peer{}
	P.peers = make([]Endpoint, 0)
	P.me = CreateEndpointSelf(0)
	P.Data = matrix.MakeMatrix(256, 256, 0, q)
	P.secret = matrix.MakeMatrix(256, 1, 1, q)
	P.Host = Host
	P.Hashes = Hashes
	P.DesiredFileName = FileToDownload

	go P.ticker()

	return &P

}

//Data must be in following format:
// total bytes: 248*256 bytes
// each chunk is continuous 4*256 bytes in Data

func MakeSeedPeer(Data []byte, i int) Endpoint {
	P := Peer{}
	fmt.Print("before endpoint")

	P.me = CreateEndpointSelf(i)
	fmt.Print("made endpoint")
	P.Data = matrix.MakeMatrix(256, 256, 0, q)
	FillMatrix(P.Data, Data)
	P.secret = matrix.MakeMatrix(256, 1, 1, q)
	return P.me
}

func FillMatrix(Matrix matrix.Matrix, Data []byte) {
	for i := 0; i < 248; i++ {
		for j := 0; j < 256; j++ {
			Matrix.Set(j, i+8, int64(Data[i*256+j]))
		}
	}
}

var ipAddress string

func handler(w http.ResponseWriter, req *http.Request) {
	ipAddress = req.Header.Get("X-Real-Ip") // Store the IP address from the request header
	if ipAddress == "" {
		ipAddress = req.Header.Get("X-Forwarded-For")
	}
	if ipAddress == "" {
		ipAddress = req.RemoteAddr
	}
}

// TODO: Create Peers own endpoint
func CreateEndpointSelf(i int) Endpoint {
	http.HandleFunc("/"+string(i), handler)
	fmt.Print("1")
	go http.ListenAndServe(":8080", nil)
	fmt.Print("2")
	ownEndpoint := new(Endpoint)

	ownEndpoint.Port = "8080"
	ownEndpoint.ServerAddress = ipAddress
	return *ownEndpoint
}

// TODO: given vector of file names return list of file names it represents
func MatrixToFileNames(M matrix.Matrix) []int {
	return make([]int, 0)
}

// TODO: given 4 matrixes return file data
func FileFromMatrixes(M []matrix.Matrix) []byte {
	return make([]byte, 0)
}

// TODO: given vector of peers return list of file names it represents
func MatrixToPeers(M matrix.Matrix) []int {
	return make([]int, 0)
}

func (P *Peer) GetFileNames(server int) []int {

	qu1 := PIR.Query(4, P.secret)
	args := PIRArgs{Qu: qu1}
	reply := PIRReply{}

	client, err := rpc.DialHTTP("tcp", P.peers[server].ServerAddress+P.peers[server].Port)
	if err != nil {
		log.Fatal("dialing:", err)
	}
	call_err := client.Call("Peer.PIRAns", &args, &reply)
	if call_err != nil {
		log.Fatal("arith error:", call_err)
	}

	FileNames := MatrixToFileNames(PIR.Reconstruct(reply.Ans, P.secret))

	qu2 := PIR.Query(5, P.secret)
	args = PIRArgs{Qu: qu2}
	reply = PIRReply{}

	pir_client, pir_err := rpc.DialHTTP("tcp", P.peers[server].ServerAddress+P.peers[server].Port)
	if pir_err != nil {
		log.Fatal("dialing:", pir_err)
	}
	pir_call_err := pir_client.Call("Peer.PIRAns", &args, &reply)
	if pir_call_err != nil {
		log.Fatal("arith error:", pir_call_err)
	}

	FileNames = append(FileNames, MatrixToFileNames(PIR.Reconstruct(reply.Ans, P.secret))...)

	return FileNames
}

func (P *Peer) GetPeers(server int) []int {

	knownPeers := make([]int, 0)
	for i := 0; i < 4; i++ {
		qu1 := PIR.Query(i, P.secret)
		args := PIRArgs{Qu: qu1}
		reply := PIRReply{}
		pir_client, pir_err := rpc.DialHTTP("tcp", P.peers[server].ServerAddress+P.peers[server].Port)
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

func (P *Peer) GetFile(server int, index int) []byte {

	fileMatrixes := make([]matrix.Matrix, 0)
	for i := 0; i < 4; i++ {
		qu1 := PIR.Query(i, P.secret)
		args := PIRArgs{Qu: qu1}
		reply := PIRReply{}
		pir_client, pir_err := rpc.DialHTTP("tcp", P.peers[server].ServerAddress+P.peers[server].Port)
		if pir_err != nil {
			log.Fatal("dialing:", pir_err)
		}
		pir_call_err := pir_client.Call("Peer.PIRAns", &args, &reply)
		if pir_call_err != nil {
			log.Fatal("arith error:", pir_call_err)
		}

		fileMatrixes = append(fileMatrixes, PIR.Reconstruct(reply.Ans, P.secret))
	}

	return FileFromMatrixes(fileMatrixes)
}

func CheckHash(File matrix.Matrix, Hash [32]byte) bool {
	columnArray := intToByte(File.GetColumn(0))

	CHash := sha256.Sum256([]byte(columnArray))

	return CHash == Hash
}

func (P *Peer) PIRAns(args *PIRArgs, reply *PIRReply) {
	reply.Ans = PIR.Ans(P.Data, args.Qu)
	return
}

func intToByte(s []int) []byte {
	r := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		r[i] = byte(s[i])
	}
	return r
}

func (P *Peer) ticker() {

	for {
		args := SendPeersArgs{}
		args.Me = P.me
		reply := SendPeersReply{}

		client, err := rpc.DialHTTP("tcp", P.Host.ServerAddress+P.Host.Port)
		if err != nil {
			log.Fatal("dialing:", err)
		}
		call_err := client.Call("Host.SendPeers", &args, &reply)
		if call_err != nil {
			log.Fatal("arith error:", call_err)
		}

		P.peers = reply.Peers

		for i := 0; i < len(P.peers); i++ {
			FileNames := P.GetFileNames(i)
			c := contains(FileNames, P.DesiredFileName)

			if c != -1 {
				P.DesiredFile = P.GetFile(i, c)
			}

		}
	}
}
