# Locking
`Race Condition` merupakan kondisi dimana satu atau lebih **thread** mengakses **shared memory** disaat yang bersamaaan. Ketika hal ini terjadi, maka value yang ada didalam memory tersebut akan menjadi kacau.

Locking (`Mutex`) merupakan salah satu tools yang bisa digunakan untuk mencegah terjadinya `Race Condition`. Mutex melakukan perubahan level akses data terhadap thread, sehingga thread yang tidak memperoleh akses tersebut akan menunggu sampai proses yang sedang berlangsung selesai.

### Contoh Race Condition
Saldo Wallet : Rp. 5000,00

- Thread 1 : akan mengurangi saldo sebesar 3000, dengan proses time (10s)
- Thread 2 : akan menambah saldo sebesar 10000, dengan proses time (5s)

| Thread 1 | Thread 2 | Value | Time |
| ------------- | ------------- | ------------- | ------------- |
| -  | -  | 5000  | -  |
| Read  | -  | 5000  | 2s  |
| -  | Read  | 5000  | 1s  |
| -3000  | -  | 5000  | 4s  |
| -  | +100000  | 5000  | 3s  |
| Write  | -  | 2000  | 4s  |
| -  | Write  | 15000  | 1s  |
Hasil akhir dari Saldo Wallet adalah: Rp. 2000,00 (Expected Value: Rp. 12.000,00) karena proses  Thread 1 paling lama.
Pada kasus ini terjadi `Race Condition` dimana kedua Thread mengakses Shared Memory yang sama.

```go
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
```
Snippet code diatas akan menghasilkan result:
```text
[T3] Tono Penambahan Dana sebesar Rp. 5000,00. Uang anda : Rp. 5000,00.
[T2] Budi Pengurangan Dana sebesar Rp. -10000,00. Uang anda : Rp. 15000,00.
[T1] Maman Penambahan Dana sebesar Rp. 10000,00. Uang anda : Rp. 15000,00.
Saldo Sekarang : Rp.15000,00.
```
dimana hasil tersebut terjadi `Race Condition` karena ke-3 proses tersebut dijalankan di waktu yang bersamaan.
- Untuk Proses **T3**, Saldo awal adalah: `Rp.10.000`, Proses ini nemambahkan dana sebesar `Rp.5.000` sehinga Expected Saldo : `Rp.15.000`.
- Untuk Proses **T2**, Saldo awal adalah: `Rp.10.000`, Proses ini mengurangi dana sebesar `Rp.10.000` sehinga Expected Saldo : `Rp.0`.
- Untuk Proses **T1**, Saldo awal adalah: `Rp.10.000`, Proses ini nemambahkan dana sebesar `Rp.10.000` sehinga Expected Saldo : `Rp.20.000`.

### Mutex
Untuk mencegah hal tersebut, dapat menggunakan sistem `Locking` dimana Shared Memory hanya boleh di akses oleh 1 proses.

```go
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
```
Snippet code diatas menghasilkan result:
```text
[T1] Maman Penambahan Dana sebesar Rp. 10000,00. Uang anda : Rp. 20000,00.
[T3] Tono Penambahan Dana sebesar Rp. 5000,00. Uang anda : Rp. 25000,00.
[T2] Budi Pengurangan Dana sebesar Rp. -10000,00. Uang anda : Rp. 15000,00.
Saldo Sekarang : Rp.15000,00.
```
hasil diatas tidak terjadi `Race Condition`.
- Proses **T1**, membaca Saldo Awal `Rp.10.000`. Proses ini menambahkan saldo sebesar `Rp.10.000`, sehingga Saldo Akhir menjadi `Rp.20.000`.
- Proses **T3**, membaca Saldo Awal `Rp.20.000`. Proses ini menambahkan saldo sebesar `Rp.5.000`, sehingga Saldo Akhir menjadi `Rp.25.000`.
- Proses **T2**, membaca Saldo Awal `Rp.25.000`. Proses ini mengurangi saldo sebesar `Rp.10.000`, sehingga Saldo Akhir menjadi `Rp.15.000`.