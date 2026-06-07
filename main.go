package main

import (
	"fmt"
	"net"
)

func handleConnection(conn net.Conn) {
	
}

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error in creating listener!")
	} else {
		fmt.Println("Successfully created Listener!")
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Failed to accept connection!")
		} else {
			fmt.Println("Successfully accepted connection!")
		}
		go handleConnection(conn)
	}
}