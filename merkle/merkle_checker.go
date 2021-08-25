package merkle

import (
	"bytes"
	"crypto/sha256"
)

func (this *Merkle) NodesToHashes(path []*Node) [][][]byte {
	hashes := [][][]byte{}
	for _, v := range path {
		lvlHashes := [][]byte{}
		nodes := this.GetChildrenOf(v)
		for _, node := range nodes {
			lvlHashes = append(lvlHashes, node.hash)
		}

		if len(lvlHashes) > 0 {
			hashes = append(hashes, lvlHashes)
		}
	}
	return hashes
}

func (this *Merkle) CheckStructure() []*Node {
	nodeErrs := []*Node{}
	for i := 1; i < len(this.nodes); i++ {
		for _, node := range this.nodes[i] {
			if !this.CheckChildren(node) {
				nodeErrs = append(nodeErrs, node)
			}
		}
	}
	return nodeErrs
}

func (this *Merkle) CheckChildren(node *Node) bool {
	if node.level == 0 {
		return true
	}

	buffer := []byte{}
	for _, child := range node.children {
		node := this.nodes[node.level-1][child]
		buffer = append(buffer, node.hash[:]...)
	}

	hash256 := sha256.Sum256(buffer)
	return bytes.Equal(hash256[:], node.hash[:])
}
