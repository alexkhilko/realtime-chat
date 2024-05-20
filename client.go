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
	for {
		fmt.Print("send> ")
		input, err := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if err != nil {
			log.Fatalln(err)
		}
		if input == "quit" || input == "exit" {
			break
		}
		_, err = conn.Write([]byte(input + "\n"))
		if err != nil {
			log.Fatal(err)
		}

		msg, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Print("recv> ", msg)
	}
}