package main

import (
	"flag"
	"fmt"
	"io/fs"
	"net"
	"os"
	"strings"
)

func fileHandler(filePath, method, requestBody string) string {
	if method == "POST" {
		os.WriteFile(filePath, []byte(requestBody), fs.ModeAppend)
		return "HTTP/1.1 201 Created\r\n\r\n"
	}
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "HTTP/1.1 404 Not Found\r\n\r\n"
	}
	return fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: application/octet-stream\r\nContent-Length: %d\r\n\r\n%s",
		len(data), data)
}

func handleConnection(c net.Conn, dir string) {
	data := make([]byte, 2048)
	_, err := c.Read(data)
	if err != nil {
		fmt.Println("Error reading data: ", err.Error())
	}

	// parsing request
	str := string(data)
	lines := strings.Split(str, "\r\n")
	fields := strings.Fields(lines[0])
	path := fields[1]
	method := fields[0]
	trimmedPath := strings.Trim(path, "/")
	pathFields := strings.Split(trimmedPath, "/")

	// building response
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
	case "files":
		fileName := pathFields[1]
		filePath := fmt.Sprintf("%s/%s", dir, fileName)
		requestBody := lines[len(lines)-1]
		response = fileHandler(filePath, method, requestBody)
	default:
		response = "HTTP/1.1 404 Not Found\r\n\r\n"
	}

	_, err = c.Write([]byte(response))
	if err != nil {
		fmt.Println("Error writing data: ", err.Error())
	}

	c.Close()
}

func main() {
	dir := flag.String("directory", "", "directory name")
	flag.Parse()

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}
	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go handleConnection(c, *dir)
	}

}
