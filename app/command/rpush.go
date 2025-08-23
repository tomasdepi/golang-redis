package command

import (
	"fmt"
	"log"
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/db"
	"github.com/codecrafters-io/redis-starter-go/app/utils"
	"github.com/tidwall/resp"
)

type RPushCommand struct {
	key      string
	elements []string
}

func ParseRPush(input []resp.Value) (RPushCommand, error) {

	if len(input) < 3 {
		return RPushCommand{}, fmt.Errorf("(error) ERR wrong number of arguments for 'rpush' command")
	}

	var parsedElements []string

	for _, e := range input[2:] {
		parsedElements = append(parsedElements, e.String())
	}

	return RPushCommand{
		key:      input[1].String(),
		elements: parsedElements,
	}, nil
}

func (rpc RPushCommand) Execute(conn net.Conn) {

	var newSlice []string

	if rv, ok := DB.Load(rpc.key); !ok {
		// key does not exist
		newSlice = rpc.elements
	} else {
		// key already exists
		newSlice = append(rv.Val.([]string), rpc.elements...)
	}

	newRedisValue := db.RedisValue{
		Val:  newSlice,
		Type: db.ARRAY_VALUE,
	}

	DB.Store(rpc.key, newRedisValue)

	ch, err := DB.PopWaiter(rpc.key)

	if err != nil {
		log.Printf("No waiter found for key %s\n", rpc.key)
	}

	ch <- newSlice[0]

	utils.WriteInteger(conn, len(newSlice))
}
