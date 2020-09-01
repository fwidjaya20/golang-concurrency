package main

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup

	fmt.Println("START Example")

	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go worker(i, &wg)
	}

	wg.Wait()
	fmt.Println("FINISH Example")
}

func worker(sequence int, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("Worker Seq : " + strconv.Itoa(sequence) + " - START")
	time.Sleep(2 * time.Second)
	fmt.Println("Worker Seq : " + strconv.Itoa(sequence) + " - DONE")
}