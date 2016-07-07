package routing_table

import (
	"net"
	"time"
)

const (
	numBitsID      = 160
	numBuckets     = 160
	slotsPerBucket = 20
)

type NodeId string

type Node struct {
	addr           *net.UDPAddr // Address of the remote node
	id             NodeId       // 20 byte ID of the remote node
	lastAccessTime time.Time    // Last looked up by current node
}

type Bucket struct {
	nodes [slotsPerBucket]Node
}

type RoutingTable struct {
	slots [numBuckets]Bucket
}
