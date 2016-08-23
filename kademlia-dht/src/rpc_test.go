package kadht

import (
	"fmt"
	"net"
	"testing"
)

type ServerConfig struct {
	node_id [20]byte
}

func start_client() {
	fmt.Println("Starting Client")
	serv_addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:6789")
	client_addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
	conn, err := net.DialUDP("udp", client_addr, serv_addr)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	//Send a ping request
	var ctx ServerConfig
	ctx.node_id = generateRandomNodeId()
	ret := SendPingRequest(conn, &ctx)
	if !ret {
		fmt.Println("Error in sending ping request")
		return
	}

	fmt.Println("Client sent ping request")
}

func TestPingPong(t *testing.T) {
	serv_addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:6789")
	if err != nil {
		t.Error("Failed in resolve addr")
		return
	}
	uconn, err := net.ListenUDP("udp", serv_addr)
	if err != nil {
		t.Error("Failed in listen")
		return
	}

	go start_client()

	//Wait for Ping message
	hdr, err := ReadMessageHeader(uconn)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Version: ", hdr.Version)
	fmt.Println("Message Type: ", hdr.MsgType)
	fmt.Println("Epoch Time: ", hdr.EpochTime)
}
