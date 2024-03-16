package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	var c net.Conn

	c, err = l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}

	data := make([]byte, 2048)
	_, err = c.Read(data)
	if err != nil {
		fmt.Println("Error reading data: ", err.Error())
	}
	str := string(data)
	fields := strings.Fields(str)
	path := fields[1]

	if path != "/" {
		_, err = c.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
		if err != nil {
			fmt.Println("Error writing data: ", err.Error())
		}
	}

	_, err = c.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
	if err != nil {
		fmt.Println("Error writing data: ", err.Error())
	}

	defer c.Close()

}
