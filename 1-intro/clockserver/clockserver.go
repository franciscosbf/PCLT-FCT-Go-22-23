package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

func main() {
	port := flag.Int("port", 8080, "server port")
	addr := fmt.Sprintf(":%v", *port)

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		// Handles asynchronously each request
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()
	for {
		_, err := io.WriteString(conn, time.Now().Format(time.RFC1123)+"\n")
		if err != nil {
			return
		}
		time.Sleep(1 * time.Second)
	}
}
