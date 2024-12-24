package main

import (
	"fmt"
	"net"
	"strings"
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

		if value.typ != "array" {
			fmt.Println("Invalid request, expected array")
			continue
		}

		if len(value.array) == 0 {
			fmt.Println("Invalid request, expected array length > 0")
			continue
		}

		// extract the command first
		command := strings.ToUpper(value.array[0].bulk)
		// then take the rest of the values
		args := value.array[1:]

		fmt.Println(value)

		writer := NewWriter(conn)

		// pick command from map
		handler, ok := Handlers[command]
		if !ok {
			fmt.Println("Invalid command: ", command)
			writer.Write(Value{typ: "string", str: "INVALID COMMAND"})
			continue
		}

		result := handler(args)
		writer.Write(result)

	}

}
