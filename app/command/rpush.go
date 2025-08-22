package command

import (
	"fmt"
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/db"
	"github.com/codecrafters-io/redis-starter-go/app/utils"
	"github.com/tidwall/resp"
)

type RPushCommand struct {
	key     string
	element string
}

func ParseRPush(input []resp.Value) (RPushCommand, error) {

	if len(input) < 3 {
		return RPushCommand{}, fmt.Errorf("(error) ERR wrong number of arguments for 'rpush' command")
	}

	return RPushCommand{
		key:     input[1].String(),
		element: input[2].String(),
	}, nil
}

func (rpc RPushCommand) Execute(conn net.Conn) {

	var listLen int

	if rv, ok := DB.Load(rpc.key); !ok {
		newRedisValue := db.RedisValue{
			Val:  []string{rpc.element},
			Type: 2,
		}

		DB.Store(rpc.key, newRedisValue)
		listLen = 1
	} else {
		newSlice := append(rv.Val.([]string), rpc.element)
		newRedisValue := db.RedisValue{
			Val:  newSlice,
			Type: 2,
		}

		DB.Store(rpc.key, newRedisValue)
		listLen = len(newSlice)
	}

	utils.WriteInteger(conn, listLen)
}
