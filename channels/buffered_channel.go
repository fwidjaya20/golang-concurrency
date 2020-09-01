package main

import "fmt"

type BuffBlockingAction struct {
	Name     string
	Sequence int
	Value    float64
}

func main()  {
	var channel = make(chan BuffBlockingAction, 2)

	buffBlockEmit(BuffBlockingAction{
		Name:     "Top-Up",
		Sequence: 1,
		Value:    5000,
	}, channel)
	buffBlockEmit(BuffBlockingAction{
		Name:     "Top-Up",
		Sequence: 2,
		Value:    5000,
	}, channel)

	act1, act2 := <-channel, <-channel

	fmt.Printf("[Receiver] %s (%d) - Rp. %v,00.\n", act1.Name, act1.Sequence, act1.Value)
	fmt.Printf("[Receiver] %s (%d) - Rp. %v,00.\n", act2.Name, act2.Sequence, act2.Value)
}

func buffBlockEmit(act BuffBlockingAction, channel chan<- BuffBlockingAction) {
	channel <- act
}