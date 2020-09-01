# Goroutines
Goroutine merupakan `Mini Thread` dimana proses eksekusinya dilakukan secara `Asynchronous` sehingga tiap proses tidak saling menunggu satu sama yang lain.

### Proses Blocking `Synchronous`
`Synchronous` menjalankan prosesnya secara sequential (bertahap).
```go
func main() {
    for i := 1; i <= 5; i++ {
        foo(i)
    }
    
    bar()
    time.Sleep(5 * time.Second)
}

func foo(num int) {
    fmt.Println("Process Foo - " + strconv.Itoa(num))
    time.Sleep(1 * time.Second)
}

func bar() {
    fmt.Println("Process Bar")
}
```
Pada snippet diatas, menghasilkan result sebagai berikut:
```text
Process Foo - 1
Process Foo - 2
Process Foo - 3
Process Foo - 4
Process Foo - 5
Process Bar       <- Blocking Occur
```

### Proses Non Blocking `Asynchronous`
```go
func main() {
    for i := 1; i <= 5; i++ {
        go foo(i)
    }
    
    bar()
    time.Sleep(5 * time.Second)
}

func foo(num int) {
    fmt.Println("Process Foo - " + strconv.Itoa(num))
    time.Sleep(1 * time.Second)
}

func bar() {
    fmt.Println("Process Bar")
}
```
Pada snippet diatas, menghasilkan result sebagai berikut:
```text
Process Bar
Process Foo - 5
Process Foo - 3
Process Foo - 4
Process Foo - 1
Process Foo - 2
```
dimana kedua fungsi `Foo` dan `Bar` dijalankan bersamaan.