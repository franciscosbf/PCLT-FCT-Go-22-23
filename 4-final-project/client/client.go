package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	const numArgs = 4
	if len(os.Args) != numArgs {
		fmt.Printf("Usage: ./%s <host:port> <message> <max>", os.Args[0])
		return
	}
	host := os.Args[1]
	message := os.Args[2]
	max, err := strconv.ParseUint(os.Args[3], 10, 64)
	if err != nil {
		fmt.Printf("%s is not a number.\n", os.Args[3])
		return
	}

	//These are here to make the compiler cooperate. Remove.
	_ = host
	_ = message
	_ = max

	printResult(0, 0)
}

// printResult prints the final result to stdout.
func printResult(hash, nonce uint64) {
	fmt.Println("Result", hash, nonce)
}

// printDisconnected prints a disconnected message to stdout.
func printDisconnected() {
	fmt.Println("Disconnected from the server.")
}
