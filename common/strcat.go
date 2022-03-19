package common

import (
	"bytes"
	"sync"
)

var bufPool = sync.Pool{
	New: func() interface{} {
		return bytes.NewBuffer([]byte{})
	},
}

func StrCat(ss ...string) string {
	if len(ss) <= 1 {
		panic("misuse of StrCat, len(ss) must be greater or equal to 2.")
	}

	buf := bufPool.Get().(*bytes.Buffer)
	buf.Reset()
	for _, s := range ss {
		buf.WriteString(s)
	}

	str := buf.String()
	bufPool.Put(buf)
	return str
}
