package main

import (
    "log"
    "net"
	"fmt"
	"bufio"
	"os"
	"strings"
)

func main() {
	addr := "localhost:7007"
	conn, err := net.Dial("tcp", addr)
	if err !=  nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("name: ")
	name, err := reader.ReadString('\n')
	name = strings.TrimSpace(name)
	if err != nil {
		log.Fatalln(err)
	}
	done := make(chan bool)

	go func(name string) {
		reader := bufio.NewReader(os.Stdin)
		for {
			input, err := reader.ReadString('\n')
			input = strings.TrimSpace(input)
			if err != nil {
				log.Fatalln(err)
				done <- true
				return
			}
			if input == "quit" || input == "exit" {
				done <- true
				return
			}
			msg := fmt.Sprintf("%s: %s \n", name, input)
			_, err = conn.Write([]byte(msg))
			if err != nil {
				log.Fatal(err)
				done <- true
				return
			}
		}
	}(name)

	go func() {
		serverReader := bufio.NewReader(conn)
		for {
			message, err := serverReader.ReadString('\n')
			if err != nil {
				fmt.Println("Error reading from server:", err)
				done <- true
				return
			}
			fmt.Print(message)
		}
	}()

	<-done
	fmt.Println("connection closed")
}