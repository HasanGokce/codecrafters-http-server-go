package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func getHeaders(rawRequest string) map[string]string {
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
	return headers
}

func handleConnection(conn net.Conn) {
	// Buffer to store incoming data
	buffer := make([]byte, 1024)

	resultBuffer, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading: ", err.Error())
		os.Exit(1)
	}

	rawRequest := string(buffer[:resultBuffer])

	headers := getHeaders(rawRequest)

	path := strings.Split(rawRequest, " ")[1]

	requestType := strings.Split(rawRequest, " ")[0]

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

	if strings.HasPrefix(path, "/echo/") {
		secondPath := splittedPath[2]

		conn.Write([]byte("HTTP/1.1 200 OK\r\n"))
		conn.Write([]byte("Content-Type: text/plain\r\n"))
		conn.Write([]byte("Compresion: gzip\r\n"))
		conn.Write([]byte("Content-Length: " + fmt.Sprint(len(secondPath)) + "\r\n\r\n"))
		conn.Write([]byte(secondPath))
		return
	}

	if strings.HasPrefix(path, "/files/") && requestType == "POST" {

		fileName := splittedPath[2]
		directory := os.Args[2]

		file, err := os.Create(directory + fileName)
		if err != nil {
			conn.Write([]byte("HTTP/1.1 500 Internal Server Error\r\n\r\n"))
			return
		}

		body := strings.Split(rawRequest, "\r\n\r\n")[1]

		file.Write([]byte(body))
		file.Close()

		fmt.Println("File created")

		conn.Write([]byte("HTTP/1.1 201 Created\r\n"))
		conn.Write([]byte("Content-Type: text/plain\r\n"))
		conn.Write([]byte("Content-Length: " + fmt.Sprint(len(body)) + "\r\n\r\n"))
		conn.Write([]byte("HTTP/1.1 " + body))

		fmt.Println(body)

		return
	}

	if strings.HasPrefix(path, "/files/") {
		fileName := splittedPath[2]
		directory := os.Args[2]

		file, err := os.ReadFile(directory + fileName)

		if err != nil {
			conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
			return
		}

		conn.Write([]byte("HTTP/1.1 200 OK\r\n"))
		conn.Write([]byte("Content-Type: application/octet-stream\r\n"))
		conn.Write([]byte("Content-Length: " + fmt.Sprint(len(file)) + "\r\n\r\n"))

		conn.Write(file)

		return
	}

	if len(splittedPath) < 3 {
		conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
		return
	}

	conn.Write([]byte("HTTP/1.1 200 OK\r\n"))
	conn.Write([]byte("Content-Type: text/plain\r\n\r\n"))
}

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port... 4221")
		os.Exit(1)
	}

	for {
		var conn net.Conn
		conn, err = l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go handleConnection(conn)
	}

}
