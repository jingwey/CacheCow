package main

import (
	"bufio"
	"io"
	"strconv"
)

const (
	STRING = '+'
	ERROR = '-'
	INTEGER = ':'
	BULK = '$'
	ARRAY  = '*'
)

type Value struct {
	typ string
	str string
	num int
	bulk string
	array []Value
}

// RESP - Redis Serialization Protocol Specification
// https://redis.io/docs/latest/develop/reference/protocol-spec/
type Resp struct {
	reader *bufio.Reader
}

func NewResp(rd io.Reader) *Resp {
	return &Resp{reader: bufio.NewReader(rd)}
}

// input := "$5\r\nAhmed\r\n"
func (r *Resp) readLine() (line []byte, n int, err error) {
	for {
		b, err := r.reader.ReadByte()
		if err != nil {
			return nil, 0, err
		}

		n += 1

		line = append(line, b)
		if len(line) >= 2 && line[len(line)-2] == '\r' {
			break
		}
	}
	return line[:len(line)-2], n, nil
}

func (r *Resp) readInteger() (x int, n int, err error) {
	line, n, err := r.readLine()
	if err != nil {
		return 0, 0, err
	}
	i64, err :=  strconv.ParseInt(string(line), 10, 64)
	if err != nil {
		return 0, n, err
	}

	return int(i64), n, nil
}