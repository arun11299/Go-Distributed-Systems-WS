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
	address        *net.UDPAddr   // Address of the remote node
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

func NewRoutingTable(server_node_id node_id.NodeId) *RoutingTable {
	&RoutingTable{
		server_id: server_node_id,
	}
}

func CreateNode(addr *net.UDPAddr, id node_id.NodeId) *Node {
	&Node{
		address:        addr,
		id:             id,
		lastAccessTime: time.Now(),
	}
}

func (this *RoutingTable) AddEntry(node *Node) bool {
	slot := node_id.commonBits(this.server_id, node.id)
	entry, found := this.findEntry(slot, node)

	if !found {
		if this.entries[slot].used == slotsPerBucket {
			//TODO: Someone should handle the eviction
			return false
		}
		this.entries[slot].append(node)
		this.entries[slot].used++
	} else {
		entry.lastAccessTime = time.Now()
	}
	return true
}

func (this *RoutingTable) RemoveEntry(node *Node) bool {
	slot := node_id.commonBits(this.server_id, node.id)
	index, found := this.findEntryByIndex(slot, node)

	if !found {
		return true
	}
	bucket := this.entries[slot]
	bucket = append(bucket[:index], bucket[index+1:]...)
	return true
}

func (this *RoutingTable) findEntry(slot uint32, node *Node) (*Node, bool) {
	bucket := this.entries[slot]
	if bucket.used == 0 {
		return nil, false
	}
	for i := 0; i < bucket.used; i++ {
		if node_id.Compare(bucket[i].id, node.id) {
			return bucket[i], true
		}
	}
	return nil, false
}

func (this *RoutingTable) findEntryByIndex(slot uint32, node *Node) (uint32, bool) {
	bucket := this.entries[slot]
	if bucket.used == 0 {
		return 0, false
	}
	for i := 0; i < bucket.used; i++ {
		if node_id.Compare(bucket[i].id, node.id) {
			return i, true
		}
	}
	return 0, false
}
