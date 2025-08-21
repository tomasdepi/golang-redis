package command

import (
	"bytes"
	"net"

	"github.com/tidwall/resp"
)

type EchoCommand struct {
	msg string
}

func parseEcho(input []resp.Value) EchoCommand {
	return EchoCommand{
		msg: input[1].String(),
	}
}

func (ec EchoCommand) Execute(conn net.Conn) {
	var buf bytes.Buffer
	wr := resp.NewWriter(&buf)

	wr.WriteString(ec.msg)
	conn.Write([]byte(buf.String()))
}
