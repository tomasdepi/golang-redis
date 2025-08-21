package handler

import (
	"bytes"
	"fmt"
	"log"
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/command"
	"github.com/tidwall/resp"
)

func HandleNewClient(conn net.Conn) {

	for {
		buffer := make([]byte, 1024)
		_, err := conn.Read(buffer)

		if err != nil {
			log.Println("Error reading:", err)
			return
		}

		parser := resp.NewReader(bytes.NewBuffer(buffer))
		respValues := []resp.Value{}

		for {
			v, _, err := parser.ReadValue()

			if err != nil {
				break
			}

			fmt.Printf("Read %s\n", v.Type())
			if v.Type() == resp.Array {
				for i, v := range v.Array() {
					fmt.Printf("  #%d %s, value: '%s'\n", i, v.Type(), v)
					respValues = append(respValues, v)
				}
			}
		}

		rc, err := command.ParseCommand(respValues)

		if err != nil {
			log.Println(err)
		}

		rc.Execute(conn)
		//processCommand(respValues, conn)
		//buffer = buffer[:0] // resets buffer
	}
}

// TODO
// 1. better buffer handling
// 2. non case sensitive commands
// 3. check command parse errors
