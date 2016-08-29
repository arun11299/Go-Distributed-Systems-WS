package kadht

import (
	"fmt"
	"net"
	"runtime"
	"testing"
	//"time"
)

type ServerConfig struct {
	node_id [20]byte
}

type signal chan int

const (
	PING_REQ_SEND = iota
	PING_REQ_RECV
	SHUTDOWN
)

func start_node(event_chan signal,
	master_chan signal,
	lip, lport string,
	rip, rport string,
	node_name string) {

	fmt.Println("Starting node...")
	remote := rip + ":" + rport
	local := lip + ":" + lport

	serv_addr, err := net.ResolveUDPAddr("udp", local)
	if err != nil {
		fmt.Println("Error in resolve address ", local)
		return
	}
	uconn, err := net.ListenUDP("udp", serv_addr)
	if err != nil {
		fmt.Println("Error in net.ListenUDP: ", err)
		return
	}
	//defer uconn.Close()
	var ctx ServerConfig
	ctx.node_id = generateRandomNodeId()

	// Connect with the other node
	remote_addr, err := net.ResolveUDPAddr("udp", remote)
	if err != nil {
		fmt.Println("ERROR")
		return
	}
	cconn, err := net.DialUDP("udp", nil, remote_addr)
	if err != nil {
		fmt.Println("Could not dial remote ", err)
	}
	defer cconn.Close()

	for {
		fmt.Println("Waiting for event")
		what := <-event_chan
		switch what {
		case PING_REQ_SEND:

			fmt.Println("Sending ping request from ", node_name)
			ret := SendPingRequest(cconn, &ctx)
			if !ret {
				fmt.Println("Failed to send ping request ", node_name)
			}

		case PING_REQ_RECV:
			fmt.Println(node_name, " got ping request")

			msg, _ := ConsumePacket(uconn)
			if msg == nil {
				fmt.Println("ERROR: Failed to read ping response ", node_name)
			} else {
				ping_req, ok := msg.(*PingRequest)
				if !ok {
					fmt.Println("ERROR: failed to convert base message to ping reply ", node_name)
					return
				}
				fmt.Println("Version: ", ping_req.base_msg.Version)
				fmt.Println("MsgType: ", ping_req.base_msg.MsgType)
				fmt.Println("EpochTime: ", ping_req.base_msg.EpochTime)
				fmt.Println("SenderId: ", ping_req.base_msg.SenderId)
				fmt.Println("RandomId: ", ping_req.base_msg.RandomId)
			}

		default:
			fmt.Println("Shutting down ", node_name)
			master_chan <- SHUTDOWN
			return
		}
	}
}

func TestPingPong(t *testing.T) {
	runtime.GOMAXPROCS(5)
	ch1 := make(chan int, 1)
	ch2 := make(chan int, 1)
	my_chan := make(chan int, 2)

	// Start the two nodes
	go start_node(ch2, my_chan, "127.0.0.1", "6790", "127.0.0.1", "6789", "node-2")
	go start_node(ch1, my_chan, "127.0.0.1", "6789", "127.0.0.1", "6790", "node-1")

	// make node-1 send ping request
	ch1 <- PING_REQ_SEND
	// make node-2 read ping response
	ch2 <- PING_REQ_RECV

	//Shutdown both
	ch1 <- SHUTDOWN
	ch2 <- SHUTDOWN

	_ = <-my_chan
	_ = <-my_chan
}
