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

	Peer := p2p.MakePeer(Host.Me, hashes, 0)
	fmt.Print(Data[0:1024])
	fmt.Print("made peer\n")

	for !compare(Peer.DesiredFile, Data[0:1024]) {
		fmt.Print("Working\n")
		time.Sleep(500 * time.Millisecond)
	}
	fmt.Print("Worked")

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
