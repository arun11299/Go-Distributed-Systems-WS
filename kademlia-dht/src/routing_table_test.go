package kadht

import (
	"crypto/sha1"
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
