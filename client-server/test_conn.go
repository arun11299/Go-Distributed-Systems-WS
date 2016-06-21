package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Wrong number of arguments")
		os.Exit(1)
	}
	hostname := os.Args[1]
	ip_addr, err := net.ResolveIPAddr("ip", hostname)
	check_error(err)

	var tcp_addr net.TCPAddr
	tcp_addr.IP = ip_addr.IP
	tcp_addr.Port = 80

	conn, err := net.DialTCP("tcp", nil, &tcp_addr)
	check_error(err)

	_, err = conn.Write([]byte("HEAD / HTTP/1.0\r\n\r\n"))
	check_error(err)

	result, err := ioutil.ReadAll(conn)
	check_error(err)

	fmt.Println(string(result))

}

func check_error(err error) {
	if err != nil {
		fmt.Printf("Error occurred: %s\n", err.Error())
		os.Exit(1)
	}
}
