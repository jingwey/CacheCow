package main

import (
	"bufio"
	"fmt"
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

func (r *Resp) Read() (Value, error) {
	valType, err := r.reader.ReadByte()
	if err != nil {
		return Value{}, err
	}

	switch valType {
	case ARRAY:
		return r.readArray()
	case BULK:
		return r.readBulk()
	default:
		fmt.Printf("Unknown type: %v", string(valType))
		return Value{}, nil
	}

}

func (r *Resp) readArray() (Value, error) {
	v := Value{}
	v.typ = "array"

	// read length of array
	length, _ , err := r.readInteger()
	if err != nil {
		return v, err
	}

	// for each line, parse and read the value
	v.array = make([]Value, length)
	for i := 0; i < length; i++ {
		val, err := r.Read()
		if err != nil {
			return v, err
		}

		v.array[i] = val
	}

	return v, nil
}

func (r *Resp) readBulk() (Value, error) {
	v := Value{}

	v.typ = "bulk"

	len, _, err := r.readInteger()
	if err != nil {
		return v, err
	}

	bulk := make([]byte, len)

	r.reader.Read(bulk)

	v.bulk = string(bulk)
	
	// read the trailing CRLF
	r.readLine()

	return v, nil
}