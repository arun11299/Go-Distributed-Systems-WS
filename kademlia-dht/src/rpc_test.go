package kadht

import (
	"fmt"
	"net"
	"testing"
	"time"
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
	fmt.Println("Wait for ping response")
	saddr, err := net.ResolveUDPAddr("udp", "127.0.0.2:6789")
	if err != nil {
		t.Error("Failed in resolve addr")
		return
	}
	uconn, err := net.ListenUDP("udp", serv_addr)
	if err != nil {
		t.Error("Failed in listen")
		return
	}

	msg, _ := ConsumePacket(conn)
	if msg == nil {
		fmt.Println("ERROR: Failed to read ping response")
		return
	}
	_, ok := msg.(*PingReply)
	if !ok {
		fmt.Println("ERROR: Failed to convert base message")
		return
	}
	fmt.Println("Received response for ping")
}

/*
func start_find_node() {
	fmt.Println("Starting find node Client")
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
		fmt.Println("Error in sending find node request")
		return
	}

	fmt.Println("Wait for Find node reply")
	msg, _ := ConsumePacket(conn)
	if msg == nil {
		fmt.Println("ERROR: Failed to read response")
		return
	}

	_, ok := msg.(*FindNodeReply)
	if !ok {
		fmt.Println("ERROR: Failed to convert base message")
		return
	}
	fmt.Println("Received response for find node")

}
*/

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
	msg, _ := ConsumePacket(uconn)
	if msg == nil {
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

	// Send ping response
	fmt.Println("Sending ping response")
	var ctx ServerConfig
	ctx.node_id = generateRandomNodeId()

	ret := SendPingResponse(uconn, omsg, &ctx)
	if !ret {
		fmt.Println("ERROR: Failed to send ping response")
		return
	}
	time.Sleep(time.Second * 10)
}

/*
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
	fmt.Println("Node ID: ", omsg.LookupNodeId)

	// Form the reply
	var ctx ServerConfig
	ctx.node_id = generateRandomNodeId()

	nodes := make([]RemoteNode, 4)
	for i := range nodes {
		nodes[i] = RemoteNode{
			Id:   generateRandomNodeId(),
			Addr: NewIpv4Addr(serv_addr),
		}
	}

	ret := SendFindNodeResponse(uconn, omsg, nodes, &ctx)
	if !ret {
		fmt.Println("Failed to send find node response")
		return
	}
	fmt.Println("Response sent")

	time.Sleep(time.Second * 10)

}
*/
