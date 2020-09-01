package main

import (
	"fmt"
	"time"
)

type Action struct {
	Name     string
	Sequence int
	Value    float64
}

func main()  {
	var channel = make(chan Action)
	defer close(channel)

	// Top Up Value
	for i := 0; i < 5; i++ {
		go emit(Action{
			Name:  "Top-Up",
			Sequence: i,
			Value: 1000,
		}, channel)
	}

	// Receiver
	go receiver(channel)

	time.Sleep(5 * time.Second)
}

func emit(act Action, channel chan<- Action) {
	fmt.Printf("[Emitting] %s (%d) - Rp. %v,00.\n", act.Name, act.Sequence, act.Value)
	channel <- act
}

func receiver(channel <-chan Action) {
	for act := range channel {
		fmt.Printf("[Receiver] %s (%d) - Rp. %v,00.\n", act.Name, act.Sequence, act.Value)
	}
}