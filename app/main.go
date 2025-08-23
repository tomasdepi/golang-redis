package main

import (
	"fmt"
	"net"
	"os"

	"github.com/tomasdepi/golang-redis/app/handler"
)

// Ensures gofmt doesn't remove the "net" and "os" imports in stage 1 (feel free to remove this!)
var _ = net.Listen
var _ = os.Exit

func main() {

	ln, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	defer ln.Close()
	fmt.Println("Ready to accept conntections!")

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go handler.HandleNewClient(conn)
	}

}
