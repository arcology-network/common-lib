package merkle

import (
	"fmt"
	"reflect"
	"testing"
)

func TestBinaryMerkle(t *testing.T) { // Create a new merkle tree with 2 branches(binary) under each non-leaf node and using sha256 hashing algorithm
	bytes := [][]byte{}
	for i := 0; i < 15; i++ {
		bytes = append(bytes, []byte(fmt.Sprint(i)))
	}

	in := NewMerkle(2, Sha256)
	in.Init(bytes)

	merkleBytes := in.Encode()                          // encode the tree
	merkle := (&Merkle{}).Decode(merkleBytes).(*Merkle) // decode the tree

	if !reflect.DeepEqual(in.nodes, merkle.nodes) {
		t.Error("Keys don't match")
	}

	proofs := merkle.NodesToHashes(merkle.GetProofNodes(Sha256([]byte(fmt.Sprint(0)))))

	// In the original tree
	seed := Sha256([]byte(fmt.Sprint(0)))
	if !merkle.Verify(proofs[:], merkle.GetRoot(), seed[:]) {
		t.Error("Error: Merkle Proofs weren't found")
	}

	// In the decoded tree
	if !in.Verify(proofs[:], in.GetRoot(), seed[:]) {
		t.Error("Error: Merkle Proofs weren't found")
	}
}

func TestQuadMerkle(t *testing.T) { // Create a new merkle tree with 4 branches under each non-leaf node and using sha256 hashing algorithm
	bytes := [][]byte{}
	for i := 0; i < 15; i++ {
		bytes = append(bytes, []byte(fmt.Sprint(i)))
	}

	merkle := NewMerkle(4, Sha256)
	merkle.Init(bytes)

	proofNodes := merkle.GetProofNodes(Sha256([]byte(fmt.Sprint(0))))
	proofs := merkle.NodesToHashes(proofNodes)

	seed := Sha256([]byte(fmt.Sprint(0)))
	if !merkle.Verify(proofs[:], merkle.GetRoot(), seed[:]) {
		t.Error("Error: Merkle Proofs weren't found")
	}
}

func TestOctodecMerkle(t *testing.T) { // Create a new merkle tree with 16 branches under each non-leaf node and using Keccak256 hashing algorithm
	bytes := [][]byte{}
	for i := 0; i < 8; i++ {
		bytes = append(bytes, []byte(fmt.Sprint(i)))
	}

	merkle := NewMerkle(2, Sha256)
	merkle.Init(bytes)

	proofNodes := merkle.GetProofNodes(Sha256([]byte(fmt.Sprint(0))))
	proofs := merkle.NodesToHashes(proofNodes)

	seed := Sha256([]byte(fmt.Sprint(0)))
	if !merkle.Verify(proofs[:], merkle.GetRoot(), seed[:]) {
		t.Error("Error: Merkle Proofs weren't found")
	}
}

func TestKeccakOctodecMerkle(t *testing.T) { // Create a new merkle tree with 16 branches under each non-leaf node and using Keccak256 hashing algorithm
	bytes := [][]byte{}
	for i := 0; i < 8; i++ {
		bytes = append(bytes, []byte(fmt.Sprint(i)))
	}

	merkle := NewMerkle(2, Keccak256)
	merkle.Init(bytes)

	proofNodes := merkle.GetProofNodes(Keccak256([]byte(fmt.Sprint(0))))
	proofs := merkle.NodesToHashes(proofNodes)

	seed := Keccak256([]byte(fmt.Sprint(0)))
	if !merkle.Verify(proofs[:], merkle.GetRoot(), seed[:]) {
		t.Error("Error: Merkle Proofs weren't found")
	}
}

func TestKeccakOctodecMerkleSingleEntry(t *testing.T) { // Create a new merkle tree with 16 branches under each non-leaf node and using Keccak256 hashing algorithm
	bytes := [][]byte{[]byte(fmt.Sprint(0))}

	merkle := NewMerkle(2, Keccak256)
	merkle.Init(bytes)

	proofNodes := merkle.GetProofNodes(Keccak256([]byte(fmt.Sprint(0))))
	proofs := merkle.NodesToHashes(proofNodes)

	seed := Keccak256([]byte(fmt.Sprint(0)))
	if !merkle.Verify(proofs[:], merkle.GetRoot(), seed[:]) {
		t.Error("Error: Merkle Proofs weren't found")
	}
}
