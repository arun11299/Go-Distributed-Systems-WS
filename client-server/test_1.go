package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s ip-addr\n", os.Args[0])
		os.Exit(1)
	}
	ip_addr_str := os.Args[1]
	addr := net.ParseIP(ip_addr_str)

	if addr == nil {
		fmt.Println("Invalid address passed ", ip_addr_str)
		os.Exit(1)
	} else {
		fmt.Printf("Parsed ip address is : %s\n", addr.String())
	}

	resolved_addr, err := net.ResolveIPAddr("ip", "www.google.com")
	if err != nil {
		fmt.Printf("DNS resolution failed\n")
		os.Exit(1)
	}
	fmt.Printf("Resolved address = %s\n", resolved_addr.String())

	addrs, err_2 := net.LookupHost("www.google.com")
	if err_2 != nil {
		fmt.Printf("Lookup failed\n")
		os.Exit(1)
	}
	for _, s := range addrs {
		fmt.Printf("%s\n", s)
	}

	tcp_addr, err_3 := net.ResolveTCPAddr("tcp4", "1.2.3.4:8080")
	if err_3 != nil {
		fmt.Println("net.ResolveTCPAddr failed!")
		os.Exit(1)
	}

	fmt.Printf("TCP address = %s\n", tcp_addr.String())
}
