package command

import (
	"fmt"
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/db"
	"github.com/codecrafters-io/redis-starter-go/app/utils"
	"github.com/tidwall/resp"
)

type LRangeCommand struct {
	key   string
	start int
	stop  int
}

func parseLRange(input []resp.Value) (LRangeCommand, error) {

	if len(input) != 4 {
		return LRangeCommand{}, fmt.Errorf("(error) ERR wrong number of arguments for 'lrange' command")
	}

	return LRangeCommand{
		key:   input[1].String(),
		start: input[2].Integer(),
		stop:  input[3].Integer(),
	}, nil

}

func (lrc LRangeCommand) Execute(conn net.Conn) {
	// If the list does not exist, an empty array is returned
	// If the start index is greater than or equal to the list's length, an empty array is returned.
	// If the stop index is greater than or equal to the list's length, the stop index is treated as the last element.
	// If the start index is greater than the stop index, the result is an empty array.

	if rv, ok := DB.Load(lrc.key); !ok {
		utils.WriteArray(conn, []string{})
		return
	} else {

		if rv.Type != db.ARRAY_VALUE {
			utils.WriteTypeOperationError(conn)
			return
		}

		slice := rv.Val.([]string)

		if lrc.start >= len(slice) || lrc.start > lrc.stop {
			utils.WriteArray(conn, []string{})
			return
		}

		// An index of -1 refers to the last element, -2 to the second last, and so on. If a negative index is out of range (i.e. >= the length of the list), it is treated as 0 (start of the list)
		var start int
		var stop int

		if lrc.start < 0 {
			start = min(0, len(slice)+start)
		}

		if lrc.stop < 0 {
			stop = max(0, len(slice)+stop)
		}

		stop = min(stop+1, len(slice)) // because LRANGE includes stop_index but golang does not

		partialSlice := slice[start:stop]

		utils.WriteArray(conn, partialSlice)
	}

}
