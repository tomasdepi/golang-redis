package command

import (
	"fmt"
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/db"
	"github.com/codecrafters-io/redis-starter-go/app/utils"
	"github.com/tidwall/resp"
)

type LpopCommand struct {
	key   string
	count int
}

func parseLpop(input []resp.Value) (LpopCommand, error) {

	var c = 1

	if len(input) > 3 {
		return LpopCommand{}, fmt.Errorf("(error) ERR wrong number of arguments for 'lpop' command")
	}

	if len(input) == 3 {
		c = input[2].Integer()
	}

	return LpopCommand{
		key:   input[1].String(),
		count: c,
	}, nil
}

func (lpop LpopCommand) Execute(conn net.Conn) {

	rv, ok := DB.Load(lpop.key)

	if !ok {
		utils.WriteNull(conn)
		return
	}

	if rv.Type != db.ARRAY_VALUE {
		utils.WriteTypeOperationError(conn)
		return
	}

	newRedisValue := db.RedisValue{
		Val:  rv.Val.([]string)[lpop.count:],
		Type: db.ARRAY_VALUE,
	}

	DB.Store(lpop.key, newRedisValue)

	itemsRemoved := rv.Val.([]string)[:lpop.count+1]

	if lpop.count == 1 {
		utils.WriteString(conn, itemsRemoved[1])
	} else {
		utils.WriteArray(conn, itemsRemoved)
	}
}
