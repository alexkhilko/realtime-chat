package main

import (
	"fmt"
	"io"
	"log"
	"net"
    "bufio"
    "sync"
)

type Client struct {
	Name string
	Conn net.Conn
}

type Message struct {
    User string
    Msg string
}

var (
    clients = make(map[string]Client)
    broadcast = make(chan Message)
    unregister = make(chan Client)
    register = make(chan Client)
    clientMutext = &sync.Mutex{}
)

func main() {
    addr := "localhost:7007"
    listener, err := net.Listen("tcp", addr)
    if err != nil {
        log.Fatalln(err)
    }
    defer listener.Close()
    log.Println("Server is running on:", addr)
    
    go handleMessages()
    
    for {
        conn, err := listener.Accept()
        fmt.Println("New connection", conn)
        if err != nil {
            log.Println("Failed to accept conn.", err)
            continue
        }
        go handleConnection(conn)
    }
}

func handleMessages() {
    for {
        select {
        case message := <-broadcast:
            clientMutext.Lock()
            for name, client := range clients {
                if name != message.User {
                    client.Conn.Write([]byte(message.Msg))
                }
            }
            clientMutext.Unlock()
        case client := <-register:
            fmt.Println("New client connected:", client.Name)
            clientMutext.Lock()
            clients[client.Name] = client
            clientMutext.Unlock()
        case client := <-unregister:
            fmt.Println("Client disconneted:", client.Name)
            clientMutext.Lock()
            delete(clients, client.Name)
            clientMutext.Unlock()
        }
    }
}

func handleConnection(conn net.Conn) {
    defer conn.Close()
    client := Client{
        Name: conn.RemoteAddr().String(),
        Conn: conn,
    }
    register <- client
    reader := bufio.NewReader(conn)
    for {
        msg, err := reader.ReadString('\n')
        if err != nil {
            if err != io.EOF {
                log.Fatalln("Failed to read data.", err)
            }
            break
        } else {
            broadcast <- Message{User: client.Name, Msg: msg}
        }
    }
    unregister <- client
}
