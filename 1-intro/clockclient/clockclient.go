package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"sync/atomic"
	"time"
)

const trials = 4

func main() {
	sleep := flag.Int64("sleep", 1000, "timeout between dials in milliseconds")
	host := flag.String("host", "", "server hostname")
	port := flag.Int("port", 8080, "server port")
	requests := flag.Int64("max", 100, "requests max (greater than 4 trials and less than 6 digits)")
	flag.Parse()

	if *sleep < 1 {
		log.Fatalf("The f*** is this \"timeout\" %v?", *sleep)
	}

	if *requests < trials || *requests > 9999 {
		log.Fatalf("The f*** is this \"requests max\" %v?", *requests)
	}

	failsCh := make(chan bool, trials)
	stopCh := make(chan bool)

	// Keeps track of all fails
	go func() {
		fails := 0

		for {
			<-failsCh // Blocks until receive a fail signal
			if fails++; fails == trials {
				stopCh <- true // We need to stop as soon as possible 
				return
			}
		}
	}()

	var msgCounter atomic.Int64
	msgCounter.Add(1)
	addr := fmt.Sprintf("%v:%v", *host, *port)
	timeout := time.Duration(*sleep)

	// Floods with requests
	for {
		select {
		case <-stopCh: // Signals the party end
			log.Fatal("It's time to give up")
		default: // Keep rockin
		}

		if msgCounter.Load() > *requests {
			return
		}

		go func() {
			conn, dialErr := net.Dial("tcp", addr)	
			if dialErr != nil {
				log.Printf(
					"error while trying to dial with %v: %v",
					addr, dialErr,
				)
				failsCh <- true // Informs that this requester got an error 
				return	
			}
			defer conn.Close()

			// Reads and outputs the received content
			reader := bufio.NewReader(conn)
			content, _, readErr := reader.ReadLine()
			if readErr != nil {
				log.Printf("error while reading received content: %v", readErr)
				failsCh <- true // Does the same thing
				return 
			}
			log.Printf("msg nº %04d: %q", msgCounter.Load(), string(content))
			msgCounter.Add(1)
		}()

		time.Sleep(timeout * time.Millisecond)
	}
}

