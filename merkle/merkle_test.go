package merkle

import (
	"fmt"
	"reflect"
	"testing"

	mempool "github.com/arcology-network/common-lib/exp/mempool"
)

func TestBinaryMerkle(t *testing.T) { // Create a new merkle tree with 2 branches(binary) under each non-leaf node and using Sha256{} hashing algorithm
	bytes := [][]byte{}
	for i := 0; i < 15; i++ {
		bytes = append(bytes, []byte(fmt.Sprint(i)))
	}

	in := NewMerkle(2, Concatenator{}, Sha256{})
	nodePool := mempool.NewMempool[*Node](1, 2, func() *Node {
		return NewNode()
	}, func(v *Node) {})

	in.Init(bytes, nodePool)

	merkleBytes := in.Encode()                          // encode the tree
	merkle := (&Merkle{}).Decode(merkleBytes).(*Merkle) // decode the tree
	merkle.SetEncoder(Concatenator{})

	// if !reflect.DeepEqual(in.nodes, merkle.nodes) {
	// 	t.Error("Keys don't match")
	// }

	_, proofs := merkle.NodesToHashes(merkle.GetProofNodes([]byte(fmt.Sprint(0))))

	// In the original tree
	target := Sha256{}.Hash([]byte(fmt.Sprint(0)))
	if !merkle.Verify(proofs[:], merkle.GetRoot(), target[:]) {
		t.Error("Error: Merkle Proofs weren't found")
	}

	// In the decoded tree
	if !in.Verify(proofs[:], in.GetRoot(), target[:]) {
		t.Error("Error: Merkle Proofs weren't found")
	}
}

func TestQuadMerkle(t *testing.T) { // Create a new merkle tree with 4 branches under each non-leaf node and using Sha256{} hashing algorithm
	bytes := [][]byte{}
	for i := 0; i < 17; i++ {
		bytes = append(bytes, []byte(fmt.Sprint(i)))
	}

	merkle := NewMerkle(4, Concatenator{}, Sha256{})
	nodePool := mempool.NewMempool[*Node](1, 2, func() *Node {
		return NewNode()
	}, func(v *Node) {})

	merkle.Init(bytes, nodePool)

	proofNodes := merkle.GetProofNodes([]byte(fmt.Sprint(16)))
	_, proofs := merkle.NodesToHashes(proofNodes)

	target := Sha256{}.Hash([]byte(fmt.Sprint(16)))
	if !merkle.Verify(proofs[:], merkle.GetRoot(), target[:]) {
		t.Error("Error: Merkle Proofs weren't found")
	}
}

func TestOctodecMerkle(t *testing.T) { // Create a new merkle tree with 16 branches under each non-leaf node and using Keccak256 hashing algorithm
	bytes := [][]byte{}
	for i := 0; i < 9; i++ {
		bytes = append(bytes, []byte(fmt.Sprint(i)))
	}

	merkle := NewMerkle(8, Concatenator{}, Sha256{})
	nodePool := mempool.NewMempool[*Node](1, 2, func() *Node {
		return NewNode()
	}, func(v *Node) {})
	merkle.Init(bytes, nodePool)

	proofNodes := merkle.GetProofNodes([]byte(fmt.Sprint(0)))
	_, proofs := merkle.NodesToHashes(proofNodes)

	target := Sha256{}.Hash([]byte(fmt.Sprint(0)))
	if !merkle.Verify(proofs[:], merkle.GetRoot(), target[:]) {
		t.Error("Error: Merkle Proofs weren't found")
	}
}

func TestKeccakOctodecMerkle(t *testing.T) { // Create a new merkle tree with 16 branches under each non-leaf node and using Keccak256 hashing algorithm
	bytes := [][]byte{}
	for i := 0; i < 8; i++ {
		bytes = append(bytes, []byte(fmt.Sprint(i)))
	}

	merkle := NewMerkle(8, Concatenator{}, Keccak256{})
	nodePool := mempool.NewMempool[*Node](1, 2, func() *Node {
		return NewNode()
	}, func(v *Node) {})
	merkle.Init(bytes, nodePool)

	proofNodes := merkle.GetProofNodes([]byte(fmt.Sprint(0)))
	_, proofs := merkle.NodesToHashes(proofNodes)

	target := Keccak256{}.Hash([]byte(fmt.Sprint(0)))
	if !merkle.Verify(proofs[:], merkle.GetRoot(), target[:]) {
		t.Error("Error: Merkle Proofs weren't found")
	}
}

func TestKeccakOctodecMerkleSingleEntry(t *testing.T) { // Create a new merkle tree with 16 branches under each non-leaf node and using Keccak256 hashing algorithm
	bytes := [][]byte{[]byte(fmt.Sprint(0))}

	merkle := NewMerkle(10, Concatenator{}, Sha256{})
	nodePool := mempool.NewMempool[*Node](1, 2, func() *Node {
		return NewNode()
	}, func(v *Node) {})
	merkle.Init(bytes, nodePool)

	proofNodes := merkle.GetProofNodes([]byte(fmt.Sprint(0)))
	_, proofs := merkle.NodesToHashes(proofNodes)

	target := Sha256{}.Hash([]byte(fmt.Sprint(0)))
	if !merkle.Verify(proofs[:], merkle.GetRoot(), target[:]) {
		t.Error("Error: Merkle Proofs weren't found")
	}
}

func TestKeccakHexadecaMerkleSingleEntry(t *testing.T) { // Create a new merkle tree with 16 branches under each non-leaf node and using Keccak256 hashing algorithm
	bytes := [][]byte{[]byte(fmt.Sprint(0))}

	merkle := NewMerkle(16, Concatenator{}, Keccak256{})
	nodePool := mempool.NewMempool[*Node](1, 2, func() *Node {
		return NewNode()
	}, func(v *Node) {})
	merkle.Init(bytes, nodePool)

	proofNodes := merkle.GetProofNodes([]byte(fmt.Sprint(0)))
	_, proofs := merkle.NodesToHashes(proofNodes)

	target := Keccak256{}.Hash([]byte(fmt.Sprint(0)))
	if !merkle.Verify(proofs[:], merkle.GetRoot(), target[:]) {
		t.Error("Error: Merkle Proofs weren't found")
	}
}

func TestKeccakHexadecaMerkleMultiEntry(t *testing.T) { // Create a new merkle tree with 16 branches under each non-leaf node and using Keccak256 hashing algorithm
	bytes := [][]byte{}
	for i := 0; i < 1000; i++ {
		bytes = append(bytes, []byte(fmt.Sprint(i)))
	}

	merkle := NewMerkle(8, Concatenator{}, Sha256{})
	nodePool := mempool.NewMempool[*Node](1, 2, func() *Node {
		return NewNode()
	}, func(v *Node) {})
	merkle.Init(bytes, nodePool)

	proofNodes := merkle.GetProofNodes([]byte(fmt.Sprint(999)))
	_, proofs := merkle.NodesToHashes(proofNodes)

	target := Sha256{}.Hash([]byte(fmt.Sprint(999)))
	if !merkle.Verify(proofs[:], merkle.GetRoot(), target[:]) {
		t.Error("Error: Merkle Proofs weren't found")
	}
}

func TestKeccakDotriacontaMerkleSingleEntry(t *testing.T) { // Create a new merkle tree with 16 branches under each non-leaf node and using Keccak256 hashing algorithm
	bytes := [][]byte{[]byte(fmt.Sprint(0))}

	merkle := NewMerkle(32, Concatenator{}, Keccak256{})
	nodePool := mempool.NewMempool[*Node](1, 2, func() *Node {
		return NewNode()
	}, func(v *Node) {})
	merkle.Init(bytes, nodePool)

	proofNodes := merkle.GetProofNodes([]byte(fmt.Sprint(0)))
	_, proofs := merkle.NodesToHashes(proofNodes)

	target := Keccak256{}.Hash([]byte(fmt.Sprint(0)))
	if !merkle.Verify(proofs[:], merkle.GetRoot(), target[:]) {
		t.Error("Error: Merkle Proofs weren't found")
	}
}

func TestMerkleRootConsistency(t *testing.T) {
	bytes := make([][]byte, 0)
	for j := 0; j < 6; j++ {
		bytes = append(bytes, []byte(fmt.Sprint(j)))
	}
	nodePool := mempool.NewMempool[*Node](1, 2, func() *Node {
		return NewNode()
	}, func(v *Node) {})
	merkle := NewMerkle(32, Concatenator{}, Sha256{})
	merkle.Init(bytes, nodePool)

	nodePool2 := mempool.NewMempool[*Node](1, 2, func() *Node {
		return NewNode()
	}, func(v *Node) {})

	tree0 := NewMerkle(2, Concatenator{}, Sha256{})
	tree0.Init(bytes, nodePool2)
	r0 := tree0.GetRoot()

	tree1 := NewMerkle(2, Concatenator{}, Sha256{})
	nodePool3 := mempool.NewMempool[*Node](1, 2, func() *Node {
		return NewNode()
	}, func(v *Node) {})

	tree1.Init(bytes, nodePool3)
	r1 := tree1.GetRoot()
	if !reflect.DeepEqual(r0, r1) {
		t.Error("Roots mismatch")
	}
}

func TestMerklePaths(t *testing.T) {
	bytes := make([][]byte, 1000)
	for i := 0; i < len(bytes); i++ {
		bytes[i] = []byte(fmt.Sprint(i))
	}
	merkle := NewMerkle(8, Concatenator{}, Sha256{})
	nodePool := mempool.NewMempool[*Node](1, 2, func() *Node {
		return NewNode()
	}, func(v *Node) {})
	merkle.Init(bytes, nodePool)

	merkleBytes := merkle.Encode()                     // encode the tree
	merkle = (&Merkle{}).Decode(merkleBytes).(*Merkle) // decode the tree
	merkle.SetEncoder(Concatenator{})

	for i := 0; i < len(bytes); i++ {
		proofNodes := merkle.GetProofNodes(bytes[i])
		_, proofs := merkle.NodesToHashes(proofNodes)

		target := Sha256{}.Hash(bytes[i])
		if !merkle.Verify(proofs[:], merkle.GetRoot(), target[:]) {
			t.Error("Error: Merkle Pro ofs weren't found")
		}
	}
}
