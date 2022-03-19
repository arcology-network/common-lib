package merkle

import (
	"sync"
)

var nodePool = sync.Pool{
	New: func() interface{} {
		return new(Node)
	},
}
