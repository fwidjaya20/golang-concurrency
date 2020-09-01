package main

import (
	"fmt"
	"strconv"
	"time"
)

func main() {
	for i := 1; i <= 5; i++ {
		syncProc1(i)
	}

	syncProc2()
	time.Sleep(5 * time.Second)
}

func syncProc1(num int) {
	fmt.Println("Process Foo - " + strconv.Itoa(num))
	time.Sleep(1 * time.Second)
}

func syncProc2() {
	fmt.Println("Process Bar")
}