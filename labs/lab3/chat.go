package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

// client struct represents a chat client with a channel for sending messages
// and a name identifier.
type client struct {
	channel chan<- string // Channel for outgoing messages
	name    string        // Name of the client
}

// Global channels for managing entering and leaving clients, and broadcasting messages.
var (
	entering = make(chan client)    // Channel for clients trying to enter the chat
	leaving  = make(chan client)    // Channel for clients trying to leave the chat
	messages = make(chan string)    // Channel for broadcasting messages to all clients
)

func main() {
	// Start listening on TCP port 8000 on localhost.
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err) // Log and exit on error
	}

	go broadcaster() // Start the broadcaster in a new goroutine
	for {
		conn, err := listener.Accept() // Accept new connections
		if err != nil {
			log.Print(err) // Log errors without stopping the server
			continue
		}
		go handleConn(conn) // Handle new connection in a separate goroutine
	}
}

// broadcaster runs in its own goroutine and manages chat state, including
// broadcasting messages and tracking entering and leaving clients.
func broadcaster() {
	clients := make(map[client]bool) // Map to keep track of connected clients
	for {
		select {
		case msg := <-messages:
			// Broadcast incoming message to all clients' channels.
			for cli := range clients {
				cli.channel <- msg
			}

		case cli := <-entering:
			clients[cli] = true // Mark client as connected
			// Generate a message with the list of all connected clients
			var list string
			for c := range clients {
				list += c.name + ", "
			}
			cli.channel <- "Current clients: " + list // Send the list to the new client

		case cli := <-leaving:
			delete(clients, cli) // Remove client from the map
			close(cli.channel)   // Close the client's channel
		}
	}
}

// handleConn handles each client connection.
func handleConn(conn net.Conn) {
	ch := make(chan string) // Create a channel for outgoing messages
	go clientWriter(conn, ch)

	who := conn.RemoteAddr().String()
	cli := client{channel: ch, name: who}

	ch <- "You are " + who
	messages <- who + " has arrived"
	entering <- cli

	input := bufio.NewScanner(conn)
	for input.Scan() {
		messages <- who + ": " + input.Text()
	}

	if err := input.Err(); err != nil {
		log.Println("reading standard input:", err)
	}

	leaving <- cli
	messages <- who + " has left"
	if err := conn.Close(); err != nil {
		log.Println("closing connection:", err)
	}
}

// clientWriter sends messages from the channel to the client's connection.
func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		if _, err := fmt.Fprintln(conn, msg); err != nil {
			log.Println("sending message to client:", err)
			return
		}
	}
}
