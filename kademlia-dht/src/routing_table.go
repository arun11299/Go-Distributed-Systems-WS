package kadht

import (
	"net"
	"time"
)

const (
	// Number of bits in node id
	numBitsID = bytesPerNodeiId * 8
	// Number of buckets in the routing table
	numBuckets = numBitsID
	// Number of entries per bucket
	entriesPerBucket = 20
	//Number of nodes (Max capped)
)

// Represents a node in the distributed
// network
type Node struct {
	address        *net.UDPAddr // Address of the remote node
	id             NodeId       // 20 byte ID of the remote node
	lastAccessTime time.Time    // Last looked up by current node
}

// Represents a single bucket in the routing table
type Bucket struct {
	entries []Node
	used    int
}

type RoutingTable struct {
	server_id NodeId
	slots     []Bucket
}

// Create a new Routing Table
func NewRoutingTable(snode_id NodeId) *RoutingTable {
	return &RoutingTable{
		server_id: snode_id,
		slots:     make([]Bucket, numBuckets),
	}
}

func CreateNode(addr *net.UDPAddr, id NodeId) *Node {
	return &Node{
		address:        addr,
		id:             id,
		lastAccessTime: time.Now(),
	}
}

// AddEntryOnly: Adds a Node to the routing table if not present.
// If the entry already exists, just updates its access time.
// Parameters:
// [in] node : The node to be added.
// [out] bool : Returns 'true' if node gets added or already present.
//              'false' if the bucket is full.
//
func (this *RoutingTable) AddEntryOnly(node *Node) bool {
	slot := commonBits(this.server_id, node.id)
	entry, found := this.findEntry(slot, node.id)

	if !found {
		if this.slots[slot].used == entriesPerBucket {
			//TODO: Someone should handle the eviction
			return false
		}
		this.slots[slot].entries = append(this.slots[slot].entries, *node)
		this.slots[slot].used++
	} else {
		entry.lastAccessTime = time.Now()
	}
	return true
}

// RemoveEntry: Removes an entry from the Routing Table.
// Parameters:
// [in] node : The node to be removed
// [out] bool : 'true' if successfully removed, 'false' otherwise
//
func (this *RoutingTable) RemoveEntry(node *Node) bool {
	slot := commonBits(this.server_id, node.id)
	index, found := this.findEntryByIndex(slot, node.id)

	if !found {
		return true
	}
	bucket := this.slots[slot].entries
	bucket = append(bucket[:index], bucket[index+1:]...)
	return true
}

// findEntry: Finds an entry in the routing table.
// Parameters:
// [in] slot : The slot index or the bucket ID.
// [in] node : The Node ID to find
// [out] Node* : The found node.
// [out] bool : 'true' if node was found, 'false' otherwise
//
func (this *RoutingTable) findEntry(slot int, id NodeId) (*Node, bool) {
	bucket := this.slots[slot]
	if bucket.used == 0 {
		return nil, false
	}
	for i := 0; i < bucket.used; i++ {
		if Compare(bucket.entries[i].id, id) {
			return &bucket.entries[i], true
		}
	}
	return nil, false
}

// findEntryByIndex: Finds the index of the node within a bucket
// Parameters:
// [in] slot : The slot index or the bucket ID.
// [in] id : The Node ID to find
// [out] int : The index of the node.
// [out] bool : 'true' if node was found, 'false' otherwise
//
func (this *RoutingTable) findEntryByIndex(slot int, id NodeId) (int, bool) {
	bucket := this.slots[slot]
	if bucket.used == 0 {
		return 0, false
	}
	for i := 0; i < bucket.used; i++ {
		if Compare(bucket.entries[i].id, id) {
			return i, true
		}
	}
	return 0, false
}
