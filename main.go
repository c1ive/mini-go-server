package main

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
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

func openHtmlFile(name string) (string, error) {
	wd, _ := os.Getwd()
	path := filepath.Join(wd, "static", name)
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(content), nil
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
	getData := func(fileName string, code int) (int, string) {
		data, err := openHtmlFile(fileName)
		if err != nil {
			fmt.Println("openHtmlFile error:", err)
			return 500, "Failed to open HTML file"
		}
		return code, data
	}
	switch path {
	case "/path":
		return getData("path.html", 200)
	case "/admin":
		return getData("admin.html", 403)
	case "/":
		return getData("index.html", 200)
	default:
		return getData("404.html", 404)
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
		conn.Write(createHttpResponse(500, ""))
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