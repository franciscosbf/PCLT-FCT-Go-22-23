package main

import (
	"fmt"
	"sync"
	"time"
)

type pipeline struct {
	m sync.RWMutex
	doneCh chan int
	toSieveCh chan int
	toMergeCh chan int
	marked map[int]*struct{}
}

var nothing = &struct{}{}

func (p *pipeline) isMarked(v int) bool {
	p.m.RLock()
	defer p.m.RUnlock()
	_, ok := p.marked[v]
	return ok
}

func (p *pipeline) producer() {
	// Shoots values like there's no tomorrow
	for v := 2; ; v++ {
		p.toSieveCh <- v 
		time.Sleep(time.Millisecond * 200)
	}
}

func (p *pipeline) sieve() {
	for {
		v := <-p.toSieveCh
		if p.isMarked(v) {
			continue // We don't give a damn about non primes
		}
		p.toMergeCh <- v

		mult := v
		// Marks all multiples of v
		go func() {
			for vv := 1; ; vv++ {
				v := mult * vv

				if p.isMarked(v) {
					continue // We don't give a damn about non primes
				}

				p.m.Lock()
				p.marked[v] = nothing
				p.m.Unlock()
			
				time.Sleep(time.Millisecond * 200)
			}
		}()
	}	
}

func (p *pipeline) merge() {
	for {
		fmt.Println(<-p.toMergeCh)
	}
}

func main() {
	pipe := &pipeline{
		toSieveCh: make(chan int),
		toMergeCh: make(chan int),
		marked: make(map[int]*struct{}),
	}

	go pipe.producer()
	go pipe.sieve()
	go pipe.merge()

	<-make(chan struct{})
}

