package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type WalletWithMutex struct {
	mutex sync.Mutex
	Balance float64
}

func main() {
	var channel = make(chan error)
	var wg sync.WaitGroup
	var wallet WalletWithMutex = WalletWithMutex{
		Balance: 10000,
	}

	wg.Add(1)
	go transactionSync("T1", "Maman", 10000, &wallet, channel, &wg)

	wg.Add(1)
	go transactionSync("T2", "Budi", -10000, &wallet, channel, &wg)

	wg.Add(1)
	go transactionSync("T3", "Tono", 5000, &wallet, channel, &wg)

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

func transactionSync(id string, stakeholder string, amount float64, wallet *WalletWithMutex, channel chan<- error, wg *sync.WaitGroup) {
	wallet.mutex.Lock()
	defer wallet.mutex.Unlock()
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