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
		fmt.Println("Failed to bind to port... 4221")
		os.Exit(1)
	}

	var conn net.Conn
	conn, err = l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}

	// Buffer to store incoming data
	buffer := make([]byte, 1024)

	resultBuffer, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading: ", err.Error())
		os.Exit(1)
	}

	rawRequest := string(buffer[:resultBuffer])

	path := strings.Split(rawRequest, " ")[1]

	// /index.html handler
	if path == "/" {
		conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
		conn.Write([]byte("<html><body><h1>Hello, World!</h1></body></html>"))
		return
	} else {
		conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
		conn.Write([]byte("<html><body><h1>404 Not Found</h1></body></html>"))
		return
	}

}
