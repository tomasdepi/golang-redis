package command

import (
	"bytes"
	"net"

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
	var buf bytes.Buffer
	wr := resp.NewWriter(&buf)

	if pc.msg != "" {
		wr.WriteString(pc.msg)
	} else {
		wr.WriteString("PONG")
	}

	conn.Write([]byte(buf.String()))
}
