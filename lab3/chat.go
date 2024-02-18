// Demonstration of channels with a chat application
// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// Chat is a server that lets clients chat with each other.

package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

type client struct {
	name    string
	channel chan<- string
}

var (
	entering = make(chan client)
	leaving  = make(chan client)
	messages = make(chan string) // all incoming client messages
)

func main() {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}

	go broadcaster()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}

func broadcaster() {
	clients := make(map[client]bool) // all connected clients
	for {
		select {
		case msg := <-messages:
			// Broadcast incoming message to all
			// clients' outgoing message channels.
			for cli := range clients {
				cli.channel <- msg
			}

		case cli := <-entering:
			clients[cli] = true
			go func() {
				cli.channel <- "Current Connected Clients:"
				for i := range clients {
					cli.channel <- i.name
				}
			}()

		case cli := <-leaving:
			delete(clients, cli)
			close(cli.channel)
		}
	}
}

func handleConn(conn net.Conn) {
	ch := make(chan string) // outgoing client messages
	go clientWriter(conn, ch)
	fmt.Fprint(conn, "Welcome to the chat room please enter your name:")
	scanner := bufio.NewScanner(conn)
	scanner.Scan()
	NAME := scanner.Text()
	CLIENT := client{name: NAME, channel: ch}
	ch <- "You are " + CLIENT.name
	messages <- NAME + " has arrived"
	entering <- CLIENT

	input := bufio.NewScanner(conn)
	for input.Scan() {
		messages <- NAME + ": " + input.Text()
	}
	// NOTE: ignoring potential errors from input.Err()

	leaving <- CLIENT
	messages <- NAME + " has left"
	conn.Close()
}

func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg) // NOTE: ignoring network errors
	}
}
