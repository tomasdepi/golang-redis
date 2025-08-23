package command

import (
	"fmt"
	"net"
	"time"

	"github.com/tidwall/resp"
	"github.com/tomasdepi/golang-redis/app/db"
	"github.com/tomasdepi/golang-redis/app/utils"
)

type GetCommand struct {
	key string
}

func parseGet(input []resp.Value) (GetCommand, error) {

	if len(input) != 2 {
		return GetCommand{}, fmt.Errorf("(error) ERR wrong number of arguments for 'get' command")
	}

	return GetCommand{
		key: input[1].String(),
	}, nil
}

func (gc GetCommand) Execute(conn net.Conn) {

	rv, ok := DB.Load(gc.key)

	if !ok {
		utils.WriteNull(conn)
		return
	}

	if rv.Type != db.SINGLE_VALUE {
		utils.WriteTypeOperationError(conn)
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
