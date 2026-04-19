package main

import (
	"fmt"
	"net"
	"os"
)

func listen(network string, address string) net.Listener {
	l, err := net.Listen(network, address)
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	return l
}

func accept(l net.Listener) net.Conn {
	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}
	return conn
}

func read_conn(conn net.Conn) {
	conn.Write([]byte("+PONG\r\n"))
}
