package handler

import (
	"bytes"
	"fmt"
	"log"
	"net"

	"github.com/tidwall/resp"
	"github.com/tomasdepi/golang-redis/app/command"
	"github.com/tomasdepi/golang-redis/app/utils"
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
			utils.WriteError(conn, err)
			log.Println(err)
			continue
		}

		rc.Execute(conn)
		//buffer = buffer[:0] // resets buffer
	}
}

// TODO
// 1. better buffer handling
// 2. check command parse errors
