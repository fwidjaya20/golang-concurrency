package main

import "fmt"

func main() {
	var example = make(chan string)

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