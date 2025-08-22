package command

import (
	"fmt"
	"net"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/db"
	"github.com/codecrafters-io/redis-starter-go/app/utils"
	"github.com/tidwall/resp"
)

const (
	GET     = "GET"
	ECHO    = "ECHO"
	SET     = "SET"
	COMMAND = "COMMAND"
	PING    = "PING"
	RPUSH   = "RPUSH"
)

var DB = db.RedisDB{}

type RedisCommand interface {
	Execute(conn net.Conn)
}

type CCommand struct {
	msg string
}

func (cc CCommand) Execute(conn net.Conn) {
	utils.WriteString(conn, "PONG")
}

func ParseCommand(input []resp.Value) (RedisCommand, error) {

	if len(input) == 0 {
		return nil, fmt.Errorf("command is empty")
	}

	c := strings.ToUpper(input[0].String())

	switch c {
	case ECHO:
		return parseEcho(input)
	case SET:
		return parseSet(input)
	case GET:
		return parseGet(input)
	case COMMAND:
		return parseCC(input)
	case PING:
		return parsePing(input)
	case RPUSH:
		return ParseRPush(input)
	default:
		return nil, fmt.Errorf("redis Command %s not supported", c)
	}
}

func parseCC(input []resp.Value) (CCommand, error) {
	return CCommand{
		msg: input[1].String(),
	}, nil
}
