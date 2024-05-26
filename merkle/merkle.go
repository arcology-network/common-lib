package merkle

import (
	"bytes"
	"crypto/sha256"

	"github.com/arcology-network/common-lib/common"
	slice "github.com/arcology-network/common-lib/exp/slice"
	"github.com/arcology-network/common-lib/mempool"
)

const (
	bufferSize  = 1024
	concurrency = 4
)

type Merkle struct {
	branch  uint32
	nodes   [][]*Node
	buffers [concurrency][]byte
	hasher  interface{ Hash([]byte) []byte }
	encoder interface{ Encode([][]byte) []byte }
	mempool *mempool.Mempool[*Node]
}

func NewMerkle(numBranches int, encoder interface{ Encode([][]byte) []byte }, hasher interface{ Hash([]byte) []byte }) *Merkle {
	merkle := &Merkle{
		branch:  uint32(numBranches),
		nodes:   [][]*Node{},
		hasher:  hasher,
		encoder: encoder,
	}
	for i := range merkle.buffers {
		merkle.buffers[i] = make([]byte, 0, bufferSize)
	}
	return merkle
}

func (this *Merkle) Hasher() interface{ Hash([]byte) []byte }                { return this.hasher }
func (this *Merkle) SetHasher(hasher interface{ Hash([]byte) []byte })       { this.hasher = hasher }
func (this *Merkle) Encoder() interface{ Encode([][]byte) []byte }           { return this.encoder }
func (this *Merkle) SetEncoder(encoder interface{ Encode([][]byte) []byte }) { this.encoder = encoder }

func (this *Merkle) Reset() *Merkle {
	for i := range this.nodes {
		this.nodes[i] = this.nodes[i][:0]
	}
	this.nodes = this.nodes[:0]
	for i := range this.buffers {
		this.buffers[i] = this.buffers[i][:0]
	}
	return this
}

func (this *Merkle) BuildParent(id uint32, children []*Node, index int, mempool *mempool.Mempool[*Node]) *Node {
	this.buffers[index] = slice.Concate(children, func(node *Node) []byte { return node.hash })

	// parent := mempool.Get().(*Node)
	parent := new(Node)
	parent.Init(id, children[0].level+1, this.hasher.Hash(this.buffers[index]))
	for i, v := range children {
		parent.children = append(parent.children, v.id)
		children[i].parent = parent.id
	}
	return parent
}

func (this *Merkle) Build(id uint32, children []*Node) []*Node {
	nodes := make([]*Node, len(children)/int(this.branch))
	worker := func(start, end, index int, args ...interface{}) {
		// mempool := this.mempool.GetPool(index)
		for i := start; i < end; i++ {
			nodes[i] = this.BuildParent(id+uint32(i), children[(i)*int(this.branch):(i+1)*int(this.branch)], index, nil)
		}
	}
	common.ParallelWorker(len(nodes), common.IfThen(len(children) < 64, 1, concurrency), worker)
	return nodes
}

func (*Merkle) Pad(original []*Node, n int) []*Node {
	for len(original)%n != 0 {
		original = append(original, original[len(original)-1])
	}
	return original
}

func (this *Merkle) newLeafNodes(data [][]byte) []*Node {
	leafNodes := make([]*Node, len(data))
	worker := func(start, end, index int, args ...interface{}) {
		// mempool := this.mempool.GetPool(index)
		for i := start; i < end; i++ {
			// leafNodes[i] = mempool.Get().(*Node)
			leafNodes[i] = new(Node)
			leafNodes[i].Init(uint32(i), 0, this.hasher.Hash(data[i]))
		}
	}
	common.ParallelWorker(len(data), common.IfThen(len(data) < 1024, 1, concurrency), worker)
	return leafNodes
}

func (this *Merkle) Init(data [][]byte, mempool interface{}) {
	if len(data) == 0 {
		return
	}

	// this.mempool = mempool
	if len(data) == 1 {
		// node := this.mempool.Get().(*Node)
		node := new(Node)
		node.Init(uint32(0), 0, this.hasher.Hash(data[0]))
		this.nodes = [][]*Node{{node}}
		return
	}

	// Initialized the leaf nodes
	this.nodes = append(this.nodes, this.Pad(this.newLeafNodes(data), int(this.branch)))

	// Build the non-leaf nodes
	for {
		nodes := this.Build(0, this.nodes[len(this.nodes)-1])
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

func (this *Merkle) GetProofNodes(key []byte) []*Node {
	hash := this.hasher.Hash(key)
	depth := 0
	mainPath := []*Node{}
	for _, v := range this.nodes[depth] {
		rgt := v.hash
		if bytes.Equal(hash[:], rgt[:]) {
			mainPath = append(mainPath, v)
			break
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
		if !slice.Contains(proofs[i], seed, bytes.Equal) {
			return false
		}
		seed = this.hasher.Hash(this.encoder.Encode(proofs[i]))
	}
	return bytes.Equal(seed[:], root[:])
}

func (this *Merkle) NodesToHashes(path []*Node) ([][]byte, [][][]byte) {
	hashes := [][][]byte{}
	subroots := make([][]byte, len(path))
	for i, v := range path {
		if childHashes := slice.Transform(this.GetChildrenOf(v), func(_ int, v *Node) []byte { return (*v).hash }); len(childHashes) > 0 {
			subroots[i] = this.hasher.Hash(this.encoder.Encode(childHashes))
			hashes = append(hashes, childHashes)
		}
	}
	return subroots, hashes
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
