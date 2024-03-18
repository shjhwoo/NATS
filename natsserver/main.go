package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/nats-io/nats-server/v2/server"
)

// https://dev.to/karanpratapsingh/embedding-nats-in-go-19o
func main() {
	//사용자 종료 감지 로직:
	opts := &server.Options{
		Port: 7070,
	}

	// Initialize new server with options
	ns, err := server.NewServer(opts)

	if err != nil {
		panic(err)
	}

	// Start the server via goroutine
	go ns.Start()

	// Wait for server to be ready for connections
	if !ns.ReadyForConnections(4 * time.Second) {
		panic("not ready for connection")
	}

	connectedClients, _ := ns.Connz(nil)
	fmt.Println("Connected clients: ", connectedClients.Conns)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	<-sigChan
}
