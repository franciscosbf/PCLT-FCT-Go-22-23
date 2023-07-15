package main

import (
	"bitcoin_miner/hash"
	"bitcoin_miner/message"
	"fmt"
)

func sampleHash() uint64 {
	return hash.Hash("The Raven That Refused to Sing (and other stories)", 5443)

}

func sampleReqMsg() *message.Message {
	return message.NewRequest("Sibelius Violin Concerto", 1904, 1905)

}

func sampleResMsg() *message.Message {
	h1 := hash.Hash("Sibelius Violin Concerto", 1904)
	h2 := hash.Hash("Sibelius Violin Concerto", 1905)

	var pow, nonce uint64
	if h1 > h2 {
		pow = h1
		nonce = 1904
	} else {
		pow = h2
		nonce = 1905
	}
	return message.NewResult(pow, nonce, 1904, 1905)
}

func main() {
	fmt.Println("A sample hash", sampleHash())
	fmt.Println("A sample request message", sampleReqMsg())
	fmt.Println("A sample reply message", sampleResMsg())

}
