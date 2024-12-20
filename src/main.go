package main

import (
	"fmt"
	"net"
)

func main() {
	l, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("starting CacheCow at :6379...")

	conn, err := l.Accept()
	if err != nil {
		fmt.Println(err)
		return
	}

	defer conn.Close() // close connection once finished

	for {

		resp := NewResp(conn)
		value, err := resp.Read()
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(value)

		conn.Write([]byte("+OK\r\n"))

	}

}
