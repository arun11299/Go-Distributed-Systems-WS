package kadht

import (
	"fmt"
	"net"
	"testing"
)

type ServerConfig struct {
	node_id [20]byte
}

func start_ping_client() {
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

func start_find_node() {
	fmt.Println("Starting Client")
	serv_addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:6790")
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
	ret := SendFindNodeRequest(conn, generateRandomNodeId(), &ctx)
	if !ret {
		fmt.Println("Error in sending ping request")
		return
	}
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

	go start_ping_client()

	//Wait for Ping message
	msg, typ := ConsumePacket(uconn)
	if typ == -1 {
		fmt.Println("Wrong message type")
		return
	}
	//cast the message
	omsg, ok := msg.(*PingRequest)
	if !ok {
		fmt.Println("Type conversion failed")
		return
	}

	fmt.Println("Version: ", omsg.base_msg.Version)
	fmt.Println("Message Type: ", omsg.base_msg.MsgType)
	fmt.Println("Epoch Time: ", omsg.base_msg.EpochTime)
}

func TestFindNodeReqRep(t *testing.T) {
	serv_addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:6790")
	if err != nil {
		t.Error("Failed in resolve addr")
		return
	}
	uconn, err := net.ListenUDP("udp", serv_addr)
	if err != nil {
		t.Error("Failed in listen")
		return
	}

	go start_find_node()
	msg, typ := ConsumePacket(uconn)
	if typ == -1 {
		fmt.Println("Wrong message type")
		return
	}
	//cast the message
	omsg, ok := msg.(*FindNodeRequest)
	if !ok {
		fmt.Println("Type conversion failed")
		return
	}

	fmt.Println("Version: ", omsg.base_msg.Version)
	fmt.Println("Message Type: ", omsg.base_msg.MsgType)
	fmt.Println("Epoch Time: ", omsg.base_msg.EpochTime)

}
