package task2

import (
	"log"
	"sync"
	"time"
)

type Track struct {
	left chan int
	right chan int
}

func PassTrough(train int) {
	log.Printf("Train %d is inside the track", train)
	time.Sleep(time.Second) // simulates that's inside the track
	log.Printf("Train %d has left the track", train)
}

func (t *Track) EnterLeft(train int) {
	log.Printf("Train %d is waiting on the left side", train)
	t.left <- train 
}

func (t *Track) EnterRight(train int) {
	log.Printf("Train %d is waiting on the right side", train)
	t.right <- train 
}

func (t *Track) Inside(wg *sync.WaitGroup) {
	side := t.left

	for {
		var train int

		select {
		case train = <-side:
			PassTrough(train)
			wg.Done()
		case <-time.After(time.Second):
		}

		if side == t.left {
			side = t.right
		} else {
			side = t.left
		}
	}
}

func Setup(wg *sync.WaitGroup) {
	var sysWg sync.WaitGroup

	track := &Track{
		left: make(chan int),
		right: make(chan int),
	}

	sysWg.Add(6)
	go track.Inside(&sysWg)
	go track.EnterLeft(1)
	go track.EnterRight(2)
	go track.EnterLeft(3)
	go track.EnterRight(4)
	go track.EnterLeft(5)
	go track.EnterRight(6)
	sysWg.Wait()

	wg.Done()
}
