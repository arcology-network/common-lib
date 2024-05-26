package addrcompressor

import (
	"math/rand"
	"strings"
	"time"

	common "github.com/arcology-network/common-lib/common"
)

func GetByDepth(originals []string, depth int) []string {
	keys := make([]string, len(originals))
	finder := func(start, end, idx int, args ...interface{}) {
		for i := start; i < end; i++ {
			idx := IndexN(originals[i], "/", depth)
			if idx >= 0 {
				keys[i] = originals[i][:idx]
			}
		}
	}
	common.ParallelWorker(len(originals), 4, finder)
	return keys
}

func GetBetweenDepths(originals []string, depth0 int, depth1 int) []string {
	keys := make([]string, len(originals))
	finder := func(start, end, idx int, args ...interface{}) {
		for i := start; i < end; i++ {
			idx0 := IndexN(originals[i], "/", depth0)
			idx1 := IndexN(originals[i], "/", depth1)
			if idx0 >= 0 && idx1 >= 0 {
				keys[i] = originals[i][idx0+1 : idx1]
			}
		}
	}
	common.ParallelWorker(len(originals), 4, finder)
	return keys
}

func LocateWildcards(patten string, wildcard string) []int {
	pos := []int{}
	wildcardPos := 0
	for {
		if wildcardPos = strings.Index(patten, wildcard); wildcardPos > -1 {
			pos = append(pos, strings.Count(patten[:wildcardPos], "/")) // num of delimiters before the next wildcard
			patten = patten[wildcardPos+len(wildcard):]
		} else {
			break
		}
	}
	return pos
}

// Find the nth occurrence of a target string
func IndexN(line string, target string, n int) int {
	if n < 0 {
		return 0
	}

	if n == len(line) {
		return len(line)
	}

	pos := 0
	for {
		if i := strings.Index(line[pos+1:], target); i > -1 && n >= 0 {
			pos += i + 1
			n--
		} else {
			break
		}
	}

	if n > 0 {
		return -1
	}
	return pos
}

// Generate a random account, testing only
func RandomAccount() string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyz0123456789")
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, 40)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func AliceAccount() string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyz0123456789")
	// rand.Seed(1)
	b := make([]rune, 40)
	for i := range b {
		b[i] = letters[1]
	}
	return string(b)
}

func BobAccount() string {
	var letters = []rune("9876543210zyxwvutsrqponmlkjihgfedcba")

	// rand.Seed(2)
	b := make([]rune, 40)
	for i := range b {
		b[i] = letters[2]
	}
	return string(b)
}

// For testing and debugging only, not a performance option
func Deepcopy(source []string) []string {
	target := make([]string, len(source))
	finder := func(start, end, idx int, args ...interface{}) {
		for i := start; i < end; i++ {
			var builder strings.Builder
			builder.WriteString(source[i])
			target[i] = builder.String()
		}
	}
	common.ParallelWorker(len(source), 4, finder)
	return target
}
