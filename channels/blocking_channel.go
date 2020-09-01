package main

import "fmt"

type BlockingAction struct {
	Name     string
	Sequence int
	Value    float64
}

func main()  {
	var channel = make(chan BlockingAction)

	blockEmit(BlockingAction{
		Name:     "Top-Up",
		Sequence: 1,
		Value:    5000,
	}, channel)
	blockEmit(BlockingAction{
		Name:     "Top-Up",
		Sequence: 2,
		Value:    5000,
	}, channel)

	act1, act2 := <-channel, <-channel

	fmt.Printf("[Receiver] %s (%d) - Rp. %v,00.\n", act1.Name, act1.Sequence, act1.Value)
	fmt.Printf("[Receiver] %s (%d) - Rp. %v,00.\n", act2.Name, act2.Sequence, act2.Value)
}

func blockEmit(act BlockingAction, channel chan<- BlockingAction) {
	channel <- act
}