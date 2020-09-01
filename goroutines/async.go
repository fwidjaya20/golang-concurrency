package main

import (
	"fmt"
	"strconv"
	"time"
)

func main() {
	for i := 1; i <= 5; i++ {
		go asyncProc1(i)
	}

	asyncProc2()
	time.Sleep(5 * time.Second)
}

func asyncProc1(num int) {
	fmt.Println("Process Foo - " + strconv.Itoa(num))
	time.Sleep(1 * time.Second)
}

func asyncProc2() {
	fmt.Println("Process Bar")
}