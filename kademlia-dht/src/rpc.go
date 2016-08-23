package kadht

import (
	"bufio"
	"errors"
	"fmt"
	"net"
)

func ReadMessageHeader(conn net.Conn) (BasicMsgHeader, error) {
	var msg_header BasicMsgHeader
	resp_reader := bufio.NewReader(conn)

	ret := msg_header.Deserialize(resp_reader)
	if !ret {
		err := errors.New("Failed to read message header from n/w")
		return msg_header, err
	}

	return msg_header, nil
}

func SendPingRequest(conn net.Conn, server_ctx *ServerConfig) bool {
	ping_req := NewPingRequest(server_ctx.node_id)
	req_writer := bufio.NewWriter(conn)

	ret := ping_req.Serialize(req_writer)
	if !ret {
		fmt.Println("ERROR: Failed to send ping request")
		return false
	}

	req_writer.Flush()
	return true
}

func SendPingResponse(conn net.Conn, ping_req *PingRequest, server_ctx *ServerConfig) bool {
	ping_resp := NewPingReply(server_ctx.node_id, ping_req)
	resp_writer := bufio.NewWriter(conn)

	ret := ping_resp.Serialize(resp_writer)
	if !ret {
		fmt.Println("ERROR: Failed to send ping response")
		return false
	}

	resp_writer.Flush()
	return true
}
