package main

import "fmt"

func jogador1() {

}

func main() {

	messages := make(chan string)

	go func() { messages <- "mensagem" }()

	msg := <-messages
	fmt.Println(msg)
}
