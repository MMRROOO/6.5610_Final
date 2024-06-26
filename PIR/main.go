package main

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	p2p "pir/PIR/P2P"
	"time"
)

func main() {
	Data := generateData(248 * 256)

	fmt.Printf("generated data %d\n", len(Data)/1024)
	Host := p2p.MakeHost(Data)
	fmt.Print("made host\n")
	hashes := Hashes(Data)
	fmt.Print("generated hashes\n")

	Peer := p2p.MakePeer(Host.Me, hashes, 1)
	fmt.Print(Data[0:1024])
	fmt.Print("made peer\n")
	Peer2 := p2p.MakePeer(Host.Me, hashes, 6)

	// p2p.MakePeer(Host.Me, hashes, 8)

	for !compare(Peer.DesiredChunk, Data[1*1024:2*1024]) {
		fmt.Print("Working1\n")
		time.Sleep(500 * time.Millisecond)
	}

	for !compare(Peer2.DesiredChunk, Data[6*1024:7*1024]) {
		fmt.Print("Working2\n")
		time.Sleep(500 * time.Millisecond)
	}
	fmt.Print("Worked")

	// for {
	// e1 := p2p.Endpoint{ServerAddress: "localhost", Port: "0000"}
	// eList := make([]p2p.Endpoint, 1)
	// eList[0] = e1
	// a := p2p.PeerEndpointToMatrix(eList)
	// b := p2p.MatrixToEndpoint(a)
	// fmt.Print(a)
	// fmt.Print(b)
	// }

	// e := p2p.CreateEndpointSelf()
	// client, err := rpc.DialHTTP("tcp", e.ServerAddress+e.Port)
	// if err != nil {
	// 	log.Fatal("dialing:", err)
	// }
	// call_err := client.Call("Peer.PIRAns", &args, &reply)
	// if call_err != nil {
	// 	log.Fatal("arith error:", call_err)
	// }

}

func compare(a []byte, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func generateData(numBytes int) []byte {
	data := make([]byte, numBytes)

	rand.Read(data)
	return data
}

func Hashes(Data []byte) []byte {
	hashes := make([]byte, (len(Data)/1024)*32)
	for i := 0; i < len(Data)/1024; i++ {
		h := sha256.Sum256(Data[i*1024 : (i+1)*1024])
		for j := 0; j < 32; j++ {
			hashes[i*32+j] = h[j]
		}
	}
	return hashes
}
