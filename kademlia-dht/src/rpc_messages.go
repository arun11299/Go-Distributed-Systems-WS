package kadht

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"time"
)

const (
	MSG_START = iota
	// All the message types identifiers
	// are listed here
	PING_REQ
	PING_RESP
	FIND_NODE_REQ
	FIND_NODE_RESP
	FIND_VALUE_REQ
	FIND_VALUE_RESP
	// End of all message types, nothing should go beyond this
	// Mark my words
	MSG_END
)

func MsgType2Str(mtype uint32) string {
	switch mtype {
	case PING_REQ:
		return "PING_REQ"
	case PING_RESP:
		return "PING_RESP"
	case FIND_NODE_REQ:
		return "FIND_NODE_REQ"
	case FIND_VALUE_REQ:
		return "FIND_VALUE_REQ"
	case FIND_VALUE_RESP:
		return "FIND_VALUE_RESP"
	case FIND_NODE_RESP:
		return "FIND_NODE_RESP"
	default:
		panic("Message Type not correct")
	}
}

const (
	// Max limit on the number of nodes a remote node can send
	// in a node reply message
	kNodes = 7
)

type Ipv4Addr struct {
	IP   [4]byte
	Port uint16
}

/*
 * RemoteNode : Stores the minimal required information
 * about the remote node.
 * TODO: Need some kind of correlation with the Node
 * present in routing_table.
 */
type RemoteNode struct {
	Id   NodeId // ID of the remote node
	Addr Ipv4Addr
}

func NewIpv4Addr(addr *net.UDPAddr) Ipv4Addr {
	var raddr Ipv4Addr
	// TODO: Assume IPv4 address. Would not work
	// properly with IPv6
	copy(raddr.IP[:], addr.IP[0:4])
	raddr.Port = uint16(addr.Port)
	return raddr
}

/*
 * Message interface that every struct implementing a
 * message type must satisfy.
 * Interface API's:
 * 1. Serialize :
 *    Takes a io.Writer and writes the structure into it
 *    in binary format.
 *    Returns 'true' if serialization is done successfully
 *    otherwise returns 'false'
 */
type IMessage interface {
	Serialize(io.Writer) bool
	Deserialize(io.Reader) bool
}

/*
 * Structure of a Basic Message
 */
type BasicMsgHeader struct {
	Version   uint32 // The message version. Hard coded to 1
	MsgType   uint32 // Type of the request or response message
	EpochTime int64  // Time at which message was created
	SenderId  NodeId // Node ID of the sender node
	RandomId  NodeId // Random ID for matching response with request context
}

/*
 * PingRequest
 */
type PingRequest struct {
	base_msg BasicMsgHeader
}

/*
 * PingReply
 */
type PingReply struct {
	base_msg BasicMsgHeader
}

/*
 * FindNodeRequest
 */
type FindNodeRequest struct {
	base_msg     BasicMsgHeader
	LookupNodeId NodeId // The ID of the node that we are looking for
}

/*
 * FindValueRequest
 */
type FindValueRequest struct {
	base_msg      BasicMsgHeader
	LookupValueId NodeId // The ID of the value taht we are looking for
}

/*
 * FindNodeReply
 */
type FindNodeReply struct {
	base_msg   BasicMsgHeader
	TotalNodes int32        // Total number of nodes in the message
	Nodes      []RemoteNode // List of nodes
}

/*
 * NewBasicMsgHeader: Creates the basic message header.
 * Parameters:
 * [in] msg_type : The type of the message to create
 * [in] sender_id : Node Id of the sending node.
 * [in] random_id : Random ID. Depends on request or reply message.
 * [out] *BasicMsgHeader : Pointer to the newly created BasicMsgHeader
 */
func NewBasicMsgHeader(msg_type uint32, sender_id, random_id NodeId) *BasicMsgHeader {
	if msg_type <= MSG_START && msg_type >= MSG_END {
		panic("Received incorrect message type: ")
	}
	now := time.Now()

	return &BasicMsgHeader{
		Version:   1,
		MsgType:   msg_type,
		EpochTime: now.Unix(),
		SenderId:  sender_id,
		RandomId:  random_id,
	}
}

/*
 * NewPingRequest : Creates a new Ping request.
 * Parameters:
 * [in] sender_id : Node Id of the sending node.
 * [out] *PingRequest : Pointer to the newly created PingRequest
 */
func NewPingRequest(sender_id NodeId) *PingRequest {
	// generate a new random_id
	return &PingRequest{
		base_msg: *NewBasicMsgHeader(PING_REQ, sender_id, generateRandomNodeId()),
	}
}

/*
 * NewPingReply : Create a new Ping response.
 * Parameters:
 * [in] sender_id : Node Id of the sending node.
 * [in] ping_req : Corresponding ping request.
 * [out] *PingReply : Pointer to the newly created PingReply
 */
func NewPingReply(sender_id NodeId, ping_req *PingRequest) *PingReply {
	// copy the random id from the request
	return &PingReply{
		base_msg: *NewBasicMsgHeader(PING_RESP, sender_id, ping_req.base_msg.RandomId),
	}
}

/*
 * NewFindNodeRequest : Create a new Find node request.
 * Parameters:
 * [in] sender_id : Node Id of the sending node.
 * [in] lookup_id : Id of the node to lookup.
 * [out] *FindNodeRequest : Pointer to the newly created FindNodeRequest
 */
func NewFindNodeRequest(sender_id, lookup_id NodeId) *FindNodeRequest {
	return &FindNodeRequest{
		base_msg:     *NewBasicMsgHeader(FIND_NODE_REQ, sender_id, generateRandomNodeId()),
		LookupNodeId: lookup_id,
	}
}

/*
 * NewFindValueRequest : Create a new Find value request.
 * Parameters:
 * [in] sender_id : Node Id of the sending node.
 * [in] lookup_id : Id of the value to lookup.
 * [out] *FindValueRequest : Pointer to the newly created FindValueRequest
 */
func NewFindValueRequest(sender_id, lookup_id NodeId) *FindValueRequest {
	return &FindValueRequest{
		base_msg:      *NewBasicMsgHeader(FIND_VALUE_REQ, sender_id, generateRandomNodeId()),
		LookupValueId: lookup_id,
	}
}

/*
 * NewFindNodeReply : Create a new Find node reply message.
 * Parameters:
 * [in] sender_id : Node Id of the sending node.
 * [in] nodes : The 'k' (atmax kNodes) close nodes
 * [in] find_node_req : The corresponding FindNodeRequest
 * [out] *FindNodeReply : Pointer to the newly created FindNodeReply
 */
func NewFindNodeReply(sender_id NodeId, nodes []RemoteNode,
	find_node_req *FindNodeRequest) *FindNodeReply {
	if len(nodes) > kNodes {
		fmt.Println("ERROR: More than allowed nodes present: ", len(nodes))
		return nil
	}
	find_node_reply := new(FindNodeReply)

	find_node_reply.base_msg = *NewBasicMsgHeader(
		FIND_NODE_RESP,
		sender_id,
		find_node_req.base_msg.RandomId,
	)

	find_node_reply.TotalNodes = int32(len(nodes))
	find_node_reply.Nodes = make([]RemoteNode, len(nodes))
	copy(find_node_reply.Nodes[:], nodes[:len(nodes)])

	return find_node_reply
}

//************** MESSAGE SERIALIZATION-DESERIALIZATION FUNCTIONS ***************//

/*
 * Serialize : Implementation of Serialize interface API for Basic message
 * type class
 * Parameters:
 * [in] writer : An io.Writer object
 * [out] bool : Returns 'true' if serialization was successfull otherwise 'false'
 */
func (this *BasicMsgHeader) Serialize(writer io.Writer) bool {
	//TODO: determine endian-ness in platform independent way.
	err := binary.Write(writer, binary.BigEndian, this)
	if err != nil {
		fmt.Println("ERROR: ", err)
		return false
	}
	return true
}

/*
 * Deserialize : Implementation of deserialize interface API for Basic message
 * type class
 * Parameters:
 * [in] reader : An io.Reader object
 * [out] bool : Returns 'true' if de-serialization was successfull otherwise 'false'
 */
func (this *BasicMsgHeader) Deserialize(reader io.Reader) bool {
	//TODO: determine endian-ness in a platform independent way.
	err := binary.Read(reader, binary.BigEndian, this)
	if err != nil {
		fmt.Println("ERROR: ", err)
		return false
	}
	return true
}

/*
 * Implementation of Serialize interface API for PingRequest
 * message type class
 */
func (this *PingRequest) Serialize(writer io.Writer) bool {
	// Forward the serialize call to the Basic message
	ret := this.base_msg.Serialize(writer)
	if !ret {
		fmt.Println("ERROR: Failed to serialize PingRequest")
	}
	return ret
}

/*
 * Implementation of the Deserialize interface API for PingRequest
 * message type class
 */
func (this *PingRequest) Deserialize(reader io.Reader) bool {
	// Header should have already been deserialized
	return true
}

/*
 * Implementation of Serialize interface API for PingReply
 * message type class
 */
func (this *PingReply) Serialize(writer io.Writer) bool {
	// Forward the serialize call to the Basic message
	ret := this.base_msg.Serialize(writer)
	if !ret {
		fmt.Println("ERROR: Failed to serialize PingReply")
	}
	return ret
}

/*
 * Implementation of Deserialize interface API for PingReply message type class
 */
func (this *PingReply) Deserialize(reader io.Reader) bool {
	// Header should have already been deserialized
	return true
}

/*
 * Implementation of Serialize interface API for FindNodeRequest
 * message type class
 */
func (this *FindNodeRequest) Serialize(writer io.Writer) bool {
	// First write the header by forwarding the call to basic message
	ret := this.base_msg.Serialize(writer)
	if !ret {
		fmt.Println("ERROR: Failed to serialize FindNodeRequest header")
		return false
	}
	// Serialize the lookup ID
	err := binary.Write(writer, binary.BigEndian, &this.LookupNodeId)
	if err != nil {
		fmt.Println("ERROR: Failed to serialize FindNodeRequest LookupNodeId")
		return false
	}
	return true
}

/*
 * Implementation of Deserialize interface API for FindNodeRequest message
 * type class
 */
func (this *FindNodeRequest) Deserialize(reader io.Reader) bool {
	// Header should have already been deserialized
	// Read the Lookup ID
	err := binary.Read(reader, binary.BigEndian, &this.LookupNodeId)
	if err != nil {
		fmt.Println("ERROR: Failed to deserialize FindNodeRequest LookupNodeId")
		return false
	}
	return true
}

/*
 * Implementation of Serialize interface API for FindValueRequest
 * message type class
 */
func (this *FindValueRequest) Serialize(writer io.Writer) bool {
	// First write the header by forwarding the call to basic message
	ret := this.base_msg.Serialize(writer)
	if !ret {
		fmt.Println("ERROR: Failed to serialize FindValueRequest header")
		return false
	}
	// Serialize the lookup ID
	err := binary.Write(writer, binary.BigEndian, &this.LookupValueId)
	if err != nil {
		fmt.Println("ERROR: Failed to serialize FindNodeRequest LookupValueId")
		return false
	}
	return true
}

/*
 * Implementation of Deserialize interface API for FindValueRequest message
 * type class
 */
func (this *FindValueRequest) Deserialize(reader io.Reader) bool {
	// Header should have already been deserialized
	// Read the lookup ID
	err := binary.Read(reader, binary.BigEndian, &this.LookupValueId)
	if err != nil {
		fmt.Println("ERROR: Failed to deserialize FindValueRequest LookupValueId")
		return false
	}
	return true
}

func (this *FindNodeReply) Serialize(writer io.Writer) bool {
	// First write the header by forwarding the call to basic message
	ret := this.base_msg.Serialize(writer)
	if !ret {
		fmt.Println("ERROR: Failed to serialize FindNodeReply header")
		return false
	}
	// Serialize Total nodes
	err := binary.Write(writer, binary.BigEndian, &this.TotalNodes)
	if err != nil {
		fmt.Println("ERROR: Failed to serialize FindNodeReply total nodes")
		return false
	}
	// Serialize the nodes
	for idx := range this.Nodes {
		err = binary.Write(writer, binary.BigEndian, &this.Nodes[idx])
		if err != nil {
			fmt.Println("ERROR: Failed to serialize FindNodeReply remote node: ", err)
			return false
		}
	}
	return true
}

func (this *FindNodeReply) Deserialize(reader io.Reader) bool {
	// Header should have already been deserialized
	// Read total nodes
	err := binary.Read(reader, binary.BigEndian, &this.TotalNodes)
	if err != nil {
		fmt.Println("ERROR: Failed to deserialize FindNodeReply total nodes")
		return false
	}
	// Read the nodes
	this.Nodes = make([]RemoteNode, this.TotalNodes)
	for idx := range this.Nodes {
		err = binary.Read(reader, binary.BigEndian, &this.Nodes[idx])
		if err != nil {
			fmt.Println("ERROR: Failed to deserialize FindNodeReply nodes")
			return false
		}
	}
	return true
}
