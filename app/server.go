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

	lines := strings.Split(rawRequest, "\r\n")

	headers := make(map[string]string)

	for _, line := range lines {
		fmt.Println("line: " + line)
		if len(line) < 1 {
			break
		}

		splittedLine := strings.Split(line, ": ")

		if len(splittedLine) == 2 {
			headers[splittedLine[0]] = splittedLine[1]
		}

	}

	fmt.Println(headers)

	path := strings.Split(rawRequest, " ")[1]

	splittedPath := strings.Split(path, "/")

	responseContentLength := 0

	if path == "/" {
		conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
		return
	}

	if path == "/user-agent" {
		fmt.Println("@" + headers["User-Agent"])
		responseContentLength = len(headers["User-Agent"])
		responseContentLengthString := fmt.Sprint(responseContentLength)
		conn.Write([]byte("HTTP/1.1 200 OK\r\n"))
		conn.Write([]byte("Content-Type: text/plain\r\n"))
		conn.Write([]byte("Content-Length: " + responseContentLengthString + "\r\n\r\n"))
		conn.Write([]byte(headers["User-Agent"]))
		return
	}

	if len(splittedPath) > 2 {
		conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
		return
	}

	if !strings.HasPrefix(path, "/echo/") {
		conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
		return
	}

	if strings.HasPrefix(path, "/echo/") {
		secondPath := splittedPath[2]

		conn.Write([]byte("HTTP/1.1 200 OK\r\n"))
		conn.Write([]byte("Content-Type: text/plain\r\n"))
		conn.Write([]byte("Content-Length: " + fmt.Sprint(len(secondPath)) + "\r\n\r\n"))
		conn.Write([]byte(secondPath))
		return
	}

}
