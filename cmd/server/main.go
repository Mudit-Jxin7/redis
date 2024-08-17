package main

import (
	"fmt"
	"net"
	"os"
	"redis/internal/server"
)

func main() {
	listener, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println("Error starting server:", err)
		os.Exit(1)
	}

	fmt.Println("Server is running on port 6379...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go server.HandleConnection(conn)
	}
}
