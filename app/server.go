package main

import (
	"fmt"
	"net"
	"os"
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

	var data []byte
	_, err = c.Read(data)
	if err != nil {
		fmt.Println("Error reading data: ", err.Error())
	}
	_, err = c.Write([]byte("HTTP/1.1 200 OK /r/n/r/n"))
	if err != nil {
		fmt.Println("Error writing data: ", err.Error())
	}

	c.Close()
}
