package handler

import (
	"bytes"
	"fmt"
	"log"
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/command"
	"github.com/tidwall/resp"
)

// var DB sync.Map

// func processCommand(command []resp.Value, conn net.Conn) {

// 	var buf bytes.Buffer
// 	wr := resp.NewWriter(&buf)

// 	switch command[0].String() {
// 	case "ECHO":
// 		wr.WriteString(command[1].String())
// 	case "SET":
// 		key := command[1].String()
// 		value := command[2].String()

// 		rv := RedisValue{
// 			Val:  value,
// 			Type: 1,
// 		}

// 		DB.Store(key, rv)

// 		wr.WriteString("OK")
// 	case "GET":
// 		key := command[1].String()

// 		if rv, exist := DB.Load(key); exist {
// 			rv := rv.(RedisValue)
// 			wr.WriteString(rv.Val.(string))
// 		} else {
// 			wr.WriteNull()
// 		}

// 	default:
// 		wr.WriteError(fmt.Errorf("ERR unknown command '%s'", command[0].String()))
// 	}

// 	conn.Write([]byte(buf.String()))
// }

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
