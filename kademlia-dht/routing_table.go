package routing_table

import (
	"net"
	"node_id"
	"time"
)

const (
	numBitsID      = node_id.bytesPerNodeiId * 8
	numBuckets     = numBitsID
	slotsPerBucket = 20
)

type Node struct {
	addr           *net.UDPAddr   // Address of the remote node
	id             node_id.NodeId // 20 byte ID of the remote node
	lastAccessTime time.Time      // Last looked up by current node
}

type Bucket struct {
	entries [slotsPerBucket]Node
	used    int
}

type RoutingTable struct {
	server_id node_id.NodeId
	slots     [numBuckets]Bucket
}

func NewRoutingTable() *RoutingTable {
	&RoutingTable{}
}

func CreateNode(addr *net.UDPAddr, id node_id.NodeId) *Node {
	&Node{addr: addr, id: id}
}

func (this *RoutingTable) AddEntry(node *Node) {
	dist := node_id.Distance(this.server_id, node.id)
	slot := dist % numBuckets

	entry, found := this.findEntry(slot, node)
	if !found {
	}
}

func (this *RoutingTable) findEntry(slot uint32, node *Node) (*Node, bool) {
	bucket := this.entries[slot]
	if bucket.used == 0 {
		return nil, false
	}
	for i := 0; i < used; i++ {
		if node_id.Compare(bucket[i].id, node.id) {
			return bucket[i], true
		}
	}
	return nil, false
}
