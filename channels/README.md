# Channel
Layer komunikasi antar `Goroutines`.

### Make a Channel
```go
var example = make(chan [data_type])
```

### Write Channel data
```go
example <- "Hello World"
```

### Read Channel Data
```go
var value = <-example
```

### Close Channel
```go
close(example)
```

### Example Using Channel
```go
func main() {
    var example = make(chan string)
    defer close(example)    

    go foo("Hello World 01", example)
    go foo("Hello World 02", example)
    go foo("Hello World 03", example)
    
    fmt.Println(<-example)
    fmt.Println(<-example)
    fmt.Println(<-example)
}

func foo(text string, c chan string) {
    c <- text
}
```

### Example communication between Goroutines
```go
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
```
Pada snippet diatas, menghasilkan result sebagai berikut:
```text
[Emitting] Top-Up (0) - Rp. 1000,00.
[Emitting] Top-Up (4) - Rp. 1000,00.
[Receiver] Top-Up (0) - Rp. 1000,00.
[Emitting] Top-Up (3) - Rp. 1000,00.
[Receiver] Top-Up (4) - Rp. 1000,00.
[Emitting] Top-Up (2) - Rp. 1000,00.
[Receiver] Top-Up (3) - Rp. 1000,00.
[Emitting] Top-Up (1) - Rp. 1000,00.
[Receiver] Top-Up (2) - Rp. 1000,00.
[Receiver] Top-Up (1) - Rp. 1000,00.
```

### Different Channel Direction
| Syntax  | Penjelasan |
| ------------- | ------------- |
| `channel chan string` 	| Dapat digunakan untuk **mengirim** dan **menerima** data. 	|
| `channel chan<- string` 	| Hanya dapat digunakan untuk **mengirim** data. 	|
| `channel <-chan string` 	| Hanya dapat digunakan untuk **menerima** data. 	|

### Channel is Blocking
```go
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
```
Result:
```text
fatal error: all goroutines are asleep - deadlock!
```
Pada dasarnya `Channel` memiliki buffer dengan size `0`, sehingga akan di `Block`.

How to resolve ? Kita bisa gunakan **Buffered Channel** dan **Wait Group**.

### Buffered Channel
```text
make(chan [tipe data], [size])
```
Example
```go
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
```

### Wait Group
`WaitGroup` digunakan untuk mengsingkronisasi `Goroutine`. Mekanisme WaitGroup sama dengan `await`.

#### 3 Main WaitGroup Function
| Function  | Explaination |
| ------------- | ------------- |
| `Add(delta)` 	| Increase the counter process by (n) given.	|
| `Done()` 	| Decrease the counter process by 1. 	|
| `Wait()` 	| Block the execution until the counter becomes 0. 	|

Example 01:
```go
func main() {
    var wg sync.WaitGroup

    fmt.Println("START")
    
    wg.Add(1) // Counter 1
    wg.Add(1) // Counter 2
    wg.Add(2) // Counter 4
    wg.Done() // Counter 3
    time.Sleep(time.Second) // Sleep a second
    wg.Done() // Counter 2
    wg.Done() // Counter 1
    time.Sleep(time.Second) // Sleep a second
    wg.Done() // Counter 0

    wg.Wait() // Block below execution until the counter is 0

    fmt.Println("FINISH")
}
```
statement `fmt.Println("FINISH")` akan di jalankan setelah `WaitGroup` counter telah bernilai 0.

Example 02:
```go
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
```
Snippet code diatas akan menghasilkan result sebagai berikut:
```text
START Example
Worker Seq : 5 - START
Worker Seq : 4 - START
Worker Seq : 3 - START
Worker Seq : 1 - START
Worker Seq : 2 - START
Worker Seq : 1 - DONE
Worker Seq : 2 - DONE
Worker Seq : 4 - DONE
Worker Seq : 5 - DONE
Worker Seq : 3 - DONE
FINISH Example <- Wait Until Worker Done
```

Example 03:
```go
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
```
Snippet code diatas akan menghasilkan result sebagai berikut:
```text
START
[Emitting] Top-Up (1) - Rp. 1000,00.
[Emitting] Top-Up (3) - Rp. 1000,00.
[Emitting] Top-Up (0) - Rp. 1000,00.
[Emitting] Top-Up (2) - Rp. 1000,00.
[Emitting] Top-Up (4) - Rp. 1000,00.
[Receiver] Top-Up (4) - Rp. 1000,00.
[Receiver] Top-Up (1) - Rp. 1000,00.
[Receiver] Top-Up (0) - Rp. 1000,00.
[Receiver] Top-Up (2) - Rp. 1000,00.
[Receiver] Top-Up (3) - Rp. 1000,00.
FINISH <- Wait until TopUp and Receiver done
```