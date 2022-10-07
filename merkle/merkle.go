package merkle

import (
	"bytes"
	"math"

	"github.com/HPISTechnologies/common-lib/common"
	"github.com/HPISTechnologies/common-lib/mempool"
)

const (
	bufferSize  = 1024
	concurrency = 4
)

type Merkle struct {
	branch  uint32
	nodes   [][]*Node
	buffers [concurrency][]byte
	hasher  func([]byte) []byte
	mempool *mempool.Mempool
}

func NewMerkle(n int, hasher func([]byte) []byte) *Merkle {
	merkle := &Merkle{
		branch: uint32(n),
		nodes:  [][]*Node{},
		hasher: hasher,
	}
	for i := range merkle.buffers {
		merkle.buffers[i] = make([]byte, 0, bufferSize)
	}
	return merkle
}

func (this *Merkle) Reset() {
	for i := range this.nodes {
		this.nodes[i] = this.nodes[i][:0]
	}
	this.nodes = this.nodes[:0]
	for i := range this.buffers {
		this.buffers[i] = this.buffers[i][:0]
	}
}

func (this *Merkle) ExtractParent(id uint32, children []*Node, index int, mempool *mempool.Mempool) *Node {
	this.buffers[index] = this.buffers[index][:0]
	for _, v := range children {
		this.buffers[index] = append(this.buffers[index], v.hash[:]...)
	}

	parent := mempool.Get().(*Node)
	parent.Init(id, children[0].level+1, this.hasher(this.buffers[index]))
	for i, v := range children {
		parent.children = append(parent.children, v.id)
		children[i].parent = parent.id
	}
	return parent
}

func (this *Merkle) Build(id uint32, children []*Node, n int) []*Node {
	if len(children) < 64 {
		return this.singleThreadedBuild(id, children, n)
	} else {
		return this.multiThreadedBuild(id, children, n)
	}
}

func (this *Merkle) singleThreadedBuild(id uint32, children []*Node, n int) []*Node {
	nodes := make([]*Node, len(children)/n)
	for i := 0; i < len(children)/n; i++ {
		nodes[i] = this.ExtractParent(id+uint32(i), children[(i)*n:(i+1)*n], 0, this.mempool)
	}
	return nodes
}

func (this *Merkle) multiThreadedBuild(id uint32, children []*Node, n int) []*Node {
	nodes := make([]*Node, len(children)/n)
	worker := func(start, end, index int, args ...interface{}) {
		mempool := this.mempool.GetTlsMempool(index)
		for i := start; i < end; i++ {
			nodes[i] = this.ExtractParent(id+uint32(i), children[(i)*n:(i+1)*n], index, mempool)
		}
	}
	common.ParallelWorker(len(nodes), concurrency, worker)
	return nodes
}

func (*Merkle) Pad(original []*Node, n int) []*Node {
	for len(original)%n != 0 {
		original = append(original, original[len(original)-1])
	}
	return original
}

func (this *Merkle) createLeafNodeSingleThreaded(data [][]byte) []*Node {
	leafNodes := make([]*Node, len(data))
	for i := 0; i < len(data); i++ {
		leafNodes[i] = this.mempool.Get().(*Node)
		leafNodes[i].Init(uint32(i), 0, this.hasher(data[i]))
	}
	return leafNodes
}

func (this *Merkle) createLeafNodeMultiThreaded(data [][]byte) []*Node {
	leafNodes := make([]*Node, len(data))
	worker := func(start, end, index int, args ...interface{}) {
		mempool := this.mempool.GetTlsMempool(index)
		for i := start; i < end; i++ {
			leafNodes[i] = mempool.Get().(*Node)
			leafNodes[i].Init(uint32(i), 0, this.hasher(data[i]))
		}
	}
	common.ParallelWorker(len(data), concurrency, worker)
	return leafNodes
}

func (this *Merkle) Init(data [][]byte, mempool *mempool.Mempool) {
	this.mempool = mempool
	if len(data) == 1 {
		node := this.mempool.Get().(*Node)
		node.Init(uint32(0), 0, this.hasher(data[0]))
		this.nodes = [][]*Node{{node}}
		return
	}

	var leafNodes []*Node
	if len(data) < 1024 {
		leafNodes = this.Pad(this.createLeafNodeSingleThreaded(data), int(this.branch))
	} else {
		leafNodes = this.Pad(this.createLeafNodeMultiThreaded(data), int(this.branch))
	}
	this.nodes = append(this.nodes, leafNodes)

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
		idx := this.IfContains(proofs[i], seed)
		if math.MaxUint32 == idx {
			return false
		}
		seed = this.ComputeHash(proofs[i])
	}
	return bytes.Equal(seed[:], root[:])
}

func (this *Merkle) ComputeHash(hashes [][]byte) []byte {
	if len(hashes) == 0 {
		return []byte{}
	}

	buffer := make([]byte, 0, len(hashes)*len(hashes[0]))
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
	return math.MaxUint32
}
