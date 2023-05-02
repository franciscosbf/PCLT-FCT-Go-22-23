package main

import "fmt"

type MsgType int

const (
	Kill MsgType = iota
	NextPrime
)

type PrimeMsg struct {
	mtype MsgType
	reply chan int
}

func Producer(ch chan<- int) {
	for i := 2; ; i++ {
		ch <- i
	}
}

func Sieve(in <-chan int, out chan<- int, prime int) {
	for {
		n := <-in
		if n % prime != 0 {
			out <- n
		}
	}
}

func Assembler(in <-chan int) chan *PrimeMsg {
	res := make(chan *PrimeMsg)

	go func() {
		for {
			msg := <-res
			switch msg.mtype {
			case Kill:
				return
			case NextPrime:
				prime := <-in
				msg.reply <- prime
				newCh := make(chan int)
				go Sieve(in, newCh, prime)
				in = newCh
			}
		}
	}()

	return res
}

func SieveMachine() chan *PrimeMsg {
	ch := make(chan int)
	go Producer(ch)
	return Assembler(ch)
}

func main() {
	resCh := SieveMachine()
	answerCh := make(chan int)

	for i := 0; i < 1000; i++ {
		resCh <- &PrimeMsg{
			mtype: NextPrime,
			reply: answerCh,
		}

		fmt.Printf("%d ", <-answerCh)
	}

	resCh <- &PrimeMsg{
		mtype: Kill,
		reply: nil,
	}

	fmt.Println()
}

