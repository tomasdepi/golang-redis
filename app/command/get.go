package command

import (
	"net"
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/utils"
	"github.com/tidwall/resp"
)

type GetCommand struct {
	key string
}

func parseGet(input []resp.Value) GetCommand {
	return GetCommand{
		key: input[1].String(),
	}
}

func (gc GetCommand) Execute(conn net.Conn) {

	rv, ok := DB.Load(gc.key)

	if !ok {
		utils.WriteNull(conn)
		return
	}

	// check expiracy
	if !rv.Expires {
		utils.WriteString(conn, rv.Val.(string))
	} else {
		if rv.ExpiresAt > time.Now().UnixMilli() {
			utils.WriteString(conn, rv.Val.(string))
		} else {
			utils.WriteNull(conn)
			DB.Delete(gc.key)
		}
	}

	// rv.ExpiresAt > time.Now().UnixMilli() {
	// 	wr.WriteNull()
	// } else {
	// 	wr.WriteString(rv.Val.(string))
	// }
}
