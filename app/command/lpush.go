package command

import (
	"fmt"
	"net"
	"slices"

	"github.com/codecrafters-io/redis-starter-go/app/db"
	"github.com/codecrafters-io/redis-starter-go/app/utils"
	"github.com/tidwall/resp"
)

type LPushCommand struct {
	key      string
	elements []string
}

func ParseLPush(input []resp.Value) (LPushCommand, error) {

	if len(input) < 3 {
		return LPushCommand{}, fmt.Errorf("(error) ERR wrong number of arguments for 'lpush' command")
	}

	var parsedElements []string

	for _, e := range input[2:] {
		parsedElements = append(parsedElements, e.String())
	}

	return LPushCommand{
		key:      input[1].String(),
		elements: parsedElements,
	}, nil
}

func (lpc LPushCommand) Execute(conn net.Conn) {

	var newSlice []string

	if rv, ok := DB.Load(lpc.key); !ok {
		// key does not exist
		newSlice = lpc.elements
		slices.Reverse(newSlice)
	} else {
		// key already exists
		sliceToAppend := lpc.elements
		slices.Reverse(sliceToAppend)
		newSlice = append(sliceToAppend, rv.Val.([]string)...)
	}

	newRedisValue := db.RedisValue{
		Val:  newSlice,
		Type: db.ARRAY_VALUE,
	}

	DB.Store(lpc.key, newRedisValue)

	utils.WriteInteger(conn, len(newSlice))
}
