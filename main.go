package main

import (
	"fmt"
	"net"
	"strings"
)

type job struct {
	connection net.Conn
}

func statusText(code int) string {
	switch code {
	case 200:
		return "200 OK"
	case 403:
		return "403 Forbidden"
	case 404:
		return "404 Not Found"
	default:
		return "500 Internal Server Error"
	}
}

func createHttpResponse(code int, body string) []byte {
	return []byte(
	"HTTP/1.1 " + statusText(code) + "\r\n" +
	"Content-Type: text/html\r\n" +
	fmt.Sprintf("Content-Length: %d\r\n", len(body)) +
	"\r\n" +
	body,
	)
}

func parseHttpRequest(httpRequest string) (string, error) {
	for line := range strings.Lines(httpRequest) {
		words := strings.Fields(line)
		if len(words) != 3 {
			return "", fmt.Errorf("Wrong http format")
		}
		if words[0] != "GET" {
			return "", fmt.Errorf("Cannot parse different then GET")
		}

		return words[1], nil
	}
	return "", fmt.Errorf("Failed to parse HTTP")
}

func handlePath(path string) (int, string) {
	switch path {
	case "/path":
			return 200, "<html><h1>PATH</h1><p>You have made a GET request to /path!</p></html>"
	case "/admin":
		return 403, "<html><h1>403 FORBIDDEN</h1><p>You are not allowed to be here...</p></html>"
	case "/favicon.ico":
		return 404, ""
	case "/":
		return 200, "This is the GO server! Try out to navigate to /path"
	default:
		return 404, "Error 404, resource not found :("
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	fmt.Println("New connection from ", conn.RemoteAddr().String())

	buf := make([]byte, 1024)

	n, err := conn.Read(buf)
	if err != nil {
		return
	}

	path, err := parseHttpRequest(string(buf[:n]))
	if err != nil {
		fmt.Println("Error in the http parsing: " + err.Error())
		return
	}

	responseCode, body := handlePath(path)
	conn.Write(createHttpResponse(responseCode, body))
}

func worker(jobs chan job) {
	for job := range jobs {
		handleConnection(job.connection)
	}
}

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error in creating listener!")
	} else {
		fmt.Println("Successfully created Listener!")
	}

	jobs := make(chan job, 10)

	for i := 0; i < 3; i++ {
		go worker(jobs)
		fmt.Printf("Worker %d started\n", i)
	}

	for {
		conn, err := ln.Accept()

		if err != nil {
			fmt.Println("Failed to accept connection!")
		} else {
			fmt.Println("Successfully accepted connection!")
		}
		
		jobs <- job{conn}
	}
}