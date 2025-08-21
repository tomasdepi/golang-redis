package command

import (
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/utils"
	"github.com/tidwall/resp"
)

type PingCommand struct {
	msg string
}

func parsePing(input []resp.Value) PingCommand {

	if len(input) > 1 {
		return PingCommand{
			msg: input[1].String(),
		}
	}

	return PingCommand{}
}

func (pc PingCommand) Execute(conn net.Conn) {
	var msg string

	if pc.msg != "" {
		msg = pc.msg
	} else {
		msg = "PONG"
	}

	utils.WriteString(conn, msg)
}
