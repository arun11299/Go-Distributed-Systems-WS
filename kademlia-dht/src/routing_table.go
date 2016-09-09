package kadht

import (
	"fmt"
	"net"
	"time"
)

const (
	// Number of bits in node id
	numBitsID = bytesPerNodeiId * 8
	// Number of buckets in the routing table
	numBuckets = numBitsID + 1
	// Number of entries per bucket
	entriesPerBucket = 10000 // (Also known as k-bucket in kademlia literature)
	// Max number of nodes a node can respond to
	// find-node query
	alphaNodes = 7
)

// Represents a node in the distributed
// network
type Node struct {
	address        net.UDPAddr // Address of the remote node
	id             NodeId      // 20 byte ID of the remote node
	lastAccessTime time.Time   // Last looked up by current node
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
		address:        *addr,
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

// LookupClosestNodes: Finds the closest 'alphaNodes' number of nodes to
// the provided lookup ID.
// Parameters:
// [in] lookup_id : The ID that needs to be looked up
// [out] []NodeId : List of upto 'alphaNodes' number of Nodes
//
func (this *RoutingTable) LookupClosestNodes(lookup_id NodeId) []net.UDPAddr {
	var result []net.UDPAddr

	if lookup_id == this.server_id {
		return result
	}

	slot := commonBits(this.server_id, lookup_id)
	var copy_len int
	result = make([]net.UDPAddr, alphaNodes)

	copy_to_result := func(from_index, slot_num, num int) {
		for idx := 0; idx < num; idx++ {
			result[from_index+idx] = this.slots[slot_num].entries[idx].address
		}
	}

	if this.slots[slot].used >= alphaNodes {
		copy_to_result(0, slot, alphaNodes)
	} else {
		copy_len = this.slots[slot].used
		copy_to_result(0, slot, copy_len)

		remaining_nodes := alphaNodes - copy_len
		xor_res := Xor(this.server_id, lookup_id)
		/*
		 * The node finding algorithm is same as that of explained in
		 * a thesis written by 'Bruno Spori'.
		 * The basic idea is to first search the k-Bucket corresponding
		 * to the 'set' bits in the xor of current and lookup Node Id's.
		 * For Eg: If the result of Xor between 2 node id's is
		 * '00010101', then:
		 * 1. First search bucket 3, followed by bucket 5 and then bucket 7.
		 * 2. Then search 0, then 1,2,4,6 buckets
		 */
		// Byte iteration
		visited_set := make(map[int]bool)
		for i := slot / 8; i < bytesPerNodeiId && copy_len < alphaNodes; i++ {
			// Bit iteration
			for j := 0; j < 8; j++ {
				if (xor_res[i] & 0x80) != 0 {
					// Look up at slot 'i+j'
					tmp_slot := this.slots[i+j]
					visited_set[i+j] = true

					if tmp_slot.used >= remaining_nodes {
						copy_to_result(copy_len, i+j, remaining_nodes)
						break
					} else {
						copy_to_result(copy_len, i+j, tmp_slot.used)
						copy_len += tmp_slot.used
						remaining_nodes = alphaNodes - copy_len
					}
				}
			} // end bit iteration
		} //end byte iteration

		if copy_len < alphaNodes {
			for i := 0; i < bytesPerNodeiId*8 && copy_len < alphaNodes; i++ {
				if _, found := visited_set[i]; !found {
					tmp_slot := this.slots[i]
					if tmp_slot.used >= remaining_nodes {
						copy_to_result(copy_len, i, remaining_nodes)
						break
					} else {
						copy_to_result(copy_len, i, tmp_slot.used)
						copy_len += tmp_slot.used
						remaining_nodes = alphaNodes - copy_len
					}
				}
			} // end for
		} // end if
	}
	return result
}

func (this *RoutingTable) printBriefStats() {
	for i := 0; i < bytesPerNodeiId; i++ {
		fmt.Println(i, " : ", this.slots[i].used)
	}
}
