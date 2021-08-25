package merkle

import (
	"bytes"
)

const (
	UINT32_MAX = 0xffffffff
)

type Merkle struct {
	branch uint32
	nodes  [][]*Node
	hasher func([]byte) []byte
}

func (this *Merkle) ExtractParent(id uint32, children []*Node) *Node {
	bytes := []byte{}
	for _, v := range children {
		hash := v.hash
		bytes = append(bytes, hash[:]...)
	}

	parent := NewNode(id, children[0].level+1, this.hasher(bytes))
	for i, v := range children {
		parent.children = append(parent.children, v.id)
		children[i].parent = parent.id
	}
	return parent
}

func (this *Merkle) Build(id uint32, children []*Node, n int) []*Node {
	nodes := []*Node{}
	for i := 1; i <= len(children)/n; i++ {
		nodes = append(nodes, this.ExtractParent(id, children[(i-1)*n:i*n]))
		id++
	}
	return nodes
}

func (*Merkle) Pad(original []*Node, n int) []*Node {
	if len(original) == n {
		return original
	}

	for i := 0; i < len(original)%n; i++ {
		original = append(original, original[len(original)-1])
	}
	return original
}

func NewMerkle(n int, hasher func([]byte) []byte) *Merkle {
	merkle := &Merkle{
		branch: uint32(n),
		nodes:  [][]*Node{},
		hasher: hasher,
	}
	return merkle
}

func (this *Merkle) Init(data [][]byte) {
	if len(data) == 1 {
		node := NewNode(uint32(0), 0, this.hasher([]byte(data[0])))
		this.nodes = [][]*Node{{node}}
		return
	}

	// Insert the leaf nodes
	leafNodes := []*Node{}
	for i, bytes := range data {
		leafNodes = append(leafNodes, NewNode(uint32(i), 0, this.hasher(bytes)))
	}
	this.nodes = append(this.nodes, this.Pad(leafNodes, int(this.branch)))

	// Build the non-leaf nodes
	for {
		nodes := this.Build(0, this.nodes[len(this.nodes)-1], int(this.branch))
		if len(nodes) == 1 {
			this.nodes = append(this.nodes, nodes)
			break
		}
		this.nodes = append(this.nodes, this.Pad(nodes, int(this.branch)))
	}
}

func (this *Merkle) GetChildrenOf(node *Node) []*Node {
	nodes := []*Node{}
	for _, childID := range node.children {
		nodes = append(nodes, this.nodes[node.level-1][childID])
	}
	return nodes
}

func (this *Merkle) GetRoot() []byte {
	if len(this.nodes) == 0 {
		return []byte{}
	}
	return this.nodes[len(this.nodes)-1][0].hash
}

func (this *Merkle) GetProofNodes(hash []byte) []*Node {
	depth := 0
	mainPath := []*Node{}
	for _, v := range this.nodes[depth] {
		rgt := v.hash
		if bytes.Equal(hash[:], rgt[:]) {
			mainPath = append(mainPath, v)
		}
	}

	for len(mainPath) > 0 {
		parentID := mainPath[len(mainPath)-1].parent
		lvl := mainPath[len(mainPath)-1].level + 1
		if int(lvl) >= len(this.nodes) {
			break
		}
		mainPath = append(mainPath, this.nodes[lvl][parentID])
	}
	return mainPath
}

func (this *Merkle) Verify(proofs [][][]byte, root []byte, seed []byte) bool {
	for i := 0; i < len(proofs); i++ {
		idx := this.IfContains(proofs[i], seed)
		if UINT32_MAX == idx {
			return false
		}
		seed = this.ComputeHash(proofs[i])
	}
	return bytes.Equal(seed[:], root[:])
}

func (this *Merkle) ComputeHash(hashes [][]byte) []byte {
	buffer := []byte{}
	for j := 0; j < len(hashes); j++ {
		buffer = append(buffer, hashes[j][:]...)
	}
	return this.hasher(buffer)
}

func (this *Merkle) IfContains(target [][]byte, seed []byte) uint32 {
	for i, v := range target {
		if bytes.Equal(seed[:], v[:]) {
			return uint32(i)
		}
	}
	return UINT32_MAX
}
