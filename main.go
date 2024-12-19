package main

import (
	"fmt"
	"io"
	"net"
	"os"
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
		buf := make([]byte, 1024)

		// read message from client
		_, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}

			fmt.Println("error reading from client: ", err.Error())
			os.Exit(1)
		}


		conn.Write([]byte("+OK\r\n"))

	}

}