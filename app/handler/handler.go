package handler

import (
	"bufio"
	"errors"
	"io"
	"log"
	"net"

	"github.com/tidwall/resp"
	"github.com/tomasdepi/golang-redis/app/command"
	"github.com/tomasdepi/golang-redis/app/utils"
)

func HandleNewClient(conn net.Conn) {
	defer conn.Close()

	reader := resp.NewReader(bufio.NewReader(conn))

	for {
		v, _, err := reader.ReadValue()
		if err != nil {
			if errors.Is(err, io.EOF) {
				log.Println("Client disconnected")
				return
			}
			log.Println("Error reading value:", err)
			utils.WriteError(conn, err)
			continue
		}

		var respValues []resp.Value
		if v.Type() == resp.Array {
			respValues = append(respValues, v.Array()...)
		} else {
			respValues = append(respValues, v)
		}

		rc, err := command.ParseCommand(respValues)

		if err != nil {
			utils.WriteError(conn, err)
			log.Println(err)
			continue
		}

		rc.Execute(conn)
	}
}
