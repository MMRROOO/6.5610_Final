package p2p

import (
	"log"
	"net/http"
	"net/rpc"
	matrix "pir/PIR/Matrix"
	"pir/PIR/PIR"
	"sync"
)

type Peer struct {
	peers  []Endpoint
	mu     sync.Mutex
	me     Endpoint
	Data   matrix.Matrix
	secret matrix.Matrix
	Hashes []byte
	Host   Endpoint
}

var q int64 = 2147483647

func MakePeer(Host Endpoint, Hashes []byte) {
	P := Peer{}
	P.peers = make([]Endpoint, 0)
	P.me = CreateEndpointSelf()
	P.Data = matrix.MakeMatrix(256, 256, 0, q)
	P.secret = matrix.MakeMatrix(256, 1, 1, q)
	P.Host = Host
	P.Hashes = Hashes

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
func CreateEndpointSelf() Endpoint {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
	ownEndpoint := new(Endpoint)
	ownEndpoint.Port = "8080"
	ownEndpoint.ServerAddress = ipAddress
	return *ownEndpoint
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

func (P *Peer) GetFile(server int, index int) []int {

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

// func CheckHash(File matrix.Matrix, Hash byte) bool {
// 	columnArray := File.GetColumn(0)
// 	CHash := sha256.Sum256([]byte(columnArray))

// 	return CHash == Hash
// }

func (P *Peer) PIRAns(args *PIRArgs, reply *PIRReply) {
	reply.Ans = PIR.Ans(P.Data, args.Qu)
	return
}
