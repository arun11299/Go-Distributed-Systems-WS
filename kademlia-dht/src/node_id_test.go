package kadht

import (
	"crypto/sha1"
	"fmt"
	"testing"
)

func TestGenerateNodeId(t *testing.T) {
	for i := 0; i < 1000; i++ {
		id := generateRandomNodeId()
		if len(id) != 20 {
			t.Error("Incorrect Node ID generated")
		}
	}
}

func TestDistance(t *testing.T) {
	a := []byte("Same string")
	b := []byte("Same string")
	hash_a := sha1.Sum(a)
	hash_b := sha1.Sum(b)
	fmt.Println("hash_a = ", hash_a)
	fmt.Println("hash_b = ", hash_b)

	dist := commonBits(hash_a, hash_b)
	if dist != 160 {
		t.Error("Distance is not Zero for matching entries")
	}
	new_b := []byte("same string")
	hash_b = sha1.Sum(new_b)
	fmt.Println("New hash_b = ", hash_b)

	dist = commonBits(hash_a, hash_b)
	t.Log("Distance = ", dist)
	if dist > 160 || dist < 0 {
		t.Error("Incorrect distance measured")
	}
}
