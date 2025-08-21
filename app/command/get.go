package command

import (
	"bytes"
	"net"
	"time"

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
