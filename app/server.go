package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

// Request struct to hold HTTP request information
type Request struct {
	Method  string
	Path    string
	Headers map[string]string
	Params  map[string]string
	Body    string
}

// Response struct to hold HTTP response information
type Response struct {
	StatusCode int
	Headers    map[string]string
	Body       string
}

// App struct to hold routes
type App struct {
	getRoutes  map[string]func(*Request, *Response)
	postRoutes map[string]func(*Request, *Response)
}

// expressgo function to start the server
func expressgo(port string, message func(string)) *App {
	l, err := net.Listen("tcp", "0.0.0.0:"+port)
	if err != nil {
		fmt.Println("Failed to bind to port... " + port)
		os.Exit(1)
	}

	message(port)

	app := &App{
		getRoutes:  make(map[string]func(*Request, *Response)),
		postRoutes: make(map[string]func(*Request, *Response)),
	}

	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				fmt.Println("Error accepting connection: ", err.Error())
				continue
			}
			go app.handleConnection(conn)
		}
	}()

	return app
}

// get function to add GET routes to the app
func (app *App) get(path string, handler func(*Request, *Response)) {
	app.getRoutes[path] = handler
}

// post function to add POST routes to the app
func (app *App) post(path string, handler func(*Request, *Response)) {
	app.postRoutes[path] = handler
}

// handleConnection function to handle incoming connections
func (app *App) handleConnection(conn net.Conn) {
	defer conn.Close()

	// Read the incoming request
	reader := bufio.NewReader(conn)
	requestLine, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading request: ", err.Error())
		return
	}

	// Parse the HTTP request line
	requestParts := strings.Fields(requestLine)
	if len(requestParts) < 2 {
		fmt.Println("Invalid request line: ", requestLine)
		return
	}

	method := requestParts[0]
	path := requestParts[1]

	// Read headers
	headers := make(map[string]string)
	for {
		line, err := reader.ReadString('\n')
		if err != nil || line == "\r\n" {
			break
		}
		headerParts := strings.SplitN(line, ": ", 2)
		if len(headerParts) == 2 {
			headers[headerParts[0]] = strings.TrimSpace(headerParts[1])
		}
	}

	// Parse parameters from the path
	params := make(map[string]string)
	var handler func(*Request, *Response)
	var ok bool

	if method == "GET" {
		for route, h := range app.getRoutes {
			if matchRoute(route, path, params) {
				handler = h
				ok = true
				break
			}
		}
	} else if method == "POST" {
		for route, h := range app.postRoutes {
			if matchRoute(route, path, params) {
				handler = h
				ok = true
				break
			}
		}
	}

	if !ok {
		// If route does not exist, send 404 response
		response := &Response{
			StatusCode: 404,
			Body:       "404 - Not Found",
			Headers:    make(map[string]string),
		}
		app.writeResponse(conn, response)
		return
	}

	// Create request object
	request := &Request{
		Method:  method,
		Path:    path,
		Headers: headers,
		Params:  params,
	}

	// Create response object
	response := &Response{
		Headers: make(map[string]string),
	}

	// Call the handler function for the route
	handler(request, response)

	// Write the response
	app.writeResponse(conn, response)
}

// matchRoute function to match routes and extract parameters
func matchRoute(route, path string, params map[string]string) bool {
	routeParts := strings.Split(route, "/")
	pathParts := strings.Split(path, "/")

	if len(routeParts) != len(pathParts) {
		return false
	}

	for i, routePart := range routeParts {
		if strings.HasPrefix(routePart, ":") {
			params[routePart[1:]] = pathParts[i]
		} else if routePart != pathParts[i] {
			return false
		}
	}

	return true
}

// writeResponse function to send HTTP response
func (app *App) writeResponse(conn net.Conn, response *Response) {

	var statusText string
	if response.StatusCode == 200 {
		statusText = "OK"
	} else if response.StatusCode == 201 {
		statusText = "Created"
	} else {
		statusText = "Not Found"

	}

	// Write status line
	conn.Write([]byte(fmt.Sprintf("HTTP/1.1 %d %s\r\n", response.StatusCode, statusText)))

	// Write headers
	for key, value := range response.Headers {
		conn.Write([]byte(fmt.Sprintf("%s: %s\r\n", key, value)))
	}

	// Write a blank line to indicate the end of headers
	conn.Write([]byte("\r\n"))

	// Write body
	conn.Write([]byte(response.Body))
}

// Main function to start the server and define routes
func main() {
	app := expressgo("4221", func(port string) {
		fmt.Println("Server started at: " + port)
	})

	app.get("/productlist", handleProductList)
	app.get("/productlist/apple", handleApple)
	app.post("/echo/:id", handleEcho)
	app.get("/echo/:id", handleEcho)
	app.post("/files/:content", handleFiles)

	select {}
}

// handleProductList function to handle /productlist route
func handleProductList(req *Request, res *Response) {
	fmt.Printf("Method: %s, Path: %s, Headers: %v, Body: %s\n", req.Method, req.Path, req.Headers, req.Body)
	res.StatusCode = 200
	res.Headers["Content-Type"] = "text/plain"
	res.Body = "Product List"
}

// handleApple function to handle /productlist/apple route
func handleApple(req *Request, res *Response) {
	res.StatusCode = 200
	res.Headers["Content-Type"] = "text/plain"
	res.Body = "Apple"
}

// handleEcho function to handle /echo/:id route
func handleEcho(req *Request, res *Response) {
	id := req.Params["id"]

	res.StatusCode = 200
	res.Headers["Content-Type"] = "text/plain"
	// Check if any compression is required
	hasGzip := strings.Contains(req.Headers["Accept-Encoding"], "gzip")
	if hasGzip {
		// Compress the response body
		res.Body = GzipCompress(id)
		res.Headers["Content-Length"] = fmt.Sprintf("%d", len(res.Body))
		res.Headers["Content-Encoding"] = "gzip"
		fmt.Println("One")
	} else {
		fmt.Println("Two")
		res.Body = fmt.Sprintf("Echo ID: %s", id)
	}
}

func handleFiles(req *Request, res *Response) {
	// We will create a new file and write the content a file from cli --directory flag
	filename := req.Params["content"]

	directory := os.Args[2]

	// Write the content to a file
	file, err := os.Create(directory + filename)
	if err != nil {
		res.StatusCode = 500
		res.Headers["Content-Type"] = "text/plain"
		res.Body = "Error creating file"
		return
	}
	// Write body to the content
	file.WriteString(req.Body)
	defer file.Close()

	res.StatusCode = 201
	res.Headers["Content-Type"] = "text/plain"
	res.Body = fmt.Sprintf("File content: %s", req.Body)
}
