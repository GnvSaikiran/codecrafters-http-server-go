package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:4000")
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

	randomString := strings.Split(strings.Trim(path, "/"), "/")[1]

	// building a response
	r := "HTTP/1.1 200 OK\r\n"
	r += "Content-Type: text/plain\r\n"
	r += "Content-Length: " + fmt.Sprint(len(randomString)) + "\r\n\r\n"
	r += randomString

	_, err = c.Write([]byte(r))
	if err != nil {
		fmt.Println("Error writing data: ", err.Error())
	}
	defer c.Close()

}
