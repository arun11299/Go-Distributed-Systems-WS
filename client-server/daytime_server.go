package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

func check_error(err error) {
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}
}

func serv_daytime(conn net.Conn) (err error) {
	defer conn.Close()
	daytime := time.Now().String()
	_, err = conn.Write([]byte(daytime))
	return err
}

func main() {
	tcp_addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:6789")
	check_error(err)

	listener, err := net.ListenTCP("tcp", tcp_addr)
	check_error(err)

	for {
		client_conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Error while accepting connection\n")
			continue
		}

		err = serv_daytime(client_conn)
		if err != nil {
			fmt.Printf("Failed to get time: %s\n", err.Error())
			continue
		}
	}
}
