package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
)

// Server context structure
type Server struct {
	port    int
	ip_addr string
	// On input event callback
	on_message_receive_cb func([]byte)
	// Internal Members
	addr *net.UDPAddr
	conn *net.UDPConn
}

// Client context structure
type Client struct {
	client_addr        *net.UDPAddr
	port               int
	on_message_send_cb func()
}

// Creates a new server instance
func NewServer(addr string, port int) *Server {
	server := &Server{ip_addr: addr, port: port}
	server.ServerSetup()
	return &server
}

/// Setups the UDP server
// This function is intended to be run
// on an independent go-routine
func (this *Server) ServerSetup() {
	udp_addr := ip_addr + ":" + strconv.Atoi(port)
	addr, err := net.ResolveUDPAddr("udp", udp_addr)
	if err != nil {
		panic("Resolve UDP address failed for server")
	}
	this.addr = addr
}

// Listen for UDP connections
// Called while setting up server
func (this *Server) Serve() {
	conn, err := net.ListenUDP("udp", this.addr)
	if err != nil {
		panic("Failed in ListenUDP")
	}
	this.conn = conn

	// Wait for incoming data
	for {
		this.WaitForInputPacket()
	}
}

func (this *Server) WaitForInputPacket() {
	var buf [1024]byte
	_, addr, err := this.conn.ReadFromUDP(buf[0:])
	if err != nil {
		//TODO: Logging
		return
	}

	if this.on_message_receive_cb != nil {
		go this.on_message_receive_cb(buf)
	}
}

// Dials into the provided address
// and returns a UDP connection
func connect() {
}
