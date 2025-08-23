package command

import (
	"fmt"
	"net"

	"github.com/tidwall/resp"
	"github.com/tomasdepi/golang-redis/app/utils"
)

type PingCommand struct {
	msg string
}

func parsePing(input []resp.Value) (PingCommand, error) {

	if len(input) > 2 {
		return PingCommand{}, fmt.Errorf("(error) ERR wrong number of arguments for 'ping' command")
	}

	if len(input) > 1 {
		return PingCommand{
			msg: input[1].String(),
		}, nil
	}

	return PingCommand{}, nil
}

func (pc PingCommand) Execute(conn net.Conn) {
	var msg string

	if pc.msg != "" {
		msg = pc.msg
	} else {
		msg = "PONG"
	}

	utils.WriteSimpleString(conn, msg)
}
