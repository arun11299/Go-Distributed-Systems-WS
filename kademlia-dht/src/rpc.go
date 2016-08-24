package kadht

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
)

/*
 * ReadMessageHeader : Reads the Basic message header from the connection.
 * Parameters:
 * [in] resp_reader : io.reader object to read bytes from.
 * [out] BasicMsgHeader : The parsed message header.
 * [out] error : If any during parsing.
 */
func ReadMessageHeader(resp_reader io.Reader) (BasicMsgHeader, error) {
	var msg_header BasicMsgHeader

	ret := msg_header.Deserialize(resp_reader)
	if !ret {
		err := errors.New("Failed to read message header from n/w")
		return msg_header, err
	}

	return msg_header, nil
}

/*
 * ConsumePacket : Consumes the packet and abstracts lots
 * of details of packet parsing.
 * This is the function that must be called for complete message
 * parsing.
 * Parameters:
 * [in] conn : The connection channel (UDP) from where to read bytes.
 * [out] IMessage : The message class type implementing IMessage interface.
 * [out] int : The message type
 *
 * TODO: This is an absolutely horrible function. Do something!!
 */
func ConsumePacket(conn net.Conn) (IMessage, int) {
	resp_reader := bufio.NewReader(conn)

	header, err := ReadMessageHeader(resp_reader)
	if err != nil {
		fmt.Println("ERROR: Parsing failed")
		return nil, -1
	}

	switch header.MsgType {
	case PING_REQ:
		ping_req := new(PingRequest)
		ping_req.base_msg = header
		ret := ping_req.Deserialize(resp_reader)
		if !ret {
			return nil, -1
		}
		return ping_req, PING_REQ

	case PING_RESP:
		ping_resp := new(PingReply)
		ping_resp.base_msg = header
		ret := ping_resp.Deserialize(resp_reader)
		if !ret {
			return nil, -1
		}
		return ping_resp, PING_RESP

	case FIND_NODE_REQ:
		find_node_req := new(FindNodeRequest)
		find_node_req.base_msg = header
		ret := find_node_req.Deserialize(resp_reader)
		if !ret {
			return nil, -1
		}
		return find_node_req, FIND_NODE_REQ

	case FIND_NODE_RESP:
		find_node_resp := new(FindNodeReply)
		find_node_resp.base_msg = header
		ret := find_node_resp.Deserialize(resp_reader)
		if !ret {
			return nil, -1
		}
		return find_node_resp, FIND_NODE_RESP

	case FIND_VALUE_REQ:
		find_value_req := new(FindValueRequest)
		find_value_req.base_msg = header
		ret := find_value_req.Deserialize(resp_reader)
		if !ret {
			return nil, -1
		}
		return find_value_req, FIND_VALUE_REQ

	default:
		panic("Invalid message type")
	}
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

func SendFindNodeRequest(conn net.Conn, lookup_id NodeId, server_ctx *ServerConfig) bool {
	find_node_req := NewFindNodeRequest(server_ctx.node_id, lookup_id)
	req_writer := bufio.NewWriter(conn)

	ret := find_node_req.Serialize(req_writer)
	if !ret {
		fmt.Println("ERROR: Failed to send find node request")
		return false
	}

	req_writer.Flush()
	return true
}

func SendFindNodeResponse(conn net.Conn, find_node_req *FindNodeRequest,
	nodes []RemoteNode, server_ctx *ServerConfig) bool {
	find_node_resp := NewFindNodeReply(server_ctx.node_id, nodes, find_node_req)
	if find_node_resp == nil {
		fmt.Println("ERROR: Failed to create find node reply")
		return false
	}
	resp_writer := bufio.NewWriter(conn)

	ret := find_node_resp.Serialize(resp_writer)
	if !ret {
		fmt.Println("ERROR: Failed to serialize find node response")
		return false
	}

	resp_writer.Flush()
	return true
}
