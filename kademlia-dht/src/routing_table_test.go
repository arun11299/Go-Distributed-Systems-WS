package kadht

import (
	"crypto/sha1"
	"fmt"
	"net"
	"strconv"
	"testing"
)

func TestBasicRouteTable(t *testing.T) {
	a := []byte("server-hash")
	hash_a := sha1.Sum(a)
	rt := NewRoutingTable(hash_a)
	for i := 0; i < 10; i++ {
		udp_addr := "127.0.0.1" + ":" + strconv.Itoa(i)
		addr, _ := net.ResolveUDPAddr("udp", udp_addr)
		a = []byte(udp_addr)
		hash_a = sha1.Sum(a)
		node := CreateNode(addr, hash_a)

		res := rt.AddEntryOnly(node)
		if !res {
			t.Error("Adding node entry failed")
		}
	}
}

func TestLookupClosestNode(t *testing.T) {
	serv_id := []byte("127.0.0.1:0")
	serv_hash := sha1.Sum(serv_id)
	rt := NewRoutingTable(serv_hash)

	// Populate some entries
	for i := 0; i < 1000; i++ {
		udp_addr := "127.0.0.1" + ":" + strconv.Itoa(i+1)
		addr, _ := net.ResolveUDPAddr("udp", udp_addr)
		tmp := []byte(udp_addr)
		tmp_hash := sha1.Sum(tmp)

		node := CreateNode(addr, tmp_hash)
		res := rt.AddEntryOnly(node)
		if !res {
			t.Error("Adding node entry failed: ", udp_addr)
		}
	}
	rt.printBriefStats()

	// Lookup for 127.0.0.1:500
	nid := sha1.Sum([]byte("127.0.0.1:500"))
	addresses := rt.LookupClosestNodes(nid)
	fmt.Println("Result set size = ", len(addresses))
	for _, e := range addresses {
		fmt.Println(e.String())
	}
}
