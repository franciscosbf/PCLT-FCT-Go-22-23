package main

import (
	"math/big"
	"math/rand"
	"testing"
	"time"
)

func TestSieve(t *testing.T) {
	primes := SieveMachine()
	rand.Seed(time.Now().UnixNano())
	numTest := rand.Intn(5000)

	for i := 0 ; i < numTest; i++ {
		prime := <-primes
		if !big.NewInt(int64(prime)).ProbablyPrime(0) {
			t.Errorf("%d isn't a prime number", prime)
		}
	}
}

