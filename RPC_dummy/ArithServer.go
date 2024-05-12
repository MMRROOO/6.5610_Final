/**
* ArithServer
 */
package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/rpc"
	"time"
)

type Values struct {
	A, B int
}
type Quotient struct {
	Quo, Rem int
}
type Arith int

func (t *Arith) Multiply(args *Values, reply *int) error {
	*reply = args.A * args.B
	return nil
}
func (t *Arith) Divide(args *Values, quo *Quotient) error {
	if args.B == 0 {
		return errors.New("divide by zero")
	}
	quo.Quo = args.A / args.B
	quo.Rem = args.A % args.B
	return nil
}
func main() {
	arith := new(Arith)
	fmt.Print("made\n")
	rpc.Register(arith)
	fmt.Print("registered\n")
	rpc.HandleHTTP()
	fmt.Print("handle\n")
	go callstuff()
	err := http.ListenAndServe(":1234", nil)
	fmt.Print("listen and serve\n")
	if err != nil {
		fmt.Println(err.Error())
	}
}

func callstuff() {
	client, err := rpc.DialHTTP("tcp", "localhost"+":1234")
	fmt.Print("dialing")
	if err != nil {
		log.Fatal("dialing:", err)
	}
	time.Sleep(10 * time.Millisecond)
	args := Values{A: 10, B: 10}
	var reply int

	ok := client.Call("Arith.Multiply", &args, &reply)
	fmt.Print("calling")

	if ok != nil {

	}

	fmt.Print(reply)
}
