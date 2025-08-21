package command

import (
	"bytes"
	"net"
	"strings"
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/db"
	"github.com/tidwall/resp"
)

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
