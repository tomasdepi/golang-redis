package command

import (
	"bytes"
	"fmt"
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/db"
	"github.com/tidwall/resp"
)

// const AVAILABLE_COMMANDS = []string{"ECHO", "SET", "GET"}

var DB = db.RedisDB{}

type RedisCommand interface {
	Execute(conn net.Conn)
}

type CCommand struct {
	msg string
}

func (cc CCommand) Execute(conn net.Conn) {
	var buf bytes.Buffer
	wr := resp.NewWriter(&buf)

	wr.WriteString("PONG")
	conn.Write([]byte(buf.String()))
}

func ParseCommand(input []resp.Value) (RedisCommand, error) {

	if len(input) == 0 {
		return nil, fmt.Errorf("command is empty")
	}

	c := input[0].String()

	switch c {
	case "ECHO":
		return parseEcho(input), nil
	case "SET":
		return parseSet(input), nil
	case "GET":
		return parseGet(input), nil
	case "COMMAND":
		return parseCC(input), nil
	case "PING":
		return parsePing(input), nil
	case "RPUSH":
		return ParseRPush(input), nil
	default:
		return nil, fmt.Errorf("redis Command %s not supported", c)
	}
}

func parseCC(input []resp.Value) CCommand {
	return CCommand{
		msg: input[1].String(),
	}
}
