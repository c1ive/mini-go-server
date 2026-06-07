package main

import (
	"fmt"
	"net"
	"strings"
)

func createHttpResponse(body string) []byte {
	return []byte(
	"HTTP/1.1 200 OK\r\n" +
	"Content-Type: text/plain\r\n" +
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

func handlePath(path string) (string) {
	switch path {
	case "/path":
			return "You have made a GET request to /path!"
	case "/":
		return "This is the GO server! Try out to navigate to /path"
	default:
		return "Error 404, resource not found :("
	}
}

func handleConnection(conn net.Conn, num int) {
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

	body := handlePath(path)
	conn.Write(createHttpResponse(body))
}

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error in creating listener!")
	} else {
		fmt.Println("Successfully created Listener!")
	}

	connections := 0

	for {
		conn, err := ln.Accept()
		connections++
		if err != nil {
			fmt.Println("Failed to accept connection!")
		} else {
			fmt.Println("Successfully accepted connection!")
		}
		go handleConnection(conn, connections)
	}
}