package main

import (
	"fmt"
	"net"
	"os"
)

func check_error(err error) {
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}
}

func handle_connection(conn net.Conn) {
	defer conn.Close()
	var buf [512]byte
	// Wait for data to be read from socket
	for {
		n, err := conn.Read(buf[0:])
		if n == 0 {
			fmt.Println("Received connection reset from the client")
			return
		}
		if err != nil {
			fmt.Printf("Error while receiving data: %s\n", err.Error())
			return
		}

		_, err = conn.Write(buf[0:n])
		if err != nil {
			fmt.Printf("Error while writing data: %s\n", err.Error())
			return
		}
	}
}

func main() {
	tcp_addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:6789")
	check_error(err)

	listener, err := net.ListenTCP("tcp", tcp_addr)
	check_error(err)

	for {
		client_conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error while accepting connection")
			continue
		}

		go handle_connection(client_conn)
	}
}
