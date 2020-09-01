package main

import (
	"fmt"
	"sync"
	"time"
)

type ActionWaitGroup struct {
	Name     string
	Sequence int
	Value    float64
}

func main() {
	var channel = make(chan ActionWaitGroup)
	defer close(channel)

	var wg sync.WaitGroup

	fmt.Println("START")

	// Top Up Value
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go emitWithWaitGroup(ActionWaitGroup{
			Name:  "Top-Up",
			Sequence: i,
			Value: 1000,
		}, channel, &wg)
	}

	// Receiving
	go receiverWithWaitGroup(channel, &wg)

	wg.Wait()
	fmt.Println("FINISH")
}

func emitWithWaitGroup(act ActionWaitGroup, channel chan<- ActionWaitGroup, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("[Emitting] %s (%d) - Rp. %v,00.\n", act.Name, act.Sequence, act.Value)
	time.Sleep(2 * time.Second)
	channel <- act
}

func receiverWithWaitGroup(channel <-chan ActionWaitGroup, wg *sync.WaitGroup) {
	for act := range channel {
		wg.Add(1)
		fmt.Printf("[Receiver] %s (%d) - Rp. %v,00.\n", act.Name, act.Sequence, act.Value)
		wg.Done()
	}
}