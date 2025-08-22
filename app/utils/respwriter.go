package utils

import (
	"bytes"
	"fmt"
	"net"

	"github.com/tidwall/resp"
)

func WriteError(conn net.Conn, er error) {
	var buf bytes.Buffer
	wr := resp.NewWriter(&buf)

	wr.WriteError(er)
	conn.Write([]byte(buf.String()))
}

func WriteInteger(conn net.Conn, i int) {
	var buf bytes.Buffer
	wr := resp.NewWriter(&buf)

	wr.WriteInteger(i)
	conn.Write([]byte(buf.String()))
}

func WriteString(conn net.Conn, s string) {
	var buf bytes.Buffer
	wr := resp.NewWriter(&buf)

	wr.WriteString(s)
	conn.Write([]byte(buf.String()))
}

func WriteSimpleString(conn net.Conn, s string) {
	var buf bytes.Buffer
	wr := resp.NewWriter(&buf)

	wr.WriteSimpleString(s)
	conn.Write([]byte(buf.String()))
}

func WriteNull(conn net.Conn) {
	var buf bytes.Buffer
	wr := resp.NewWriter(&buf)

	wr.WriteNull()
	conn.Write([]byte(buf.String()))
}

func WriteTypeOperationError(conn net.Conn) {
	WriteError(conn, fmt.Errorf("WRONGTYPE Operation against a key holding the wrong kind of value"))
}

func WriteArray(conn net.Conn, values []string) {
	var buf bytes.Buffer
	wr := resp.NewWriter(&buf)

	respValues := make([]resp.Value, len(values))

	for i, v := range values {
		respValues[i] = resp.StringValue(v)
	}

	wr.WriteArray(respValues)
	conn.Write([]byte(buf.String()))
}

func WriteEmptyArray(conn net.Conn) {
	WriteArray(conn, []string{})
}
