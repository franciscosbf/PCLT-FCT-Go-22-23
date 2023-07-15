package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"syscall"
	"time"

	"bitcoin_miner/message"
)

var (
	logger *log.Logger
	cpus   int = runtime.NumCPU()
)

const (
	hdlrWaitTime = 2 * time.Second
	buffSize     = 32
)

// readRawMessage reads data from a given connection.
func readRawMessage(conn *net.TCPConn) ([]byte, error) {
	var (
		n    int
		err  error
		buff = make([]byte, buffSize)
	)

	// Set read deadline to prevent blocking
	conn.SetReadDeadline(time.Now().Add(hdlrWaitTime))

	for err == nil {
		readBuff := make([]byte, buffSize)
		if n, err = conn.Read(readBuff); n > 0 {
			buff = append(buff, readBuff...)
		}
	}

	if err == io.EOF || errors.Is(err, os.ErrDeadlineExceeded) {
		return buff, nil
	}

	return nil, err
}

// serveConn handles a given connection
func serveConn(conn *net.TCPConn) {
	defer conn.Close()

	data, err := readRawMessage(conn)
	if err != nil {
		logger.Printf("error while reading message: %v", err)
		return
	}

	msg, err := message.FromJSON(data)
	if err != nil {
		logger.Printf("invalid message: %s", string(data))
		return
	}

	_ = msg // TODO:
}

// connReceiver behaves as a handler
// waiting to serve a connection.
func connReceiver(
	listener *net.TCPListener,
	closeCh <-chan os.Signal,
) error {
	for {
		// Stops upon some error or keeps rocking.
		select {
		case <-closeCh:
			return nil
		default:
		}

		conn, err := listener.AcceptTCP()
		if err != nil {
			return err
		}

		go serveConn(conn)
	}
}

// listen waits for connections to the server with
// multiple handlers to deal with incoming requests.
func listen(addr *net.TCPAddr) error {
	termCh := make(chan os.Signal, 1)

	// Initializes the listener
	listener, err := net.ListenTCP("tcp4", addr)
	if err != nil {
		return err
	}

	// Marks signCh to receive the following signals:
	signal.Notify(
		termCh, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	return connReceiver(listener, termCh)
}

// startServer runs the server core.
func startServer(port int) error {
	// Builds a tcp addr representation on any host.
	addr, _ := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d", port))

	return listen(addr)
}

func main() {
	// File logger for debugging purposes.
	const (
		name = "server.log"
		flag = os.O_RDWR | os.O_CREATE
		perm = os.FileMode(0666)
	)

	file, err := os.OpenFile(name, flag, perm)
	if err != nil {
		return
	}
	defer file.Close()

	logger = log.New(file, "", log.Lshortfile|log.Lmicroseconds)

	const numArgs = 2
	if len(os.Args) != numArgs {
		fmt.Printf("Usage: ./%s <port>", os.Args[0])
		return
	}

	port, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println("Port must be a number")
		return
	}

	if err := startServer(port); err != nil {
		log.Printf("internal error: %v", err)
	}
}
