package kadht

import (
	"encoding/binary"
	"fmt"
	"io"
	"time"
)

const (
	MSG_START = iota
	// All the message types identifiers
	// are listed here
	PING_REQ
	PING_RESP
	// End of all message types, nothing should go beyond this
	// Mark my words
	MSG_END
)

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
	// Forward the de-serialize call to the basic message
	ret := this.base_msg.Deserialize(reader)
	if !ret {
		fmt.Println("ERROR: Failed to deserialize PingRequest")
	}
	return ret
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
	// Forward the de-serialize call to the basic message
	ret := this.base_msg.Deserialize(reader)
	if !ret {
		fmt.Println("ERROR: Failed to deserialize PingReply")
	}
	return ret
}
