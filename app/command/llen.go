package command

import (
	"fmt"
	"net"

	"github.com/tidwall/resp"
	"github.com/tomasdepi/golang-redis/app/db"
	"github.com/tomasdepi/golang-redis/app/utils"
)

type LlenCommand struct {
	key string
}

func parseLlen(input []resp.Value) (LlenCommand, error) {

	if len(input) != 2 {
		return LlenCommand{}, fmt.Errorf("(error) ERR wrong number of arguments for 'llen' command")
	}

	return LlenCommand{
		key: input[1].String(),
	}, nil
}

func (gc LlenCommand) Execute(conn net.Conn) {

	rv, ok := DB.Load(gc.key)

	if !ok {
		utils.WriteInteger(conn, 0)
		return
	}

	if rv.Type != db.ARRAY_VALUE {
		utils.WriteTypeOperationError(conn)
		return
	}

	utils.WriteInteger(conn, len(rv.Val.([]string)))
}
