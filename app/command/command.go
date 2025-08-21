package command

import (
	"bytes"
	"fmt"
	"net"
	"strings"
	"time"

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

type EchoCommand struct {
	msg string
}

type SetCommand struct {
	key string
	val string

	// ex int
	px int64

	// nx bool
	// xx bool

	// get bool
}

type GetCommand struct {
	key string
}

func (cc CCommand) Execute(conn net.Conn) {
	var buf bytes.Buffer
	wr := resp.NewWriter(&buf)

	wr.WriteString("PONG")
	conn.Write([]byte(buf.String()))
}

func (ec EchoCommand) Execute(conn net.Conn) {
	var buf bytes.Buffer
	wr := resp.NewWriter(&buf)

	wr.WriteString(ec.msg)
	conn.Write([]byte(buf.String()))
}

func (sc SetCommand) Execute(conn net.Conn) {

	rv := db.RedisValue{
		Val:       sc.val,
		Type:      1,
		ExpiresAt: time.Now().UnixMilli() + sc.px,
		Expires:   sc.px != 0,
	}

	DB.Store(sc.key, rv)

	var buf bytes.Buffer
	wr := resp.NewWriter(&buf)

	wr.WriteString("OK")
	conn.Write([]byte(buf.String()))
}

func (gc GetCommand) Execute(conn net.Conn) {

	var buf bytes.Buffer
	wr := resp.NewWriter(&buf)

	rv, ok := DB.Load(gc.key)

	if !ok {
		wr.WriteNull()
		conn.Write([]byte(buf.String()))
		return
	}

	// check expiracy
	if !rv.Expires {
		wr.WriteString(rv.Val.(string))
	} else {
		if rv.ExpiresAt > time.Now().UnixMilli() {
			wr.WriteString(rv.Val.(string))
		} else {
			wr.WriteNull()
			DB.Delete(gc.key)
		}
	}

	// rv.ExpiresAt > time.Now().UnixMilli() {
	// 	wr.WriteNull()
	// } else {
	// 	wr.WriteString(rv.Val.(string))
	// }

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
	default:
		return nil, fmt.Errorf("redis Command %s not supported", c)
	}
}

func parseCC(input []resp.Value) CCommand {
	return CCommand{
		msg: input[1].String(),
	}
}

func parseEcho(input []resp.Value) EchoCommand {
	return EchoCommand{
		msg: input[1].String(),
	}
}

func parseSet(input []resp.Value) SetCommand {
	sc := SetCommand{
		key: input[1].String(),
		val: input[2].String(),
	}

	if len(input) > 3 {
		opValue := strings.ToUpper(input[3].String())
		if opValue == "PX" {
			sc.px = int64(input[4].Integer())
		}
	}

	return sc
}

func parseGet(input []resp.Value) GetCommand {
	return GetCommand{
		key: input[1].String(),
	}
}
