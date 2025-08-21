package command

import (
	"fmt"
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/utils"
	"github.com/tidwall/resp"
)

type EchoCommand struct {
	msg string
}

func parseEcho(input []resp.Value) (EchoCommand, error) {

	if len(input) != 2 {
		return EchoCommand{}, fmt.Errorf("(error) ERR wrong number of arguments for 'echo' command")
	}

	return EchoCommand{
		msg: input[1].String(),
	}, nil
}

func (ec EchoCommand) Execute(conn net.Conn) {
	utils.WriteString(conn, ec.msg)
}
