package command

import (
	"net"
	"strings"
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/db"
	"github.com/codecrafters-io/redis-starter-go/app/utils"
	"github.com/tidwall/resp"
)

type SetCommand struct {
	key string
	val string

	// ex int
	px int64

	// nx bool
	// xx bool

	// get bool
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

func (sc SetCommand) Execute(conn net.Conn) {

	rv := db.RedisValue{
		Val:       sc.val,
		Type:      1,
		ExpiresAt: time.Now().UnixMilli() + sc.px,
		Expires:   sc.px != 0,
	}

	DB.Store(sc.key, rv)

	utils.WriteString(conn, "OK")
}
