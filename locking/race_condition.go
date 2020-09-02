package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type Wallet struct {
	Balance float64
}

func main() {
	var channel = make(chan error)
	var wg sync.WaitGroup
	var wallet Wallet = Wallet{
		Balance: 10000,
	}

	wg.Add(1)
	go transaction("T1", "Maman", 10000, &wallet, channel, &wg)

	wg.Add(1)
	go transaction("T2", "Budi", -10000, &wallet, channel, &wg)

	wg.Add(1)
	go transaction("T3", "Tono", 5000, &wallet, channel, &wg)

	go func() {
		wg.Wait()
		close(channel)
	}()

	for err := range channel {
		if err != nil {
			fmt.Println(err)
		}
	}

	fmt.Printf("Saldo Sekarang : Rp.%v,00.\n", wallet.Balance)
}

func transaction(id string, stakeholder string, amount float64, wallet *Wallet, channel chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()

	state := "Penambahan Dana"

	if amount < 0 {
		state = "Pengurangan Dana"
	}

	time.Sleep(2 * time.Second)

	balanceProcessed := wallet.Balance + amount

	if balanceProcessed < 0 {
		channel <- errors.New(fmt.Sprintf("[%s] Saldo tidak cukup.", id))
		return
	}

	wallet.Balance = balanceProcessed
	fmt.Printf("[%s] %s %s sebesar Rp. %v,00. Uang anda : Rp. %v,00.\n", id, stakeholder, state, amount, wallet.Balance)
	channel <- nil
}