package command

import (
	"fmt"
	"net"
	"time"

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

	itemsRemoved := pop(rv, lpop.key, lpop.count)

	if lpop.count == 1 {
		utils.WriteString(conn, itemsRemoved[0])
	} else {
		utils.WriteArray(conn, itemsRemoved)
	}
}

type BLpopCommand struct {
	key     string
	timeout int
}

func parseBLpop(input []resp.Value) (BLpopCommand, error) {

	if len(input) != 3 {
		return BLpopCommand{}, fmt.Errorf("(error) ERR wrong number of arguments for 'blpop' command")
	}

	return BLpopCommand{
		key:     input[1].String(),
		timeout: input[2].Integer(),
	}, nil
}

func (blpop BLpopCommand) Execute(conn net.Conn) {

	rv, ok := DB.Load(blpop.key)

	if !ok {
		// key does not exist, so blocking

		var maxDuration time.Duration = 1<<63 - 1

		if blpop.timeout != 0 {
			maxDuration = time.Duration(blpop.timeout * int(time.Second))
		}

		ch := make(chan string, 1)
		DB.AddWaiter(blpop.key, ch)

		select {
		case val := <-ch:
			utils.WriteArray(conn, []string{blpop.key, val})
		case <-time.After(maxDuration):
			utils.WriteNull(conn)
		}

	} else {
		// key exist so pop value
		itemsRemoved := pop(rv, blpop.key, 1)

		utils.WriteArray(conn, []string{blpop.key, itemsRemoved[0]})
	}
}

func pop(rv db.RedisValue, key string, count int) []string {

	newRedisValue := db.RedisValue{
		Val:  rv.Val.([]string)[count:],
		Type: db.ARRAY_VALUE,
	}

	DB.Store(key, newRedisValue)

	return rv.Val.([]string)[:count]
}
