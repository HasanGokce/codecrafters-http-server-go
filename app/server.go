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

	// Get remote address (IP and port)
	remoteAddr := conn.RemoteAddr()
	fmt.Println("Connection from:", remoteAddr)

	// Buffer to store incoming data
	buffer := make([]byte, 1024)

	resultBuffer, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading: ", err.Error())
		os.Exit(1)
	}

	rawRequest := string(buffer[:resultBuffer])

	fmt.Println(rawRequest)

	path := strings.Split(rawRequest, " ")[1]

	splittedPath := strings.Split(path, "/")

	fmt.Println(splittedPath)

	fmt.Println(len(splittedPath))

	if len(splittedPath) < 3 {
		conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
		return
	}

	secondPath := splittedPath[2]

	fmt.Println(secondPath)

	if path == "/" {
		conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
		return
	}

	if !strings.HasPrefix(path, "/echo/") {
		conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
		return
	}

	if strings.HasPrefix(path, "/echo/") {
		conn.Write([]byte("HTTP/1.1 200 OK\r\n"))
		conn.Write([]byte("Content-Type: text/plain\r\n"))
		conn.Write([]byte("Content-Length: " + fmt.Sprint(len(secondPath)) + "\r\n\r\n"))
		conn.Write([]byte(secondPath))
		return
	}

}
