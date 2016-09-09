package kadht

import (
	"crypto/sha1"
	"os"
	"time"
)

const (
	bytesPerNodeiId = 20
)

type NodeId [bytesPerNodeiId]byte

// Generate a random Node ID
// This is particularly useful for
// a newly started node when it does not
// have any ID. Otherwise for reboot/restart cases
// the node id once generated must be saved in a
// persistent medium.
func generateRandomNodeId() NodeId {
	host_name, err := os.Hostname()
	if err != nil {
		panic("OS hostname not found!!")
	}
	data := []byte(time.Now().String() + host_name)
	hash := sha1.Sum(data)
	return hash
}

// Calculates the number of common bits
// between 2 node Id's. Based upon the
// number of common bits, a bucket in the
// routing table is assigned.
func commonBits(a, b NodeId) int {
	var i int
	common_bits := 0

	for i = 0; i < bytesPerNodeiId; i++ {
		if a[i] != b[i] {
			break
		}
	}
	// matching bytes mutiplied by bits per bytes
	if i == bytesPerNodeiId {
		return bytesPerNodeiId * 8
	}

	common_bits = i * 8
	// Now find if there are any matching bits
	// at the ith index byte which didnt compare
	// wholly
	j := 0
	res := a[i] ^ b[i]
	for (res & 0x80) == 0 {
		res <<= 1
		j++
	}
	common_bits += j

	return common_bits
}

// Compares 2 node Id's
// If they are equal, returns true
// Otherwise returns false
func Compare(a, b NodeId) bool {
	res := commonBits(a, b)
	if res == bytesPerNodeiId*8 {
		return true
	}
	return false
}

func Xor(a, b NodeId) (res NodeId) {
	for i := 0; i < bytesPerNodeiId; i++ {
		res[i] = a[i] ^ b[i]
	}
	return res
}

func toString(a NodeId) string {
	return string(a[:])
}
