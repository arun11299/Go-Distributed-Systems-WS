package main

import (
	"net"
	"os"
	"sync"
)

func read_routine(conn net.Conn) {
	defer conn.Close()
	var buf [512]byte
	for {
		n, err := conn.Read(buf[0:])
		if n == 0 {
			break
		}
		if err != nil {
			break
		}
	}
}

func write_routine(conn net.Conn) {
	defer conn.Close()
	for {
		conn.Write([]byte("Hello"))
	}
}

func handle_server_conn(conn net.Conn) {
	var wg sync.WaitGroup

	go read_routine(conn)
	go write_routine(conn)

	wg.Add(2)
	wg.Wait()
}

func main() {
	serv_addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:6789")
	if err != nil {
		os.Exit(1)
	}
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		conn, err := net.DialTCP("tcp", nil, serv_addr)
		if err != nil {
			os.Exit(1)
		}
		wg.Add(1)
		go handle_server_conn(conn)
	}

	wg.Wait()
}
