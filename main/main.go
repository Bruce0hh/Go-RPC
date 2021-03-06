package main

import (
	gorpc "Go-RPC"
	"Go-RPC/codec"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"
)

func main() {
	addr := make(chan string)
	go startServer(addr)

	conn, _ := net.Dial("tcp", <-addr)
	defer func() { _ = conn.Close() }()

	time.Sleep(time.Second)
	_ = json.NewEncoder(conn).Encode(gorpc.DefaultOption)
	cc := codec.NewGobCodec(conn)

	for i := 0; i < 5; i++ {
		h := &codec.Header{
			ServiceMethod: "Foo.Sum",
			Seq:           uint64(i),
		}

		_ = cc.Write(h, fmt.Sprintf("rpc req %d", h.Seq))
		_ = cc.ReaderHeader(h)
		var reply string
		_ = cc.ReadBody(&reply)
		log.Println("reply: ", reply)
	}
}

func startServer(addr chan string) {
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatal("network error: ", err)
	}
	log.Println("start rpc server on", l.Addr())
	addr <- l.Addr().String()
	gorpc.Accept(l)
}
