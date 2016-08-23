package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"os"
)

type BaseMessage struct {
	Version uint32
	Magic   [4]byte
}

func start_client() {
	fmt.Println("Starting Client")
	serv_addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:6789")
	client_addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")

	conn, err := net.DialUDP("udp", client_addr, serv_addr)
	if err != nil {
		fmt.Println("DialUDP failed")
		os.Exit(1)
	}
	defer conn.Close()

	sendPkt := BaseMessage{
		Version: 42,
		Magic:   [4]byte{'A', 'r', 'u', 'n'},
	}

	err = binary.Write(conn, binary.BigEndian, &(sendPkt))
	if err != nil {
		fmt.Println("Binary write failed: ", err)
		os.Exit(1)
	}

	/*
		wrBytes, err := conn.Write(sendPkt.magic)
		fmt.Println("Written ", wrBytes)

		if err != nil {
			fmt.Println("Binary write failed for magic: ", err)
			os.Exit(1)
		}
	*/
	fmt.Println("Client sent data")
}

func main() {
	serv_addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:6789")
	if err != nil {
		fmt.Println("Error in net.ResolveUDPAddr")
		os.Exit(1)
	}

	uconn, err := net.ListenUDP("udp", serv_addr)
	if err != nil {
		fmt.Println("Error in net.ListenUDP")
		os.Exit(1)
	}

	var pktbuf BaseMessage

	// Start the cllient
	go start_client()

	fmt.Println("Server going to read")
	err = binary.Read(uconn, binary.BigEndian, &(pktbuf))
	if err != nil {
		fmt.Println("Error in binary read of n/w packet ", err)
		os.Exit(1)
	}

	//pktbuf.magic = make([]byte, 4)
	//io.ReadFull(uconn, pktbuf.magic)

	fmt.Println("Read packet version:", pktbuf.Version)
	fmt.Println("Read packet magic:", pktbuf.Magic)

}
