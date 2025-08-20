package main

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"os"
	"sync"

	"github.com/tidwall/resp"
)

// Ensures gofmt doesn't remove the "net" and "os" imports in stage 1 (feel free to remove this!)
var _ = net.Listen
var _ = os.Exit

var DB sync.Map

func processCommand(command []resp.Value, conn net.Conn) {

	switch command[0].String() {
	case "ECHO":
		//r := fmt.Sprintf("+%s\r\n", command[1].String())
		var buf bytes.Buffer
		wr := resp.NewWriter(&buf)
		wr.WriteString(command[1].String())
		conn.Write([]byte(buf.String()))
	case "SET":
		key := command[1].String()
		value := command[2].String()
		DB.Store(key, value)

		var buf bytes.Buffer
		wr := resp.NewWriter(&buf)
		wr.WriteString("OK")
		conn.Write([]byte(buf.String()))
	case "GET":
		key := command[1].String()
		var buf bytes.Buffer
		wr := resp.NewWriter(&buf)

		value, exist := DB.Load(key)

		if !exist {
			wr.WriteNull()
		} else {
			wr.WriteString(value.(string))
		}

		conn.Write([]byte(buf.String()))

	default:
		fmt.Println("No conozco el comando")
		conn.Write([]byte("+PONG\r\n"))

	}
}

func handleNewClient(conn net.Conn) {

	for {
		buffer := make([]byte, 1024)

		_, err := conn.Read(buffer)

		if err != nil {
			log.Println("Error reading:", err)
			return
		}

		parser := resp.NewReader(bytes.NewBuffer(buffer))
		command := []resp.Value{}

		for {
			v, _, err := parser.ReadValue()

			if err != nil {
				break
			}

			fmt.Printf("Read %s\n", v.Type())
			if v.Type() == resp.Array {
				for i, v := range v.Array() {
					fmt.Printf("  #%d %s, value: '%s'\n", i, v.Type(), v)
					command = append(command, v)
				}
			}
		}

		processCommand(command, conn)
	}
}

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	ln, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		// defer conn.Close()

		go handleNewClient(conn)
	}

}
