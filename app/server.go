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

	for {
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
		lines := strings.Split(str, "\r\n")
		path := strings.Fields(lines[0])[1]
		trimmedPath := strings.Trim(path, "/")
		pathFields := strings.Split(trimmedPath, "/")

		var response string
		switch pathFields[0] {
		case "":
			response = "HTTP/1.1 200 OK\r\n\r\n"
		case "echo":
			i := strings.Index(trimmedPath, "/")
			randomString := trimmedPath[i+1:]
			response = fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s",
				len(randomString), randomString)
		case "user-agent":
			userAgent := strings.Fields(lines[2])[1]
			response = fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s",
				len(userAgent), userAgent)
		default:
			response = "HTTP/1.1 404 Not Found\r\n\r\n"
		}

		_, err = c.Write([]byte(response))
		if err != nil {
			fmt.Println("Error writing data: ", err.Error())
		}

		c.Close()
	}
}
