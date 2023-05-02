package main

import (
	"flag"
	"log"
	"pclt/miniproject/task1"
	"pclt/miniproject/task2"
	"pclt/miniproject/task3"
	"sync"
)

func runTask1(wg *sync.WaitGroup) {
	task1.Setup(wg)
}

func runTask2(wg *sync.WaitGroup) {
	task2.Setup(wg)
}

func runTask3(wg *sync.WaitGroup) {
	task3.Setup(wg)
}

func main() {
	task := flag.Int("task", 1, "choose between 1 and 3 (task number)")
	flag.Parse()

	if *task < 1 || *task > 3 {
		log.Fatalf("Invalid task number %d. Only accepts 1, 2 or 3", *task)
	}

	log.Printf("Running task %d", *task)
	var wg sync.WaitGroup
	wg.Add(1)
	switch *task {
	case 1:
		runTask1(&wg)
	case 2:
		runTask2(&wg)
	default:
		runTask3(&wg)
	}
	wg.Wait()
}
