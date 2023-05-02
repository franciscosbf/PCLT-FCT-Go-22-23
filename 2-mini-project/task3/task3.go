package task3

import (
	"math/rand"
	"log"
	"sync"
	"time"
)

type Track struct {
	left chan int
	right chan int
}

func PassTrough(train int, wgProg *sync.WaitGroup, wgStack *sync.WaitGroup) {
	defer wgProg.Done()
	defer wgStack.Done() // called first

	log.Printf("Train %d is inside the track", train)
	deviation := rand.Intn(100) + 100 // simulation purposes
	time.Sleep(time.Duration(deviation) * time.Millisecond)
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

func (t *Track) TrackSystem(wg *sync.WaitGroup) {
	side := t.left

	for {
		timeout := time.After(time.Second)

		var wgStack sync.WaitGroup

		for {
			select {
			// tries to stack as much trains as possible in a second
			case train := <-side: 
				wgStack.Add(1)
				go PassTrough(train, wg, &wgStack)
			case <-timeout:
				goto wait // stops accepting from the current side
			}
		}

		// Let then all pass trough the track 
		// before switching to the oposite side
		wait: wgStack.Wait() 

		// Switch sides
		if side == t.left {
			side = t.right
		} else {
			side = t.left
		}
	}
}

func Setup(wg *sync.WaitGroup) {
	var trainsWg sync.WaitGroup

	rand.Seed(time.Now().UnixNano()) // simulation purposes

	track := &Track{
		left: make(chan int),
		right: make(chan int),
	}

	go track.TrackSystem(&trainsWg)
	trainsWg.Add(6)
	go track.EnterLeft(1)
	go track.EnterRight(2)
	go track.EnterLeft(3)
	go track.EnterRight(4)
	go track.EnterLeft(5)
	go track.EnterRight(6)
	trainsWg.Wait()

	wg.Done()
}
