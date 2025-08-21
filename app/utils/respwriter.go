package utils

import (
	"bytes"
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

func WriteNull(conn net.Conn) {
	var buf bytes.Buffer
	wr := resp.NewWriter(&buf)

	wr.WriteNull()
	conn.Write([]byte(buf.String()))
}
