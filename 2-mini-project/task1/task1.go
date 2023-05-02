package task1

import (
	"log"
	"sync"
	"time"
)

func Train(id int, track chan<- int, wg *sync.WaitGroup) {
	log.Printf("Train %d waiting to enter the track", id)

	track <- id	

	wg.Done()
}

func Track(done <-chan bool) chan int {
	track := make(chan int)

	go func() {
		for {
			var train int
			select {
			case train = <-track:
				log.Printf("Train %d is inside the track", train)
				time.Sleep(time.Second) // simulates that's inside the track
			case <-done:
				return
			}

			log.Printf("Train %d has left the track", train)
		}
	}()

	return track
}

func Setup(wg *sync.WaitGroup) {
	var sysWg sync.WaitGroup
	done := make(chan bool)

	track := Track(done)
	sysWg.Add(2)
	go Train(1, track, &sysWg)
	go Train(2, track, &sysWg)
	sysWg.Wait()
	done <- true

	wg.Done()
}
